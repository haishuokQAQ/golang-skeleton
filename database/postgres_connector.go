package database

import (
    `database/sql`
    `errors`
    `fmt`
    `github.com/haishuokQAQ/golang-skeleton/constant`
    `github.com/haishuokQAQ/golang-skeleton/utils`
    "github.com/lib/pq"
    `strings`
)

var (
    postgresTableNamesSql          = `select table_name from information_schema.tables where table_schema = 'public';`
    postgresSpecifiedTableNamesSql = `select table_name from information_schema.tables where table_schema = 'public' and table_name =any($1);`

    postgresTableColumnsSql = `select column_name, is_nullable, data_type, false
from information_schema.columns where table_schema = 'public' and table_name = $1 order by ordinal_position;`
)

type PostgresConnector struct {
    db     *sql.DB
    dbName string
}

func (p *PostgresConnector) ConnectionDB(host string, port int, dbName, username,password string) error {
    dsn := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=disable binary_parameters=yes", host, port, username, dbName, password)
    dbName, err := utils.GetDbNameFromDSN(dsn)
    if err != nil {
        return err
    }
    p.dbName = dbName
    fmt.Println("Postgres Connecting dsn : " + dsn)
    db, err := sql.Open(constant.DBTypePostgres, dsn)
    if err != nil {
        return err
    }
    if err = db.Ping(); err != nil {
        return err
    }
    p.db = db
    return nil
}

func (p *PostgresConnector) AllTableData() (TableMetaDataList, error) {
    rows, err := p.db.Query(postgresTableNamesSql)
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
        tableColumnsInfo, err := p.GetTableColumns(tableName)
        if err != nil {
            return nil, err
        }
        rev = append(rev, &TableMetaData{Name: tableName, Columns: tableColumnsInfo})
    }

    return rev, rows.Err()
}
func (p *PostgresConnector) SpecifiedTables(tableNameList []string) (TableMetaDataList, error) {
    if len(tableNameList) == 0 {
        return nil, errors.New("tableNameList is empty")
    }
    rows, err := p.db.Query(postgresSpecifiedTableNamesSql, pq.Array(tableNameList))
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
        tableColumnsInfo, err := p.GetTableColumns(tableName)
        if err != nil {
            return nil, err
        }
        rev = append(rev, &TableMetaData{Name: tableName, Columns: tableColumnsInfo})
    }

    return rev, rows.Err()
}

func (p *PostgresConnector) GetTableColumns(tableName string) (ColumnMetaDataList, error) {
    rows, err := p.db.Query(postgresTableColumnsSql, tableName)
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

