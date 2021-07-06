package db_access

import (
    `bytes`
    `fmt`
    `github.com/haishuokQAQ/golang-skeleton/database`
    `github.com/haishuokQAQ/golang-skeleton/utils`
    `strings`
    `text/template`
    `time`
)

type ObjectGenerator struct {
}

func (gen *ObjectGenerator) Generate(meta *database.TableMetaData) (code string, childMeta *ChildDBAccessMeta, err error) {
    temp, err := template.New(fmt.Sprintf("%s_db_access", meta.Name)).
        Funcs(template.FuncMap{
            "CamelizeStr":    utils.CamelizeStr,
            "FirstCharacter": utils.FirstCharacter,
            "Replace": func(old, new, src string) string {
                return strings.ReplaceAll(src, old, new)
            },
            "Add": func(a, b int) int {
                return a + b
            },
            "now": func() string {
                return time.Now().Format(time.RFC3339)
            },
        }).
        Parse(childDBAccessTemplate)
    if err != nil {
        return "", nil, err
    }
    childMeta = &ChildDBAccessMeta{
        Name:   meta.Name,
        GoType: fmt.Sprintf("%s_db_access", meta.Name),
    }
    buffer := &bytes.Buffer{}
    err = temp.Execute(buffer, childMeta)
    if err != nil {
        return "", nil, err
    }
    return buffer.String(), childMeta, nil
}

type ChildDBAccessMeta struct {
    Name      string   `json:"name"`
    GoType    string   `json:"go_type"`
    Imports   []string `json:"imports"`
    Functions string   `json:"functions"`
}

var childDBAccessTemplate = `
package dal
{{ if .Imports }}
import (
{{- range .Imports}}
	"{{.}}"
{{- end}}
)
{{end}}

{{$structName := CamelizeStr .GoType true}}

type {{$structName}} struct {

}

{{- range .Functions}}
	{{.}}
{{- end}}
`

type SelectParamConfig struct {
    Column        *database.ColumnMetaData
    ConditionType string
    UseParameter  bool
    Value         interface{}
}

type SelectConfig struct {
    Param       []*SelectParamConfig
    ReturnBatch bool
}

func generateSelectFunction() string {
    return ""
}

type selectFunctionParam struct {
    firstChar string
    structName string
    functionName string
    paramStr string
    returnType string
}

var selectFunctionTemplate = `
func ({{.firstChar}} *{{.structName }}) {{.functionName}}(ctx context.Context, {{.paramStr}}) ({{.returnType}}, error) {
    result := {{.resultInitCode}}
    err := GetDB(ctx).Model(&{{modelClass}}{}).
{{- range .Where}}
	Where("{{.column}} {{.operator}} {{.placeholder}}", {{.value}}).
{{- end}}
    {{.selectFuction}}({{.addressTaker}}result).Error
    if err != nil {
    {{- range .errorFilterCode}}
        {{.}}
    {{- end}}
        return nil, err
    }
	return result, nil
}
`