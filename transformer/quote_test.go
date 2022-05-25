package transformer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/transform"
)

func replace(s string) string {
	return strings.ReplaceAll(s, "a", "x")
}

func TestProcessLiteral(t *testing.T) {
	tr := ProcessLiteral(replace)
	for _, c := range []struct {
		in, exp string
	}{
		{
			in:  "a = \"abc\" and b = 'abc'",
			exp: "x = \"abc\" xnd b = 'abc'",
		},
		{
			in:  "a = \"'abc'\" and b = 'abc'",
			exp: "x = \"'abc'\" xnd b = 'abc'",
		},
		{
			in:  "a = '\"abc\"' and b = 'abc'",
			exp: "x = '\"abc\"' xnd b = 'abc'",
		},
		{
			in:  "a = 'abc\"' and b = 'abc'",
			exp: "x = 'abc\"' xnd b = 'abc'",
		},
		{
			in:  "a = '\"abc' and b = 'abc'",
			exp: "x = '\"abc' xnd b = 'abc'",
		},
		{
			in:  "a = \"abc'\" and b = 'abc'",
			exp: "x = \"abc'\" xnd b = 'abc'",
		},
		{
			in:  "a = \"'abc\" and b = 'abc'",
			exp: "x = \"'abc\" xnd b = 'abc'",
		},
		{
			in:  "a = \"\\\"abc\" and b = 'abc'",
			exp: "x = \"\\\"abc\" xnd b = 'abc'",
		},
		{
			in:  "a = '\\'abc' and b = 'abc'",
			exp: "x = '\\'abc' xnd b = 'abc'",
		},
		{
			in:  "a = '\\\\'abc' and b = 'abc'",
			exp: "x = '\\\\'xbc' and b = 'xbc'",
		},
	} {
		out, _, _ := transform.String(tr, c.in)
		assert.Equal(t, c.exp, out)
		//t.Logf("Result:\n\tin:\t%s\n\tout:\t%s\n", c.in, out)
	}
}
