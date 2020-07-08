package parser

import (
	"regexp"
	"strings"
)

type tagType int

const (
	tLen    tagType = iota // string: exact length of string
	tRegexp                // string: regexp
	tIn                    // string, int: value must be in variety
	tMin                   // int: not less than
	tMax                   // int: not more than
)

type tag struct {
	tType tagType
	value string
}

func newTag(key, val string) (tag, bool) {
	var t tagType

	switch key {
	case "len":
		t = tLen
	case "regexp":
		t = tRegexp
	case "in":
		t = tIn
	case "min":
		t = tMin
	case "max":
		t = tMax
	default:
		return tag{}, false
	}

	return tag{
		tType: t,
		value: val,
	}, true
}

func parseTags(str, tagToken string) []tag {
	tRegWhole := regexp.MustCompile(`.*` + tagToken + `:"(\S+?)".*`)
	tRegOne := regexp.MustCompile(`(\w+):(\S+)`)

	var tags []tag
	wholeTag := tRegWhole.FindStringSubmatch(str)
	if len(wholeTag) < 1 {
		return tags
	}

	for _, str := range strings.Split(wholeTag[1], "|") {
		matches := tRegOne.FindStringSubmatch(str)

		t, ok := newTag(matches[1], matches[2])
		if !ok {
			continue
		}

		tags = append(tags, t)
	}

	return tags
}
