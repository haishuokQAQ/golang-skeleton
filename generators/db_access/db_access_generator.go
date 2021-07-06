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

type DBAccessConfig struct {
    DalPath    string
    DBType     string
    ORM        string
    ORMVersion string
    Tables     []*database.TableMetaData
}

type DbAccessGenerator struct {
    objectGenerator ObjectGenerator
}

type DBAccessGenerateMeta struct {
    DalObjects []*ChildDBAccessMeta
    Imports    []string
    ORMConnectCode string
}



func (gen *DbAccessGenerator) Generate(config *DBAccessConfig) error {
    err := utils.MkdirPathIfNotExist(config.DalPath)
    if err != nil {
        return err
    }
    // 根据db类型生成generator
    generator := NewGenerator(config.DBType)
    // 生成meta
    meta := &DBAccessGenerateMeta{
        DalObjects:     []*ChildDBAccessMeta{},
        Imports:        []string{},
        ORMConnectCode: "",
    }
    code, ormImports, err := GenerateConnectCode(&GenerateConfig{DBGenerator: generator}, config.ORM, config.ORMVersion)
    if err != nil {
        return err
    }
    meta.ORMConnectCode = code
    meta.Imports = append(meta.Imports, ormImports...)
    childGen := &ObjectGenerator{}
    for _, table := range config.Tables {
        innerCode, childMeta, err := childGen.Generate(table)
        if err != nil {
            return err
        }
        meta.DalObjects = append(meta.DalObjects, childMeta)
        // 如果存在则删除
        path := fmt.Sprintf("%s/%s.go",strings.TrimSuffix(config.DalPath, "/"), childMeta.Name)
        err = utils.DeleteFileIfExist(path)
        if err != nil {
            return err
        }
        err = utils.SaveFile(strings.TrimSuffix(config.DalPath, "/"), fmt.Sprintf("%s.go", childMeta.Name), []byte(innerCode))
        if err != nil {
            return err
        }
    }
    temp, err := template.New("db_access").
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
        Parse(dbAccessTemplate)
    if err != nil {
        return err
    }

    buffer := &bytes.Buffer{}
    err = temp.Execute(buffer, meta)
    if err != nil {
        return err
    }
    path := fmt.Sprintf("%s/%s.go",strings.TrimSuffix(config.DalPath, "/"), "db_access")
    err = utils.DeleteFileIfExist(path)
    if err != nil {
        return err
    }
//    err = ioutil.WriteFile(path, []byte(str), fs.ModeDevice)
    //    err = utils.SaveToFile(strings.TrimSuffix(config.DalPath, "/"), fmt.Sprintf("%s.go", "db_access"), buffer.Bytes())
    err = utils.SaveFile(strings.TrimSuffix(config.DalPath, "/"), fmt.Sprintf("%s.go", "db_access"), buffer.Bytes())
    if err != nil {
        return err
    }
    return nil
}


var dbAccessTemplate = `package dal

import (
    "fmt"
{{- range .Imports}}
    {{.}}
{{- end}}
)

const TransactionContextKey = "transaction_object"

var StdDBAccess *DBAccess

type DBAccess struct {
    db *gorm.DB
{{- range .DalObjects}}
    {{CamelizeStr .Name true}} *{{CamelizeStr .GoType true}}
{{- end}}
}

func (da *DBAccess) DB() *gorm.DB {
	return da.db
}

func (da *DBAccess) BeginTransaction() *gorm.DB {
	return da.db.Begin()
}

func GetDB(ctx context.Context) *gorm.DB {
    var tx *gorm.DB
	txObj, exists := ctx.Value(TransactionContextKey).(*gorm.DB)
	if exists {
		tx = txObj
	}
	if tx != nil {
		// tx 在初始化的时候已经注入过 trace logger
		return tx
	}
	//_, exists = ctxUtil.GetTrace(ctx)
	//if exists {
		// 存在 trace 对象，注入trace logger
		// 克隆，不会克隆底层的 sqldb 对象，不用担心连接池问题
		//db := stdDBAccess.db.New()
		//db.SetLogger(gormCommon.NewGormLoggerWithLevel(stdDBAccess.Logger.WithTraceInCtx(ctx), log.LevelInfo))
		//return db
	//} else {
		// 没有 trace 对象，不处理
	//	return stdDBAccess.db
	//}
	return stdDBAccess.db
}

// ConnectDB is used to open db connection
func ConnectDB(ip string, port int, username string, password string, dbName string) error {

	if StdDBAccess != nil {
		StdDBAccess.db.Close()
	}
    
    {{.ORMConnectCode}}
	StdDBAccess = &DBAccess{
		db:                       db,
		{{- range .DalObjects}}
            {{CamelizeStr .Name true}} :&{{CamelizeStr .GoType true}}{},
        {{- end}}
	}
	return nil
}`
