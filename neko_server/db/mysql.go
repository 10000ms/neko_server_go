package db

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "time"
)

type NekoDbSettingType = map[string]map[string]string

func CreateDbConnect(
    Username string,
    Password string,
    Network string,
    Server string,
    Port string,
    Database string,
) *sql.DB {
    conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", Username, Password, Network, Server, Port, Database)
    DB, err := sql.Open("mysql", conn)
    if err != nil {
        fmt.Println("connection to mysql failed:", err)
        panic(err)
    }

    DB.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超时的连接就close
    DB.SetMaxOpenConns(100)                  //设置最大连接数
    return DB
}
