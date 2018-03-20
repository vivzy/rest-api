package models

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "gopkg.in/gorp.v1"
  "log"
)

func initDb() *gorp.DbMap {
  db, err := sql.Open("mysql", "gouser:vishnu@/instructions")
  CheckErr(err, "sql.Open failed")

  dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
  dbmap.AddTableWithName(User{}, "user").SetKeys(true, "Id")

  err = dbmap.CreateTablesIfNotExists()
  CheckErr(err, "Create table failed")
  return dbmap
}

func CheckErr(err error, msg string) {
  if err != nil {
    log.Fatalln(msg, err)
  }
}

func LogErr(err error, msg string) {
  if err != nil {
    log.Printf(msg, err)
  }
}
