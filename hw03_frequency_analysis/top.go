package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	ErrRegExCompilation = errors.New("failed to compile regex pattern")
)

type token struct {
	counter int
	value   string
}

type parser struct {
	regex  *regexp.Regexp
	data   map[string]int
	rating []token
	star   bool
}

func newParser(pattern string, star bool) (*parser, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRegExCompilation, err)
	}

	return &parser{
		regex: reg,
		data:  make(map[string]int),
		star:  star,
	}, nil
}

func (p *parser) parse(str string) {
	if str == "" {
		return
	}

	for _, token := range p.regex.FindAllString(str, -1) {
		if p.star {
			p.data[strings.ToLower(token)]++
		} else {
			p.data[token]++
		}
	}
}

func (p *parser) rate(n int) []string {
	if len(p.data) == 0 {
		return nil
	}

	p.rating = make([]token, 0, len(p.data))

	for value, counter := range p.data {
		p.rating = append(p.rating, token{counter, value})
	}

	sort.Slice(p.rating, func(i, j int) bool {
		if p.rating[i].counter == p.rating[j].counter {
			return p.rating[i].value < p.rating[j].value
		}
		return p.rating[i].counter > p.rating[j].counter
	})

	r := make([]string, 0, n)
	for i := 0; i < n; i++ {
		r = append(r, p.rating[i].value)
	}

	return r
}

func Top10(input string, isStar bool) ([]string, error) {
	var pattern string

	if isStar {
		pattern = `[\p{L}\d][\p{L}\d-]*[\p{L}\d]*`
	} else {
		pattern = `\S+`
	}

	p, err := newParser(pattern, isStar)
	if err != nil {
		return nil, err
	}

	p.parse(input)

	return p.rate(10), nil
}
