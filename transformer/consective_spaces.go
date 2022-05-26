package transformer

import (
	"golang.org/x/text/transform"
)

type removeConsecutiveSpaces struct{}

var _ transform.Transformer = (*removeConsecutiveSpaces)(nil)

func (t *removeConsecutiveSpaces) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	nSrc = len(src)
	nDstBuf := len(dst)
	pre := byte(' ')
	for i := 0; i < nSrc; i++ {
		if src[i] == ' ' && pre == ' ' {
			continue
		}
		if nDst >= nDstBuf {
			return nDst, i, transform.ErrShortDst
		}
		dst[nDst] = src[i]
		pre = dst[nDst]
		nDst++
	}
	if atEOF && dst[nDst-1] == ' ' {
		nDst--
	}
	return
}

func (t *removeConsecutiveSpaces) Reset() {}

// RemoveConsecutiveSpaces removes all consecutive spaces.
// This handles only U+0020 for performance. You should normalize spaces in advance.
func RemoveConsecutiveSpaces() transform.Transformer {
	return &removeConsecutiveSpaces{}
}
