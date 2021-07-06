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

type ModelGenerator struct {

}

func (gen *ModelGenerator) GenerateModel(meta *database.TableMetaData, modelBasePath string) error{
    params := map[string]string{
        "packageName": "db",
    }
    t, err := template.New("tableTemplate").Funcs(template.FuncMap{
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
        "param": func(name string) interface{} {
            if v, ok := params[name]; ok {
                return v
            }
            return ""
        },
    }).Parse(tableModelTemplate)
    if err != nil {
        return err
    }
    var buf bytes.Buffer
    if err := t.Execute(&buf, meta); err != nil {
        return  err
    }
    if err := utils.SaveFile(modelBasePath, meta.Name+".go", buf.Bytes()); err != nil {
        fmt.Printf("save file error: %#v", err.Error())
        return err
    }
    return nil
}

var tableModelTemplate = `
package {{param "packageName"}}
{{ if .Imports }}
import (
{{- range .Imports}}
	"{{.}}"
{{- end}}
)
{{end}}

{{$structName := CamelizeStr .Name true}}

type {{$structName}} struct {
{{- range .Columns}}
	{{CamelizeStr .Name true}} {{.GoType}} ` + "{{.Tag}}" + `
{{- end}}
}
{{$firstChar := FirstCharacter .Name}}

func ({{$firstChar}} {{$structName }}) TableName() string {
	return "{{.Name}}"
}`