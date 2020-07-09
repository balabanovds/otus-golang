package generator

import (
	"io"
	"text/template"
)

const (
	tmplLayout = `{{- block "tmplHeader" .}}{{end}}
{{- $fns := .Funcs }}
{{- range .Structs}}
	{{- block "tmplValidateFunc" merge . $fns }}{{end}}
{{- end}}`

	tmplValidateFunc = `
{{define "tmplValidateFunc"}}
func ({{.Struct.Short}} {{.Struct.Name}}) Validate() ([]ValidationError, error) {
	var vErrors []ValidationError
{{ $s := .Struct}}
{{- range $tmpl, $funcs := .Funcs }}
	{{- range $i, $fn := $funcs }}
		{{- $data := call $fn $s }}
		{{- if $data }}
			{{- range $data }}
				{{- if eq $tmpl "tmplFuncForOne"}}
					{{- block  "tmplFuncForOne" . }}{{- end}}
				{{- else}}
					{{- block  "tmplFuncForSlice" . }}{{- end}}
				{{- end}}
				
			{{- end}}
		{{- end}}
	{{- end}}
{{- end}}

	return vErrors, nil
}
{{- end}}`

	tmplFuncForOne = `
{{- define "tmplFuncForOne"}}
	{
		vErr, err := {{.Func}}("{{.Short}}.{{.Field}}", {{.Short}}.{{.Field}}, "{{.Value}}")
		if err != nil {
			return nil, err
		}
		if vErr != nil {
			vErrors = append(vErrors, *vErr)
		}
	}
{{- end}}
`

	tmplFuncForSlice = `
{{- define "tmplFuncForSlice"}}
	{
		for _, v := range {{.Short}}.{{.Field}} {
			vErr, err := {{.Func}}("{{.Short}}.{{.Field}}", v, "{{.Value}}")
			if err != nil {
				return nil, err
			}
			if vErr != nil {
				vErrors = append(vErrors, *vErr)
				break
			}
		}
	}
{{- end}}
`
)

func generateTemplates(data templateData, w io.Writer) error {
	funcMap := template.FuncMap{
		"merge": merge,
	}

	t, err := template.New("master").Funcs(funcMap).Parse(tmplLayout)
	if err != nil {
		return err
	}

	partials := []string{
		tmplHeader,
		tmplValidateFunc,
		tmplFuncForOne,
		tmplFuncForSlice,
	}

	for _, str := range partials {
		t, err = t.Parse(str)
		if err != nil {
			return err
		}
	}

	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func genFuncsMap() map[string][]interface{} {
	return map[string][]interface{}{
		"tmplFuncForOne": {
			genLen,
			genRegex,
			genInStr,
			genMin,
			genMax,
			genInInt,
		},
		"tmplFuncForSlice": {
			genLenSlice,
			genRegexSlice,
			genInStrSlice,
			genMinSlice,
			genMaxSlice,
			genInIntSlice,
		},
	}
}
