package db_access

import (
    `errors`
    `fmt`
    `github.com/haishuokQAQ/golang-skeleton/constant`
)

type GenerateConfig struct {
    DBGenerator DBDataGenerator
}

func GenerateGormV1ConnectCode(conf *GenerateConfig) (code string, imports []string){
    imports = []string{}
    if conf.DBGenerator.GetType() == constant.DBTypeMysql {
        imports = append(imports, `"github.com/jinzhu/gorm"`, `_ "github.com/jinzhu/gorm/dialects/mysql"`)
    } else {
        imports = append(imports, `"github.com/jinzhu/gorm"`, `_ "github.com/jinzhu/gorm/dialects/postgres"`)
    }
   return fmt.Sprintf(gormV1ConnectTemplate, conf.DBGenerator.GenerateDSN(), conf.DBGenerator.GetType()), imports
}

func GenerateGormV2ConnectCode(conf *GenerateConfig) (code string, imports []string){
    imports = []string{}
    if conf.DBGenerator.GetType() == constant.DBTypeMysql {
        imports = append(imports, `"gorm.io/gorm/schema"`, `"gorm.io/driver/mysql"`,`_ "gorm.io/driver/mysql"`, `"gorm.io/gorm"`)
    } else {
        imports = append(imports, `"gorm.io/gorm/schema"`, `"gorm.io/driver/postgres"`,`_ "gorm.io/driver/postgres"`, `"gorm.io/gorm"`)
    }
    return fmt.Sprintf(gormV2ConnectTemplate, conf.DBGenerator.GenerateDSN(), conf.DBGenerator.GetType(), conf.DBGenerator.GetType()), imports
}

func GenerateConnectCodeForGorm(conf *GenerateConfig, ormVersion string) (code string, imports []string, err error) {
    if ormVersion == "v1" {
        code, imports =  GenerateGormV1ConnectCode(conf)
        return
    } else if ormVersion == "v2" {
        code, imports =  GenerateGormV2ConnectCode(conf)
        return
    }
    return "", nil, errors.New(fmt.Sprintf("不支持的gorm版本%+v", ormVersion))
}

func GenerateConnectCode(conf *GenerateConfig, orm, ormVersion string) (string, []string, error) {
    if orm == "gorm" {
        return GenerateConnectCodeForGorm(conf, ormVersion)
    } else {
        return "", nil, errors.New(fmt.Sprintf("不支持的orm框架%s", orm))
    }
}

var gormV1ConnectTemplate = `
    %s
    db, err := gorm.Open("%s", dsn)
    if err != nil {
        return err
    }
    db.SingularTable(true)
    db.LogMode(false)
    db = db.Set("gorm:save_associations", false).Set("gorm:association_save_reference", false)
    db.DB().SetConnMaxLifetime(1 * time.Hour)
`


var gormV2ConnectTemplate = `
    %s
    db, err := gorm.Open(%s.Dialector{
        Config: &%s.Config{
            DSN: dsn,
        },
    }, &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true,
        },
        FullSaveAssociations:                     false,
        Logger:                                   nil,
        PrepareStmt:                              false,
        DisableAutomaticPing:                     false,
        DisableForeignKeyConstraintWhenMigrating: false,
        DisableNestedTransaction:                 false,
        AllowGlobalUpdate:                        false,
        QueryFields:                              false,
        Plugins:                                  nil,
    })
    if err != nil {
        return err
    }
    db = db.Set("gorm:save_associations", false).Set("gorm:association_save_reference", false)
`