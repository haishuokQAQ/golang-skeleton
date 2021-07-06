package database

type DBConnector interface {
    ConnectionDB(host string, port int, dbName, user,password string) error
    AllTableData() (TableMetaDataList, error)
    SpecifiedTables(tableName []string) (TableMetaDataList, error)
    GetTableColumns(tableName string) (ColumnMetaDataList, error)
}

