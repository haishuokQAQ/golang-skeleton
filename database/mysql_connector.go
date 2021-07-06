package database

import (
    "database/sql"
    "errors"
    "fmt"
    `github.com/haishuokQAQ/golang-skeleton/constant`
    "strings"

    _ "github.com/go-sql-driver/mysql"
)

const (
    mysqlTableNamesSql          = `select table_name from information_schema.tables where table_schema = ? and table_type = 'BASE TABLE';`
    mysqlSpecifiedTableNamesSql = `select table_name from information_schema.tables where table_schema = ? and table_name in ('%s') and table_type = 'BASE TABLE';`
    mysqlTableColumnsSql        = `select column_name,
is_nullable, if(column_type = 'tinyint(1)', 'boolean', data_type),
column_type like '%unsigned%'
from information_schema.columns
where table_schema = ? and  table_name = ?
order by ordinal_position;
`
)

type MysqlConnector struct {
    db     *sql.DB
    dbName string
}

func (m *MysqlConnector) ConnectionDB(host string, port int, dbName, username,password string) error {
    dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?", username, password, host, port, dbName)
    fmt.Println("MySQL Connecting dsn : " + dsn)
    m.dbName = dbName
    db, err := sql.Open(constant.DBTypeMysql, dsn)
    if err != nil {
        return err
    }
    if err := db.Ping(); err != nil {
        return err
    }
    m.db = db
    return nil
}

func (m *MysqlConnector) AllTableData() (TableMetaDataList, error) {
    rows, err := m.db.Query(mysqlTableNamesSql, m.dbName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    rev := TableMetaDataList{}
    for rows.Next() {
        var tableName string
        if err := rows.Scan(&tableName); err != nil {
            return nil, err
        }
        tableColumnsInfo, err := m.GetTableColumns(tableName)
        if err != nil {
            return nil, err
        }
        rev = append(rev, &TableMetaData{Name: tableName, Columns: tableColumnsInfo})
    }

    return rev, rows.Err()
}

func (m *MysqlConnector) SpecifiedTables(tableNameList []string) (TableMetaDataList, error) {
    if len(tableNameList) == 0 {
        return nil, errors.New("tableNameList is empty")
    }
    sqlStr := fmt.Sprintf(mysqlSpecifiedTableNamesSql, strings.Join(tableNameList, "','"))
    rows, err := m.db.Query(sqlStr, m.dbName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    rev := TableMetaDataList{}
    for rows.Next() {
        var tableName string
        if err := rows.Scan(&tableName); err != nil {
            return nil, err
        }
        tableColumnsInfo, err := m.GetTableColumns(tableName)
        if err != nil {
            return nil, err
        }
        rev = append(rev, &TableMetaData{Name: tableName, Columns: tableColumnsInfo})
    }

    return rev, rows.Err()
}

func (m *MysqlConnector) GetTableColumns(tableName string) (ColumnMetaDataList, error) {
    rows, err := m.db.Query(mysqlTableColumnsSql, m.dbName, tableName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    rev := ColumnMetaDataList{}
    for rows.Next() {
        var name, isNullable, dataType string
        var isUnsigned bool
        if err := rows.Scan(&name, &isNullable, &dataType, &isUnsigned); err != nil {
            return nil, err
        }
        rev = append(rev, NewColumnMetaData(name,
            strings.ToLower(isNullable) == "yes", dataType, isUnsigned, tableName))
    }
    return rev, rows.Err()
}
