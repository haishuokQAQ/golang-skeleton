package main

import (
    `bytes`
    `fmt`
    `github.com/haishuokQAQ/golang-skeleton/utils`
    `strings`
    `text/template`
    `time`
)

func main() {
    bytes, err := GenerateTemplate(tempStr, &struct {
        Imports []string `json:"imports"`
        DalObjects []struct{Name string `json:"name"`
        GoType string `json:"dal_objects"`
        Tag string `json:"tag"`}
    }{
        Imports: []string{
            `"context"`,
            `"errors"`,
            `"fmt"`,
            `"time"`,
            fmt.Sprintf("%s",`_ "github.com/jinzhu/gorm/dialects/postgres"`),
        },
        DalObjects: []struct {
            Name   string `json:"name"`
            GoType string `json:"dal_objects"`
            Tag    string `json:"tag"`
        }{
            {Name: "test", GoType: "string", Tag: fmt.Sprintf("`%s`",`json:"test", "custom":"true"`)},
        },
    }, nil)
    if err != nil {
        panic(err)
    }
   fmt.Println(string(bytes))
}


func GenerateTemplate(templateText string, templateData interface{}, params map[string]interface{}) ([]byte, error) {
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
    }).Parse(templateText)
    if err != nil {
        return nil, err
    }
    var buf bytes.Buffer
    if err := t.Execute(&buf, templateData); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

var tempStr = `package db_access

import (
     "github.com/jinzhu/gorm"
{{- range .Imports}}
    "{{.}}"
{{- end}}
)

var StdDBAccess *DBAccess

type DBAccess struct {
    DB *gorm.DB
{{- range .DalObjects}}
    {{CamelizeStr .Name true}} {{.GoType}} {{.Tag}}
{{- end}}
}

func (da *DBAccess) DB() *gorm.DB {
	return da.db_access
}

func (da *DBAccess) BeginTransaction() *gorm.DB {
	return da.db_access.Begin()
}

// ConnectDB is used to open db_access connection
func ConnectDB(ip string, port int, username string, password string, dbname string) error {

	if StdDBAccess != nil {
		StdDBAccess.db_access.Close()
	}

	dsn := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=disable binary_parameters=yes", ip, port, username, dbname, password)

	db_access, err := gorm.Open("postgres", dsn)
	if err != nil {
		return err
	}
	db_access.SingularTable(true)
	db_access.LogMode(false)
	db_access = db_access.Set("gorm:save_associations", false).Set("gorm:association_save_reference", false)
	db_access.DB().SetConnMaxLifetime(1 * time.Hour)
	StdDBAccess = &DBAccess{
		db_access:                       db_access,
		{{- range .DalObjects}}
            {{CamelizeStr .Name true}} &{{.GoType}}
        {{- end}}
	}
	return nil
}`