package generator

import (
	"github.com/balabanovds/otus-golang/hw09_generator_of_validators/pkg/parser"
)

type data struct {
	Func     string
	Short    string
	Field    string
	Value    string
	IsOrigin bool
	Origin   string
}

func genLen(s parser.ParsedStruct) []data {
	return getData(s, parser.FString, parser.TagLen, "validateLen")
}

func genLenSlice(s parser.ParsedStruct) []data {
	return getData(s, parser.FSliceString, parser.TagLen, "validateLen")
}

func genRegex(s parser.ParsedStruct) []data {
	return getData(s, parser.FString, parser.TagRegexp, "validateRegex")
}

func genRegexSlice(s parser.ParsedStruct) []data {
	return getData(s, parser.FSliceString, parser.TagRegexp, "validateRegex")
}

func genInStr(s parser.ParsedStruct) []data {
	return getData(s, parser.FString, parser.TagInStr, "validateInString")
}
func genInStrSlice(s parser.ParsedStruct) []data {
	return getData(s, parser.FSliceString, parser.TagInStr, "validateInString")
}

func genInInt(s parser.ParsedStruct) []data {
	return getData(s, parser.FInt, parser.TagInInt, "validateInInt")
}
func genInIntSlice(s parser.ParsedStruct) []data {
	return getData(s, parser.FSliceInt, parser.TagInInt, "validateInInt")
}

func genMin(s parser.ParsedStruct) []data {
	return getData(s, parser.FInt, parser.TagMin, "validateMin")
}
func genMinSlice(s parser.ParsedStruct) []data {
	return getData(s, parser.FSliceInt, parser.TagMin, "validateMin")
}

func genMax(s parser.ParsedStruct) []data {
	return getData(s, parser.FInt, parser.TagMax, "validateMax")
}
func genMaxSlice(s parser.ParsedStruct) []data {
	return getData(s, parser.FSliceInt, parser.TagMax, "validateMax")
}

func getData(s parser.ParsedStruct, fieldType parser.FType, tagType parser.TagType, fn string) []data {
	var result []data
	for _, field := range s.Fields {
		for _, tag := range field.Tags {
			if tag.Type == tagType && field.Type.Key == fieldType {
				result = append(result, data{
					Func:     fn,
					Short:    s.Short,
					Field:    field.Name,
					Value:    tag.Value,
					IsOrigin: field.Type.Root == field.Type.Declared,
					Origin:   field.Type.Root,
				})
			}
		}
	}
	return result
}

func merge(s parser.ParsedStruct, funcs map[string][]interface{}) struct {
	Struct parser.ParsedStruct
	Funcs  map[string][]interface{}
} {
	return struct {
		Struct parser.ParsedStruct
		Funcs  map[string][]interface{}
	}{s, funcs}
} //nolint:gofumpt
