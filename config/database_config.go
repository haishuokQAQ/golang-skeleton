package config

type DatabaseConfig struct {
    DBType          string   `json:"db_type"`
    DBName          string   `json:"db_name"`
    Host            string   `json:"host"`
    Port            int      `json:"port"`
    UserName        string   `json:"user_name"`
    Password        string   `json:"password"`
    SpecifiedTables []string `json:"specified_tables"`
    ORM string `json:"orm"`
    ORMVersion string `json:"orm_version"`
}
