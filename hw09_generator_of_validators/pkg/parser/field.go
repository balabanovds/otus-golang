package parser

type Field struct {
	Name string
	Type *FieldType
	Tags []Tag
}

func newField(name string, t *FieldType, tags []Tag) *Field {
	return &Field{
		Name: name,
		Type: t,
		Tags: tags,
	}
}

func (f *Field) Validate() error {
	switch f.Type.Key {
	case FString, FSliceString:
		if inTags(f.Tags, TagMin, TagMax, TagInInt) {
			return ErrParseTagsNotValid
		}
	case FInt, FSliceInt:
		if inTags(f.Tags, TagLen, TagRegexp, TagInStr) {
			return ErrParseTagsNotValid
		}
	case FUnknown:
		return ErrParseFieldUnknownType
	}
	return nil
}
