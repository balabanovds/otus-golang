package generator

import (
	p "github.com/balabanovds/otus-golang/hw09_generator_of_validators/pkg/parser"
)

type data struct {
	Func     string
	Short    string
	Field    string
	Value    string
	IsOrigin bool
	Origin   string
}

func genLen(s p.ParsedStruct) []data {
	return getData(s, p.FString, p.TagLen, "validateLen")
}

func genLenSlice(s p.ParsedStruct) []data {
	return getData(s, p.FSliceString, p.TagLen, "validateLen")
}

func genRegex(s p.ParsedStruct) []data {
	return getData(s, p.FString, p.TagRegexp, "validateRegex")
}

func genRegexSlice(s p.ParsedStruct) []data {
	return getData(s, p.FSliceString, p.TagRegexp, "validateRegex")
}

func genInStr(s p.ParsedStruct) []data {
	return getData(s, p.FString, p.TagInStr, "validateInString")
}
func genInStrSlice(s p.ParsedStruct) []data {
	return getData(s, p.FSliceString, p.TagInStr, "validateInString")
}

func genInInt(s p.ParsedStruct) []data {
	return getData(s, p.FInt, p.TagInInt, "validateInInt")
}
func genInIntSlice(s p.ParsedStruct) []data {
	return getData(s, p.FSliceInt, p.TagInInt, "validateInInt")
}

func genMin(s p.ParsedStruct) []data {
	return getData(s, p.FInt, p.TagMin, "validateMin")
}
func genMinSlice(s p.ParsedStruct) []data {
	return getData(s, p.FSliceInt, p.TagMin, "validateMin")
}

func genMax(s p.ParsedStruct) []data {
	return getData(s, p.FInt, p.TagMax, "validateMax")
}
func genMaxSlice(s p.ParsedStruct) []data {
	return getData(s, p.FSliceInt, p.TagMax, "validateMax")
}

func getData(s p.ParsedStruct, fieldType p.FType, tagType p.TagType, fn string) []data {
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

func merge(s p.ParsedStruct, funcs map[string][]interface{}) struct {
	Struct p.ParsedStruct
	Funcs  map[string][]interface{}
} {
	return struct {
		Struct p.ParsedStruct
		Funcs  map[string][]interface{}
	}{s, funcs}
} //nolint:gofumpt
