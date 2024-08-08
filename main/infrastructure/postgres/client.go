package postgres

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
    "tradingACE/main/infrastructure/server"
)

func Init() *sql.DB {
    database := server.SystemConfig.Database
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
        database.Host, database.Port, database.User, database.Password, database.Name, database.SSLMode, database.TimeZone)
    client, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatalf("Could not create a connection:%v", err)
    }

    err = client.Ping()
    if err != nil {
        log.Fatalf("Could not connect to PostgreSQL:%v", err)
    }

    log.Println("Connected to PostgreSQL")
    return client
}
