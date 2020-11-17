// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"lute/ast"
	"lute/lex"
)

// Caret 插入符 \u2038。
const Caret = "‸"

// Zwsp 零宽空格。
const Zwsp = "\u200b"

// Parse 会将 markdown 原始文本字节数组解析为一颗语法树。
func Parse(name string, markdown []byte, options *Options) (tree *Tree) {
	tree = &Tree{Name: name, Context: &Context{Option: options}}
	tree.Context.Tree = tree
	tree.lexer = lex.NewLexer(markdown)
	tree.Root = &ast.Node{Type: ast.NodeDocument}
	tree.parseBlocks()
	tree.parseInlines()
	tree.lexer = nil
	return
}

// Context 用于维护块级元素解析过程中使用到的公共数据。
type Context struct {
	Tree   *Tree    // 关联的语法树
	Option *Options // 解析渲染选项

	LinkRefDefs   map[string]*ast.Node // 链接引用定义集
	FootnotesDefs []*ast.Node          // 脚注定义集

	Tip                                                               *ast.Node // 末梢节点
	oldtip                                                            *ast.Node // 老的末梢节点
	currentLine                                                       []byte    // 当前行
	currentLineLen                                                    int       // 当前行长
	lineNum, offset, column, nextNonspace, nextNonspaceColumn, indent int       // 解析时用到的行号、下标、缩进空格数等
	indented, blank, partiallyConsumedTab, allClosed                  bool      // 是否是缩进行、空行等标识
	lastMatchedContainer                                              *ast.Node // 最后一个匹配的块节点
}

// InlineContext 描述了行级元素解析上下文。
type InlineContext struct {
	tokens     []byte     // 当前解析的 Tokens
	tokensLen  int        // 当前解析的 Tokens 长度
	pos        int        // 当前解析到的 token 位置
	lineNum    int        // 当前解析的起始行号
	columnNum  int        // 当前解析的起始列号
	delimiters *delimiter // 分隔符栈，用于强调解析
	brackets   *delimiter // 括号栈，用于图片和链接解析
}

// advanceOffset 用于移动 count 个字符位置，columns 指定了遇到 tab 时是否需要空格进行补偿偏移。
func (context *Context) advanceOffset(count int, columns bool) {
	var currentLine = context.currentLine
	var charsToTab, charsToAdvance int
	var c byte
	for 0 < count {
		c = currentLine[context.offset]
		if lex.ItemTab == c {
			charsToTab = 4 - (context.column % 4)
			if columns {
				context.partiallyConsumedTab = charsToTab > count
				if context.partiallyConsumedTab {
					charsToAdvance = count
				} else {
					charsToAdvance = charsToTab
					context.offset++
				}
				context.column += charsToAdvance
				count -= charsToAdvance
			} else {
				context.partiallyConsumedTab = false
				context.column += charsToTab
				context.offset++
				count--
			}
		} else {
			context.partiallyConsumedTab = false
			context.offset++
			context.column++ // 假定是 ASCII，因为块开始标记符都是 ASCII
			count--
		}
	}
}

// advanceNextNonspace 用于预移动到下一个非空字符位置。
func (context *Context) advanceNextNonspace() {
	context.offset = context.nextNonspace
	context.column = context.nextNonspaceColumn
	context.partiallyConsumedTab = false
}

// findNextNonspace 用于查找下一个非空字符。
func (context *Context) findNextNonspace() {
	i := context.offset
	cols := context.column

	var token byte
	for {
		token = context.currentLine[i]
		if lex.ItemSpace == token {
			i++
			cols++
		} else if lex.ItemTab == token {
			i++
			cols += 4 - (cols % 4)
		} else {
			break
		}
	}

	context.blank = lex.ItemNewline == token
	context.nextNonspace = i
	context.nextNonspaceColumn = cols
	context.indent = context.nextNonspaceColumn - context.column
	context.indented = context.indent >= 4
}

