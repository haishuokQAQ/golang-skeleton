package procedures

import (
    `github.com/haishuokQAQ/golang-skeleton/config`
    `testing`
)

func TestExecuteForDatabase(t *testing.T) {
   err := ExecuteForDatabase(&config.DatabaseConfig{
        DBType:          "mysql",
        DBName:          "block_chain",
        Host:            "192.168.1.101",
        Port:            3306,
        UserName:        "root",
        Password:        "Khs19940718!",
        SpecifiedTables: nil,
        ORM:             "gorm",
        ORMVersion:      "v1",
    },"./temp")
   if err != nil {
       panic(err)
   }
}

func TestExecuteForMysql(t *testing.T) {
    err := ExecuteForDatabase(&config.DatabaseConfig{
        DBType:          "mysql",
        DBName:          "ttoss",
        Host:            "10.184.24.223",
        Port:            3306,
        UserName:        "mysql",
        Password:        "P@ssw0rd!",
        SpecifiedTables: nil,
        ORM:             "gorm",
        ORMVersion:      "v1",
    },"./temp")
    if err != nil {
        panic(err)
    }
}