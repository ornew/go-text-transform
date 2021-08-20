package jpnorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/transform"
)

func TestRemoveConsecutiveSpaces(t *testing.T) {
	tr := RemoveConsecutiveSpaces()
	for _, c := range []struct {
		in, exp string
	}{
		{
			in:  "a  ",
			exp: "a",
		},
		{
			in:  "  a",
			exp: "a",
		},
		{
			in:  "a  b  ",
			exp: "a b",
		},
		{
			in:  "  a  b",
			exp: "a b",
		},
		{
			in:  "  a  ",
			exp: "a",
		},
		{
			in: "0 1  2   3    4     5      6       7        8         9          " +
				"10 11  12   13    14     15      16       17        18         19          " +
				"20 21  22   23    24     25      26       27        28         29          " +
				"30 31  32   33    34     35      36       37        38         39          ",
			exp: "0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20" +
				" 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39",
		},
	} {
		out, _, _ := transform.String(tr, c.in)
		assert.Equal(t, c.exp, out)
	}
}
