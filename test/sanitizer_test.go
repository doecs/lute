// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"testing"

	"lute"
)

var sanitizerTests = []parseTest{

	{"5", "<iframe src='javascript:parent.require(\"child_process\").exec(\"open -a Calculator\")'></iframe>", "<iframe></iframe>\n"},
	{"4", "![Escape SRC - onerror](\"onerror=\"alert('ImageOnError'))", "<p><img src=\"%22onerror=%22alert(&#39;ImageOnError&#39;)\" alt=\"Escape SRC - onerror\"/></p>\n"},
	{"3", "<EMBED SRC=\"data:image/svg+xml;base64,mock payload\" type=\"image/svg+xml\" AllowScriptAccess=\"always\"></EMBED>", "<p><embed></embed></p>\n"},
	{"2", "<FOo>bar", "<p><foo>bar</p>\n"},
	{"1", "<img onerror=\"alert(1)\" src=\"bar.png\" />", "<img src=\"bar.png\"/>\n"},
	{"0", "foo<script>alert(1)</script>bar", "<p>foo alert(1) bar</p>\n"},
}

func TestSanitizer(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.Sanitize = true

	for _, test := range sanitizerTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var sanitizerVditorTests = []parseTest{

	{"7", "<iframe src=\"//player.bilibili.com/player.html?aid=test&page=1\" scrolling=\"no\" border=\"0\" frameborder=\"no\" framespacing=\"0\" allowfullscreen=\"true\"> </iframe>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;iframe src=&quot;//player.bilibili.com/player.html?aid=test&amp;amp;page=1&quot;&gt; &lt;/iframe&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><iframe src=\"//player.bilibili.com/player.html?aid=test&amp;page=1\"> </iframe></pre></div>"},
	{"6", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;img src=\"test1<wbr>\" onerror=\"alert('XSS')\"&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"1\"><img src=\"test\" onerror=\"alert('XSS')\"></pre></div>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;img src=&quot;test1<wbr>&quot;&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><img src=\"test1\"></pre></div>"},
	{"5", "<iframe src='javascript:parent.require(\"child_process\").exec(\"open -a Calculator\")'></iframe>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;iframe&gt;&lt;/iframe&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><iframe></iframe></pre></div>"},
	{"4", "![Escape SRC - onerror](\"onerror=\"alert('ImageOnError'))", "<p data-block=\"0\"><img src=\"\" alt=\"Escape SRC - onerror\"/>\n</p>"},
	{"3", "<EMBED SRC=\"data:image/svg+xml;base64,mock payload\" type=\"image/svg+xml\" AllowScriptAccess=\"always\"></EMBED>", "<p data-block=\"0\">\u200b<span class=\"vditor-wysiwyg__block\" data-type=\"html-inline\"><code data-type=\"html-inline\">\u200b&lt;embed/&gt;</code></span>\u200b\n</p>"},
	{"2", "<FOo>bar", "<p data-block=\"0\">foobar\n</p>"},
	{"1", "<img onerror=\"alert(1)\" src=\"bar.png\" />", "<p data-block=\"0\"><img src=\"bar.png\" alt=\"\"/>\n</p>"},
	{"0", "foo<script>alert(1)</script>bar", "<p data-block=\"0\">foo\n</p><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>  </code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\">  </pre></div><p data-block=\"0\">bar\n</p>"},
}

func TestSanitizerVditor(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.Sanitize = true

	for _, test := range sanitizerVditorTests {
		html := luteEngine.SpinVditorDOM(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
