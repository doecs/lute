// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"unicode"
	"unicode/utf8"

	"lute/ast"
	"lute/util"
)

// ChinesePunct 会把文本节点 textNode 中的中文间的英文标点换成对应的中文标点。
func (r *BaseRenderer) ChinesePunct(textNode *ast.Node) {
	text := util.BytesToStr(textNode.Tokens)
	text = chinesePunct0(text)
	textNode.Tokens = util.StrToBytes(text)
}

func chinesePunct0(text string) (ret string) {
	runes := []rune(text)
	length := len(runes)
	for i, r := range runes {
		if ('.' == r || '!' == r || '?' == r) && i+1 < length {
			if '.' == runes[i+1] || '!' == runes[i+1] || '?' == runes[i+1] {
				// 连续英文标点符号出现在中文后不优化
				ret += string(r)
				continue
			} else if isFileExt(i+1, length, &runes) {
				// 中文.合法扩展名 的形式不进行转换
				ret += string(r)
				continue
			}
		}
		ret = chinesePunct00(ret, r)
	}
	return
}

func chinesePunct00(prefix string, nextChar rune) string {
	if 0 == len(prefix) {
		return string(nextChar)
	}

	nextCharIsEnglishComma := ',' == nextChar
	nextCharIsEnglishPeriod := '.' == nextChar
	nextCharIsEnglishColon := ':' == nextChar
	nextCharIsEnglishBang := '!' == nextChar
	nextCharIsEnglishQuestion := '?' == nextChar

	currentChar, size := utf8.DecodeLastRuneInString(prefix)
	if 1 == size && (',' == currentChar) && unicode.Is(unicode.Han, nextChar) {
		// test,测试 => test，测试
		return prefix[:len(prefix)-1] + "，" + string(nextChar)
	}

	if !nextCharIsEnglishComma && !nextCharIsEnglishPeriod && !nextCharIsEnglishColon && !nextCharIsEnglishBang && !nextCharIsEnglishQuestion {
		return prefix + string(nextChar)
	}

	if !unicode.Is(unicode.Han, currentChar) {
		return prefix + string(nextChar)
	}

	if nextCharIsEnglishComma {
		return prefix + "，"
	} else if nextCharIsEnglishPeriod {
		return prefix + "。"
	} else if nextCharIsEnglishColon {
		return prefix + "："
	} else if nextCharIsEnglishBang {
		return prefix + "！"
	} else if nextCharIsEnglishQuestion {
		return prefix + "？"
	}
	return prefix + string(nextChar)
}
