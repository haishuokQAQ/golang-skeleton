package db_access

import (
    `fmt`
    `github.com/haishuokQAQ/golang-skeleton/database`
    `testing`
)

func TestObjectGenerator_Generate(t *testing.T) {
    generator := &ObjectGenerator{}
    code, meta, err := generator.Generate(&database.TableMetaData{
        Name:    "test_table",
        Columns: nil,
        Indexes: nil,
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(*meta)
    fmt.Println(code)
}
