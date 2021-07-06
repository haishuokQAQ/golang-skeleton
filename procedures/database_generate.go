package procedures

import (
    `fmt`
    `github.com/haishuokQAQ/golang-skeleton/config`
    `github.com/haishuokQAQ/golang-skeleton/constant`
    `github.com/haishuokQAQ/golang-skeleton/database`
    `github.com/haishuokQAQ/golang-skeleton/generators/db_access`
    `github.com/haishuokQAQ/golang-skeleton/utils`
    `strings`
)

func ExecuteForDatabase(dbConf *config.DatabaseConfig, projectBasePath string) error {
    var connector database.DBConnector
    if dbConf.DBType == constant.DBTypeMysql {
        connector = &database.MysqlConnector{}
    } else {
        connector = &database.PostgresConnector{}
    }
    err := connector.ConnectionDB(dbConf.Host, dbConf.Port,dbConf.DBName, dbConf.UserName, dbConf.Password)
    if err != nil {
        return err
    }
    var metas []*database.TableMetaData
    if len(dbConf.SpecifiedTables) > 0{
        metas, err = connector.SpecifiedTables(dbConf.SpecifiedTables)
        if err != nil {
            return err
        }
    } else {
        metas, err = connector.AllTableData()
        if err != nil {
            return err
        }
    }
    // 生成model
    modelGen := &db_access.ModelGenerator{}
    modelBasePath := fmt.Sprintf("%s/model/db/", strings.TrimSuffix(projectBasePath, "/"))
    if len(metas) > 0 {
        // 创建文件夹
        err = utils.MkdirPathIfNotExist(modelBasePath)
        if err != nil {
            return err
        }
        // 清理文件夹
       /* err = utils.CleanUpGenFiles(modelBasePath)
        if err != nil {
            return err
        }*/
    }
    for _, meta := range metas {
        err = modelGen.GenerateModel(meta, modelBasePath)
        if err != nil {
            return err
        }
    }
    // 生成dao
    generator := &db_access.DbAccessGenerator{}
    err = generator.Generate(&db_access.DBAccessConfig{
        DalPath:    fmt.Sprintf("%s/dal",strings.TrimSuffix(projectBasePath, "/")),
        DBType:     dbConf.DBType,
        ORM:        dbConf.ORM,
        ORMVersion: dbConf.ORMVersion,
        Tables:     metas,
    })
    if err != nil {
        return err
    }
    return nil
}
