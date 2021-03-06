//
// md2html_test.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// tests for md2html.go.
//
// Copyright © 2018 Frank Storbeck. All rights reserved.
// Code licensed under the BSD License:
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
package main

import (
	"testing"
)

func TestStyling(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "", want: ""},
		{s: "**bb**", want: "<strong>bb</strong>"},
		{s: "xx**bb**", want: "xx<strong>bb</strong>"},
		{s: "**bb**yy", want: "<strong>bb</strong>yy"},
		{s: "xx**bb**yy", want: "xx<strong>bb</strong>yy"},
		{s: "**pp", want: "**pp"},
		{s: "xx**pp", want: "xx<em></em>pp"},
		{s: "**ppyy", want: "**ppyy"},
		{s: "xx__ppyy", want: "xx<em></em>ppyy"},
		{s: "**bb**xx**bb**__bb__**bb**yy", want: "<strong>bb</strong>xx<strong>bb</strong><strong>bb</strong><strong>bb</strong>yy"},
		{s: "**bb***bb****", want: "<strong>bb</strong>*bb<strong></strong>"},
		{s: "~~bb~~", want: "<del>bb</del>"},
		{s: "xx~~bb~~", want: "xx<del>bb</del>"},
		{s: "~~bb~~yy", want: "<del>bb</del>yy"},
	}
	for _, tst := range tests {
		got := StrongEmDel(tst.s)
		if got != tst.want {
			t.Errorf("StrongEmDel(%q) generates:\n%q\nshould be:\n%q\n", tst.s, got, tst.want)
		}
	}
}

func TestUniCode(t *testing.T) {
	tests := []struct {
		s    string
		esc  bool
		want string
	}{
		{s: "a \\* b", esc: true, want: "a U+002A b"},
		{s: "\\* b", esc: true, want: "U+002A b"},
		{s: "a \\*", esc: true, want: "a U+002A"},
		{s: "a \\_ b", esc: true, want: "a \\_ b"},
		{s: "a \\* b", esc: false, want: "a \\U+002A b"},
		{s: "\\* b", esc: false, want: "\\U+002A b"},
		{s: "a \\*", esc: false, want: "a \\U+002A"},
		{s: "a \\_ b", esc: true, want: "a \\_ b"},
	}
	for _, tst := range tests {
		got := CodeUni(tst.s, []byte{'*'}, tst.esc)
		if got != tst.want {
			t.Errorf("CodeUni(%q, []byte{'*'}, %t) generates:\n%q\nshould be:\n%q\n", tst.s, tst.esc, got, tst.want)
		}
	}
}

func TestUnEscape(t *testing.T) {
	tests := []struct {
		s    string
		esc  bool
		want string
	}{
		{s: "a U+002A b", esc: false, want: "a * b"},
		{s: "U+002A b", esc: false, want: "* b"},
		{s: "a U+002A", esc: false, want: "a *"},
		{s: "a U+002B b", esc: false, want: "a U+002B b"},
		{s: "a U+002A b", esc: true, want: "a \\* b"},
		{s: "U+002A b", esc: true, want: "\\* b"},
		{s: "a U+002A", esc: true, want: "a \\*"},
		{s: "a U+002B b", esc: true, want: "a U+002B b"},
	}
	for _, tst := range tests {
		got := DecodeUni(tst.s, []byte{'*'}, tst.esc)
		if got != tst.want {
			t.Errorf("DecodeUni(%q, []byte{'*'}) generates:\n%q\nshould be:\n%q\n", tst.s, got, tst.want)
		}
	}
}

func TestOlnlyRunes(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{s: "", want: false},
		{s: "==", want: false},
		{s: "===", want: true},
		{s: "=======", want: true},
		{s: "=======x", want: false},
		{s: "====x==", want: false},
	}

	for _, tst := range tests {
		got := OnlyRunes(tst.s, '=')
		if got != tst.want {
			t.Errorf("'OnlyRunes(%q, '=')' generates: %t, should be: %t",
				tst.s, got, tst.want)
		}
	}
}

func TestCountLeading(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{s: "", want: 0},
		{s: "x", want: 0},
		{s: "#x", want: 1},
		{s: "######x", want: 6},
		{s: "#######x", want: 0},
		{s: " #####", want: 0},
	}

	for _, tst := range tests {
		n := CountLeading(tst.s, '#', 6)
		if n != tst.want {
			t.Errorf("'CountLeading(%q)' generates: %d, should be: %d",
				tst.s, n, tst.want)
		}
	}
}