// closeUnmatchedBlocks 最终化所有未匹配的块节点。
func (context *Context) closeUnmatchedBlocks() {
	if !context.allClosed {
		for context.oldtip != context.lastMatchedContainer {
			parent := context.oldtip.Parent
			context.finalize(context.oldtip, context.lineNum-1)
			context.oldtip = parent
		}
		context.allClosed = true
	}
}

// finalize 执行 block 的最终化处理。调用该方法会将 context.Tip 置为 block 的父节点。
func (context *Context) finalize(block *ast.Node, lineNum int) {
	parent := block.Parent
	block.Close = true

	// 节点最终化处理。比如围栏代码块提取 info 部分；HTML 代码块剔除结尾空格；段落需要解析链接引用定义等。
	switch block.Type {
	case ast.NodeCodeBlock:
		codeBlockFinalize(block)
	case ast.NodeHTMLBlock:
		htmlBlockFinalize(block)
	case ast.NodeParagraph:
		insertTable := paragraphFinalize(block, context)
		if insertTable {
			return
		}
	case ast.NodeMathBlock:
		mathBlockFinalize(block)
	case ast.NodeList:
		listFinalize(block)
	}

	context.Tip = parent
}

// addChildMarker 将构造一个 NodeType 节点并作为子节点添加到末梢节点 context.Tip 上。
func (context *Context) addChildMarker(nodeType ast.NodeType, tokens []byte) (ret *ast.Node) {
	ret = &ast.Node{Type: nodeType, Tokens: tokens, Close: true}
	context.Tip.AppendChild(ret)
	return ret
}

// addChild 将构造一个 NodeType 节点并作为子节点添加到末梢节点 context.Tip 上。如果末梢不能接受子节点（非块级容器不能添加子节点），则最终化该末梢
// 节点并向父节点方向尝试，直到找到一个能接受该子节点的节点为止。添加完成后该子节点会被设置为新的末梢节点。
func (context *Context) addChild(nodeType ast.NodeType, offset int) (ret *ast.Node) {
	for !context.Tip.CanContain(nodeType) {
		context.finalize(context.Tip, context.lineNum-1) // 注意调用 finalize 会向父节点方向进行迭代
	}

	ret = &ast.Node{Type: nodeType}
	context.Tip.AppendChild(ret)
	context.Tip = ret
	return ret
}

// listsMatch 用户判断指定的 listData 和 itemData 是否可归属于同一个列表。
func (context *Context) listsMatch(listData, itemData *ast.ListData) bool {
	return listData.Typ == itemData.Typ &&
		((0 == listData.Delimiter && 0 == itemData.Delimiter) || listData.Delimiter == itemData.Delimiter) &&
		listData.BulletChar == itemData.BulletChar
}

// Tree 描述了 Markdown 抽象语法树结构。
type Tree struct {
	Name          string         // 名称，可以为空
	Root          *ast.Node      // 根节点
	Context       *Context       // 块级解析上下文
	lexer         *lex.Lexer     // 词法分析器
	inlineContext *InlineContext // 行级解析上下文
}

