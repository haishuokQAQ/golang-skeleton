package db_access

import (
    `github.com/haishuokQAQ/golang-skeleton/database`
    `testing`
)

func TestDbAccessGenerator_Generate(t *testing.T) {
    gen := &DbAccessGenerator{}
    err := gen.Generate(&DBAccessConfig{
        DalPath:    "./tmp/dal",
        DBType:     constant.DBTypePostgres,
        ORM:        "gorm",
        ORMVersion: "v2",
        Tables: []*database.TableMetaData{
            &database.TableMetaData{
                Name:    "test_table",
                Columns: nil,
                Indexes: nil,
            },
        },
    })
    if err != nil {
        panic(err)
    }
}
