package db_access

import (
    `github.com/haishuokQAQ/golang-skeleton/constant`
)

type DBDataGenerator interface {
    GenerateDSN() string
    GetType() string
}

func NewGenerator(dbType string) DBDataGenerator {
    if dbType == constant.DBTypeMysql {
        return &MysqlDataGenerator{}
    } else {
        return &PostgresDataGenerator{}
    }
}

type MysqlDataGenerator struct {

}

func (gen *MysqlDataGenerator) GenerateDSN() string{
    return `dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?", username, password, ip, port, dbName)`
}

func (gen *MysqlDataGenerator) GetType() string{
    return constant.DBTypeMysql
}

type PostgresDataGenerator struct {

}

func (gen *PostgresDataGenerator) GenerateDSN() string{
    return `dsn := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=disable binary_parameters=yes", ip, port, username, dbname, password)`
}

func (gen *PostgresDataGenerator) GetType() string{
    return constant.DBTypePostgres
}