func TestBuild(t *testing.T) {
	tests := []struct {
		s    []string
		want string
	}{
		// Headers
		{s: []string{"hdr1", "==="}, want: "r{h1:id=\"hdr1\"{hdr1} p{}}"},
		{s: []string{"hdr2", "---"}, want: "r{h2:id=\"hdr2\"{hdr2} p{}}"},
		{s: []string{"aa", "# hdr1", "bb"}, want: "r{p{aa} h1:id=\"hdr1\"{hdr1} p{bb}}"},
		{s: []string{"aa", "### hdr3", "bb"}, want: "r{p{aa} h3:id=\"hdr3\"{hdr3} p{bb}}"},
		{s: []string{"###### hdr6"}, want: "r{h6:id=\"hdr6\"{hdr6} p{}}"},
		{s: []string{"####### hdr7"}, want: "r{p{####### hdr7}}"},

		// Quoting
		{s: []string{"> quote"},
			want: "r{blockquote{quote }}"},
		{s: []string{"aa", "> quote1", "> quote2", "bb"},
			want: "r{p{aa} blockquote{quote1  quote2 } p{bb}}"},
		{s: []string{"aa", "> quote1", "", "> quote2", "bb"},
			want: "r{p{aa} blockquote{quote1 } blockquote{quote2 } p{bb}}"},
		{s: []string{"aa`cc`bb"}, want: "r{p{aa<code>cc</code>bb}}"},
		{s: []string{"aa", "```", "a1", "a2", "```", "bb"}, want: "r{p{aa} pre{code{a1 a2}} p{bb}}"},

		// Lists
		{s: []string{"aa", "* 1", "* 2", "  + 2.1", "  + 2.2", "    - 2.2.1", "  + 2.3", "* 3", "cc"},
			want: "r{p{aa} ul{li{1} li{2} ul{li{2.1} li{2.2} ul{li{2.2.1}} li{2.3}} li{3}} p{cc}}"},
		{s: []string{"aa", "  * 1", "   l1", "  * 2", "   l2", "c"},
			want: "r{p{aa} ul{li{1 l1} li{2 l2}} p{c}}"},
		{s: []string{"a", "", "- b", "", "  * c", "", "  d", "", "- e", "", "f"},
			want: "r{p{a} ul{li{b <p></p>} ul{li{c <p></p>}} p{d <p></p>} li{e <p></p>}} p{f}}"},
		{s: []string{"* 1", "* 2", "a", "* p", "* q"},
			want: "r{ul{li{1} li{2}} p{a} ul{li{p} li{q}}}"},
		{s: []string{"a", "* 1", "* 2", "  + a", "* 3", "", "  b", " c", "* 4", "d"},
			want: "r{p{a} ul{li{1} li{2} ul{li{a}} li{3 <p></p> b c} li{4}} p{d}}"},
		{s: []string{"aa", "1. 1", "2. 2", "   2. 2.1", "   2. 2.2", "cc"},
			want: "r{p{aa} ol{li{1} li{2} ol{li{2.1} li{2.2}}} p{cc}}"},
		{s: []string{"aa", "1. 1", "2. 2", "   - 2.1", "   - 2.2", "cc"},
			want: "r{p{aa} ol{li{1} li{2} ul{li{2.1} li{2.2}}} p{cc}}"},

		// Tables
		{s: []string{"s", "| A | B |", "| --- | --- |", "| a | b |", "", "e"},
			want: "r{p{s table:style=\"width: 100%\"{tr{th{A} th{B}} tr{td{a} td{b}}} e}}"},
		{s: []string{"s", "| A | B |", "|| --- |  | --- ||", "| a | b |", "", "e"},
			want: "r{p{s table:style=\"width: 100%\"{tr{th{A} th{B}} tr{td{a} td{b}}} e}}"},
		{s: []string{"s", "| A | B | C | D |", "| --- | :--- | ---: | :---: |", "| a | b | c | d |", "e"},
			want: "r{p{s table:style=\"width: 100%\"{tr{th{A} th:style=\"text-align: left\"{B} th:style=\"text-align: right\"{C} th:style=\"text-align: center\"{D}} tr{td{a} td:style=\"text-align: left\"{b} td:style=\"text-align: right\"{c} td:style=\"text-align: center\"{d}}} e}}"},
	}
	for _, tst := range tests {
		ht := NewHTMLTree("r")
		ht.br, _ = ht.br.AddBranch(-1, "p")

		for _, s := range tst.s {
			err := ht.Build(s)
			if err != nil {
				t.Fatalf("Build(%q) returns error: %s, should be nil", s, err)
			}
		}

		got := ht.root.String()
		if got != tst.want {
			t.Errorf("Build(%q)... generates:\n%q\nshould be:\n%q\n",
				tst.s[0], got, tst.want)
		}
	}
}

func TestImages(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "aa ![im](lnk) bb",
			want: "aa <img src=\"lnk\" alt=\"im\"/> bb"},
		{s: "![i1](l1)![i2](l2)",
			want: "<img src=\"l1\" alt=\"i1\"/><img src=\"l2\" alt=\"i2\"/>"},
	}

	for _, tst := range tests {
		got := Images(tst.s)
		if got != tst.want {
			t.Errorf("Images(%q) generates:\n%q\nshould be:\n%q\n", tst.s, got, tst.want)
		}
	}
}

func TestInlineCodes(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "aa `code` bb", want: "aa <code>code</code> bb"},
	}

	for _, tst := range tests {
		got := InlineCodes(tst.s)
		if got != tst.want {
			t.Errorf("InlineCodes(%q) generates:\n%q\nshould be:\n%q\n", tst.s, got, tst.want)
		}
	}
}

func TestLinks(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "aa [txt](link) bb", want: "aa <a href=\"link\">txt</a> bb"},
		{s: "[t1](l1)[t2](l2)", want: "<a href=\"l1\">t1</a><a href=\"l2\">t2</a>"},
	}

	for _, tst := range tests {
		got := Links(tst.s)
		if got != tst.want {
			t.Errorf("Links(%q) generates:\n%q\nshould be:\n%q\n", tst.s, got, tst.want)
		}
	}
}
