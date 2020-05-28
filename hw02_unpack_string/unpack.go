package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

type token struct {
	prev      *token
	mult      int
	str       string
	char      rune
	protected bool
}

func newToken(r rune, prev *token) *token {
	return &token{
		prev: prev,
		mult: 1,
		char: r,
	}
}

func (t *token) eval() error {
	// check first char
	if t.prev == nil {
		if unicode.IsDigit(t.char) {
			return ErrInvalidString
		}
		if unicode.IsLetter(t.char) {
			t.str = string(t.char)
		}
		return nil
	}

	// set protected flag
	if t.prev.char == '\\' && !t.prev.protected {
		t.protected = true
	}

	// check digit
	if unicode.IsDigit(t.char) {
		// "aaa10b" should fail
		if !t.prev.protected && unicode.IsDigit(t.prev.char) {
			return ErrInvalidString
		}

		// `qwe\4\5` => `qwe45`
		if t.protected {
			t.str = string(t.char)
			return nil
		}
		t.prev.mult = int(t.char - '0')
		return nil
	}

	// check escaped letter
	if unicode.IsLetter(t.char) && t.protected {
		t.str = "\\" + string(t.char)
		return nil
	}

	// skip not protected escape char
	if t.char == '\\' && !t.protected {
		return nil
	}

	t.str = string(t.char)

	return nil
}

func (t *token) String() string {
	var b strings.Builder
	for i := 0; i < t.mult; i++ {
		b.WriteString(t.str)
	}
	return b.String()
}

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	runes := []rune(str)
	tokens := make([]*token, 0)

	for i, r := range runes {
		var prev *token = nil
		if i > 0 {
			prev = tokens[i-1]
		}
		tokens = append(tokens, newToken(r, prev))
	}

	var b strings.Builder
	for _, t := range tokens {
		if err := t.eval(); err != nil {
			return "", err
		}
	}

	for _, t := range tokens {
		b.WriteString(t.String())
	}

	return b.String(), nil
}
