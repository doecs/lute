// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

// EChartsJSONRenderer 描述了 JSON 渲染器。
type EChartsJSONRenderer struct {
	*BaseRenderer
}

// newEChartsJSONRenderer 创建一个 ECharts JSON 渲染器。
func (lute *Lute) newEChartsJSONRenderer(treeRoot *Node) Renderer {
	ret := &EChartsJSONRenderer{&BaseRenderer{rendererFuncs: map[int]RendererFunc{}, option: lute.options, treeRoot: treeRoot}}
	ret.rendererFuncs[NodeDocument] = ret.renderDocumentEChartsJSON
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphEChartsJSON
	ret.rendererFuncs[NodeText] = ret.renderTextEChartsJSON
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanEChartsJSON
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockEChartsJSON
	ret.rendererFuncs[NodeMathBlock] = ret.renderMathBlockEChartsJSON
	ret.rendererFuncs[NodeInlineMath] = ret.renderInlineMathEChartsJSON
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasisEChartsJSON
	ret.rendererFuncs[NodeStrong] = ret.renderStrongEChartsJSON
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquoteEChartsJSON
	ret.rendererFuncs[NodeHeading] = ret.renderHeadingEChartsJSON
	ret.rendererFuncs[NodeList] = ret.renderListEChartsJSON
	ret.rendererFuncs[NodeListItem] = ret.renderListItemEChartsJSON
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreakEChartsJSON
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreakEChartsJSON
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreakEChartsJSON
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTMLEChartsJSON
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTMLEChartsJSON
	ret.rendererFuncs[NodeLink] = ret.renderLinkEChartsJSON
	ret.rendererFuncs[NodeImage] = ret.renderImageEChartsJSON
	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethroughEChartsJSON
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarkerEChartsJSON
	ret.rendererFuncs[NodeTable] = ret.renderTableEChartsJSON
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHeadEChartsJSON
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRowEChartsJSON
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCellEChartsJSON
	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicodeEChartsJSON
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImgEChartsJSON
	return ret
}

func (r *EChartsJSONRenderer) renderInlineMathEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Inline Math\nspan", false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderMathBlockEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Math Block\ndiv", false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderEmojiImgEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Emoji Img\n", false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderEmojiUnicodeEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Emoji Unicode\n", false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderTableCellEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Table Cell\ntd", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderTableRowEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Table Row\ntr", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderTableHeadEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Table Head\nthead", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderTableEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Table\ntable", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderStrikethroughEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Strikethrough\ndel", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderImageEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Image\nimg", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderLinkEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Link\na", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderHTMLEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("HTML Block\n", false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderInlineHTMLEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Inline HTML\n", false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderDocumentEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if (entering) {
		r.writeString(`{
  "name": "doc",
  "children": [`)
	} else {
		r.writeString("]}")
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderParagraphEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Paragraph\np", true)
		r.openChildren()
	} else {
		if nil != node.next {
			r.closeChildren(true)
		}
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderTextEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		text := fromItems(node.tokens)
		length := len(text)
		if 16 <= length {
			length = 16 // 不考虑 rune 解码
		}
		r.val("Text\n"+text[:length], false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderCodeSpanEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Code Span\ncode", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderEmphasisEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Emphasis\nem", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderStrongEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Strong\nstrong", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderBlockquoteEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Blockquote\nblockquote", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderHeadingEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		h := "<h" + " 123456"[node.headingLevel:node.headingLevel+1]
		r.val("Heading\n"+h, false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderListEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("List Item\n"+fromItems(node.listData.marker), false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderListItemEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("List Item\n["+fromItems(node.listData.marker)+"]", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderTaskListItemMarkerEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		check := " "
		if node.taskListItemChecked {
			check = "X"
		}
		r.val("Task List Item Marker\n["+check+"]", false)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderThematicBreakEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Thematic Break\nhr", false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderHardBreakEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Hard Break\nbr", false)

	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderSoftBreakEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Soft Break\n\\\\n", false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderCodeBlockEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.val("Code Block\npre.code", false)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) val(val string, comma bool) {
	r.writeString("\"name\": \"" + val + "\"")
	if comma {
		r.comma()
	}
}

func (r *EChartsJSONRenderer) openObj() {
	r.writeByte('{')
}

func (r *EChartsJSONRenderer) closeObj(comma bool) {
	r.writeByte('}')
	if comma {
		r.comma()
	}
}

func (r *EChartsJSONRenderer) openChildren() {
	r.writeString("\"children\": [")
}

func (r *EChartsJSONRenderer) closeChildren(comma bool) {
	r.writeByte(']')
	if comma {
		r.comma()
	}
}

func (r *EChartsJSONRenderer) comma() {
	r.writeString(",\n")
}