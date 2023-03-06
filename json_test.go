package json_utils

import "testing"

func TestSanitizeJson(t *testing.T) {
	cases := []struct { in, want string 
	}{
		// validity 
		{ "[\"a\", \"b\", \"c\"]", "[\"a\", \"b\", \"c\"]" },
		// remove comments
		{ "//comment\n{\"a\":\"b\"}", "\n{\"a\":\"b\"}" },
		{ "/*//comment*/{\"a\":\"b\"}", "{\"a\":\"b\"}" },
		{ "{\"a\":\"b\"//comment\n}", "{\"a\":\"b\"\n}" },
		{ "{\"a\":\"b\"/*comment*/}", "{\"a\":\"b\"}" },
		{ "{\"a\"/*\n\n\ncomment\r\n*/:\"b\"}", "{\"a\":\"b\"}" },
		{ "/*!\n * comment\n */\n{\"a\":\"b\"}", "\n{\"a\":\"b\"}" },
		{ "{/*comment*/\"a\":\"b\"}", "{\"a\":\"b\"}" },
		// doesn\'t strip comments inside strings
		{ "{\"a\":\"b//c\"}", "{\"a\":\"b//c\"}" },
		{ "{\"a\":\"b/*c*/\"}", "{\"a\":\"b/*c*/\"}" },
		{ "{\"/*a\":\"b\"}", "{\"/*a\":\"b\"}" },
		{ "{\"\\\"/*a\":\"b\"}", "{\"\\\"/*a\":\"b\"}" },
		// escaped slashes
		{ "{\"\\\\\":\"https://foobar.com\"}", "{\"\\\\\":\"https://foobar.com\"}" },
		{ "{\"foo\\\"\":\"https://foobar.com\"}", "{\"foo\\\"\":\"https://foobar.com\"}" },
		// line endings - no comments
		{ "{\"a\":\"b\"\n}", "{\"a\":\"b\"\n}" },
		{ "{\"a\":\"b\"\r\n}", "{\"a\":\"b\"\r\n}" },
		// line endings - single line comment
		{ "{\"a\":\"b\"//c\n}", "{\"a\":\"b\"\n}" },
		{ "{\"a\":\"b\"//c\r\n}", "{\"a\":\"b\"\r\n}" },
		{ "{\"a\":\"b\"//c\r}", "{\"a\":\"b\"\r}" },
		// line endings - single line block comment
		{ "{\"a\":\"b\"/*c*/\n}", "{\"a\":\"b\"\n}" },
		{ "{\"a\":\"b\"/*c*/\r\n}", "{\"a\":\"b\"\r\n}" },
		{ "{\"a\":\"b\"/*c*/\r}", "{\"a\":\"b\"\r}" },
		// line endings - multi line block comment
		{ "{\"a\":\"b\",/*c\nc2*/\"x\":\"y\"\n}", "{\"a\":\"b\",\"x\":\"y\"\n}" },
		{ "{\"a\":\"b\",/*c\r\nc2*/\"x\":\"y\"\r\n}", "{\"a\":\"b\",\"x\":\"y\"\r\n}" },
		// strips trailing commas
		{ "{\"x\":true,}", "{\"x\":true}" },
		{ "{\"x\":true,\n  }", "{\"x\":true\n  }" },
		{ "[true, false,]", "[true, false]" },
		{ "{\n  \"array\": [\n    true,\n    false,\n  ],\n}", "{\n  \"array\": [\n    true,\n    false\n  ]\n}" },
		{ "{\n  \"array\": [\n    true,\n    false /* comment */ ,\n /*comment*/ ],\n}", "{\n  \"array\": [\n    true,\n    false  \n  ]\n}" },
	}
	for _, c := range cases {
		got := string(stripComments([]byte(c.in)))
		if got != c.want {
			t.Errorf("stripComments(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}