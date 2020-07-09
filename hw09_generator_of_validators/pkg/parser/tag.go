package parser

import (
	"regexp"
	"strings"
)

type TagType int

const (
	// TagLen tag for string: field must be exact length
	TagLen TagType = iota
	// TagRegExp tag for string: field must match regexp value
	TagRegexp
	// TagInStr for string: field must be in variety of values
	TagInStr
	// TagInInt for int: field must be in variety of values
	TagInInt
	// TagMin for int: field must be not less than value
	TagMin
	// TagMax for int: field must be not more than value
	TagMax
)

type Tag struct {
	Type  TagType
	Value string
}

func newTag(key, val string, fieldType FType) (Tag, bool) {
	var t TagType

	switch key {
	case "len":
		t = TagLen
	case "regexp":
		t = TagRegexp
	case "in":
		if fieldType == FString || fieldType == FSliceString {
			t = TagInStr
		} else if fieldType == FInt || fieldType == FSliceInt {
			t = TagInInt
		}
	case "min":
		t = TagMin
	case "max":
		t = TagMax
	default:
		return Tag{}, false
	}

	return Tag{
		Type:  t,
		Value: val,
	}, true
}

func parseTags(str, tagToken string, fieldType FType) []Tag {
	tRegWhole := regexp.MustCompile(`.*` + tagToken + `:"(\S+?)".*`)
	tRegOne := regexp.MustCompile(`(\w+):(\S+)`)

	var tags []Tag
	wholeTag := tRegWhole.FindStringSubmatch(str)
	if len(wholeTag) < 1 {
		return tags
	}

	for _, str := range strings.Split(wholeTag[1], "|") {
		matches := tRegOne.FindStringSubmatch(str)

		t, ok := newTag(matches[1], matches[2], fieldType)
		if !ok {
			continue
		}

		tags = append(tags, t)
	}

	return tags
}

func inTags(tags []Tag, tType ...TagType) bool {
	for _, tag := range tags {
		for _, t := range tType {
			if tag.Type == t {
				return true
			}
		}
	}
	return false
}
