package transformer

import (
	"golang.org/x/text/transform"
)

type ProcessFn = func(string) string

type processLiteral struct {
	fn ProcessFn

	// mode string literal
	mstr byte
	// mode escape string
	mesc bool
}

var _ transform.Transformer = (*processLiteral)(nil)

func (t *processLiteral) Transform(d, s []byte, eof bool) (nd int, ns int, err error) {
	nsb := len(s)
	ndb := len(d)
	var cd, cs int
	flush := func() (int, int, error) {
		if ns == cs {
			return cd, ns, nil
		}
		n := cs - ns
		if cd+n > ndb {
			return cd, ns, transform.ErrShortDst
		}
		_ = copy(d[cd:cd+n], s[ns:cs])
		return cd + n, cs, nil
	}
	flushT := func(fn func(string) string) (int, int, error) {
		if ns == cs {
			return cd, ns, nil
		}
		out := fn(string(s[ns:cs]))
		n := len(out)
		if cd+n > ndb {
			return cd, ns, transform.ErrShortDst
		}
		_ = copy(d[cd:cd+n], []byte(out))
		return cd + n, cs, nil
	}
	for i := 0; i < nsb; i++ {
		cs = i
		switch s[i] {
		case '\'':
			if !t.mesc {
				if t.mstr == 0 {
					// start 'string'
					cd, ns, err = flushT(t.fn)
					if err != nil {
						return
					}
					nd = cd
					t.mstr = '\''
					continue
				} else if t.mstr == '\'' {
					// end 'string'
					cd, ns, err = flush()
					if err != nil {
						return
					}
					nd = cd
					t.mstr = 0
					continue
				}
			}
		case '"':
			if !t.mesc {
				if t.mstr == 0 {
					// start "string"
					cd, ns, err = flushT(t.fn)
					if err != nil {
						return
					}
					nd = cd
					t.mstr = '"'
					continue
				} else if t.mstr == '"' {
					// end "string"
					cd, ns, err = flush()
					if err != nil {
						return
					}
					nd = cd
					t.mstr = 0
					continue
				}
			}
		}
		if s[i] == '\\' {
			if t.mesc {
				t.mesc = false // special case \\
			} else {
				t.mesc = true // start \escape
			}
		} else {
			t.mesc = false // end \escape
		}
	}
	if eof {
		cs++
		if t.mstr != 0 {
			cd, ns, err = flushT(t.fn)
		} else {
			cd, ns, err = flush()
		}
		if err != nil {
			return
		}
	}
	return cd, cs, nil
}

func (t *processLiteral) Reset() {
	t.mstr = 0
	t.mesc = false
}

func ProcessLiteral(fn ProcessFn) transform.Transformer {
	return &processLiteral{
		fn: fn,
	}
}
