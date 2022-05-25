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
	// previous value
	//prev byte
}

var _ transform.Transformer = (*processLiteral)(nil)

func (t *processLiteral) Transform(d, s []byte, eof bool) (nd int, ns int, err error) {
	nsb := len(s)
	ndb := len(d)
	var cd, cs int
	//fmt.Printf("Transform\n\tdst:\t%s\n\tsrc:\t%s'\n\teof:\t%t\n\tns:\t%d\n\tnb:\t%d\n\tcursor d=%d:s=%d\n", d, s, eof, nsb, ndb, cd, cs)
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
	/*
		write := func(s string) (int, error) {
			n := len(s)
			if cd+n > ndb {
				fmt.Printf("ShortDst dst=%s\n\tnd=%d,ndb=%d,cd=%d,cs=%d,n=%d\n", d[:cd], nd, ndb, cd, cs, n)
				return cd, transform.ErrShortDst
			}
			c := copy(d[cd:cd+n], s[:n])
			fmt.Printf("Write %d=%d data=%s\n", c, n, s[:n])
			return cd + n, nil
		}
	*/
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
					/*
						cd, err = write("[")
						if err != nil {
							return
						}
						nd = cd
					*/
					//ns++ // skip first '
					// NOTE change mode after output
					t.mstr = '\''
					continue
				} else if t.mstr == '\'' {
					// end 'string'
					cd, ns, err = flush()
					if err != nil {
						return
					}
					nd = cd
					/*
						* cd, err = write("]")
						if err != nil {
							return
						}
					*/
					nd = cd
					//ns++ // skip last '
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
					/*
						cd, err = write("<")
						if err != nil {
							return
						}
						nd = cd
					*/
					//ns++ // skip first "
					t.mstr = '"'
					continue
				} else if t.mstr == '"' {
					// end "string"
					cd, ns, err = flush()
					if err != nil {
						return
					}
					nd = cd
					/*
						cd, err = write(">")
						if err != nil {
							return
						}
						nd = cd
					*/
					//ns++ // skip last "
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
		//d[cd] = s[i]
		//cd++

		//t.prev = s[i]
	}
	if eof {
		cs++
		cd, ns, err = flush()
		if err != nil {
			return
		}
	}
	//fmt.Printf("Return\n\tdst:\t%s\n\tsrc:\t%s'\n\teof:\t%t\n\tns:\t%d\n\tnb:\t%d\n\tcursor d=%d:s=%d\n", d, s, eof, nsb, ndb, cd, cs)
	return cd, cs, nil
}

func (t *processLiteral) Reset() {
	//fmt.Printf("\n\nReset============================================\n")
	t.mstr = 0
	t.mesc = false
	//t.prev = 0
}

func ProcessLiteral(fn ProcessFn) transform.Transformer {
	return &processLiteral{
		fn: fn,
	}
}