// Options 描述了一些列解析和渲染选项。
type Options struct {
	// GFMTable 设置是否打开“GFM 表”支持。
	GFMTable bool
	// GFMTaskListItem 设置是否打开“GFM 任务列表项”支持。
	GFMTaskListItem bool
	// GFMTaskListItemClass 作为 GFM 任务列表项类名，默认为 "vditor-task"。
	GFMTaskListItemClass string
	// GFMStrikethrough 设置是否打开“GFM 删除线”支持。
	GFMStrikethrough bool
	// GFMAutoLink 设置是否打开“GFM 自动链接”支持。
	GFMAutoLink bool
	// SoftBreak2HardBreak 设置是否将软换行（\n）渲染为硬换行（<br />）。
	SoftBreak2HardBreak bool
	// CodeSyntaxHighlight 设置是否对代码块进行语法高亮。
	CodeSyntaxHighlight bool
	// CodeSyntaxHighlightDetectLang bool
	CodeSyntaxHighlightDetectLang bool
	// CodeSyntaxHighlightInlineStyle 设置语法高亮是否为内联样式，默认不内联。
	CodeSyntaxHighlightInlineStyle bool
	// CodeSyntaxHightLineNum 设置语法高亮是否显示行号，默认不显示。
	CodeSyntaxHighlightLineNum bool
	// CodeSyntaxHighlightStyleName 指定语法高亮样式名，默认为 "github"。
	CodeSyntaxHighlightStyleName string
	// Footnotes 设置是否打开“脚注”支持。
	Footnotes bool
	// ToC 设置是否打开“目录”支持。
	ToC bool
	// HeadingID 设置是否打开“自定义标题 ID”支持。
	HeadingID bool
	// AutoSpace 设置是否对普通文本中的中西文间自动插入空格。
	// https://github.com/sparanoid/chinese-copywriting-guidelines
	AutoSpace bool
	// FixTermTypo 设置是否对普通文本中出现的术语进行修正。
	// https://github.com/sparanoid/chinese-copywriting-guidelines
	// 注意：开启术语修正的话会默认在中西文之间插入空格。
	FixTermTypo bool
	// ChinesePunct 设置是否对普通文本中出现中文后跟英文逗号句号等标点替换为中文对应标点。
	ChinesePunct bool
	// Emoji 设置是否对 Emoji 别名替换为原生 Unicode 字符。
	Emoji bool
	// AliasEmoji 存储 ASCII 别名到表情 Unicode 映射。
	AliasEmoji map[string]string
	// EmojiAlias 存储表情 Unicode 到 ASCII 别名映射。
	EmojiAlias map[string]string
	// EmojiSite 设置图片 Emoji URL 的路径前缀。
	EmojiSite string
	// HeadingAnchor 设置是否对标题生成链接锚点。
	HeadingAnchor bool
	// Terms 将传入的 terms 合并覆盖到已有的 Terms 字典。
	Terms map[string]string
	// Vditor 所见即所得支持
	VditorWYSIWYG bool
	// Vditor 即时渲染支持
	VditorIR bool
	// InlineMathAllowDigitAfterOpenMarker 设置内联数学公式是否允许起始 $ 后紧跟数字 https://github.com/b3log/lute/issues/38
	InlineMathAllowDigitAfterOpenMarker bool
	// LinkBase 设置链接、图片的基础路径。如果用户在链接或者图片地址中使用相对路径（没有协议前缀且不以 / 开头）并且 LinkBase 不为空则会用该值作为前缀。
	// 比如 LinkBase 设置为 http://domain.com/，对于 ![foo](bar.png) 则渲染为 <img src="http://domain.com/bar.png" alt="foo" />
	LinkBase string
	// VditorCodeBlockPreview 设置 Vditor 代码块是否需要渲染预览部分
	VditorCodeBlockPreview bool
	// RenderListMarker 设置在渲染 OL、UL 时是否添加 data-marker 属性 https://github.com/88250/lute/issues/48
	RenderListMarker bool
	// Setext 设置是否解析 Setext 标题 https://github.com/88250/lute/issues/50
	Setext bool
	// Sanitize 设置是否启用 XSS 安全过滤 https://github.com/88250/lute/issues/51
	Sanitize bool
	// ImageLazyLoading 设置图片懒加载时使用的图片路径，配置该字段后将启用图片懒加载。
	// 图片 src 的值会复制给新属性 data-src，然后使用该参数值作为 src 的值 https://github.com/88250/lute/issues/55
	ImageLazyLoading string
	// ChineseParagraphBeginningSpace 设置是否使用传统中文排版“段落开头空两格”
	ChineseParagraphBeginningSpace bool
}

func (context *Context) ParentTip() {
	if tip := context.Tip.Parent; nil != tip {
		context.Tip = context.Tip.Parent
	}
}
