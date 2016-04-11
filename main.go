package main

import (
  "log"
  "fmt"
  "time"
  "github.com/gin-gonic/gin"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/BurntSushi/toml"
)

type Config struct {
  User  string `toml:"user"`
  Port  string `toml:"port"`
  Host  string `toml:"host"`
  Pass  string `toml:"pass"`
  DB    string `toml:"db"`
}

type Request struct {
  RID string `json:"rid"`
  EID int64 `json:"event"`
  REF string `json:"referer"`
  ENV string `json:"environment"`
}

func setAccessHeader(c *gin.Context) {
  c.Header("Access-Control-Allow-Origin", "*")
  c.Header("Access-Control-Allow-Credentials", "true")
  c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
  c.Header("Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS")
}

func main() {
  router := gin.Default()
  router.POST("/", addEvent)
  router.Run(":8080")
}

func addEvent(c *gin.Context) {
  setAccessHeader(c)
  var req Request
  c.BindJSON(&req)
  var conf Config
  _, err := toml.DecodeFile("./config.toml", &conf)
  if err != nil {
    log.Fatalln(err)
  }
  log.Println("open db")
  var con string
  con = conf.User+":"+conf.Pass+"@"+conf.Host+"/"+conf.DB
  log.Println(con)
  db, err := sql.Open("mysql", con)
  if err != nil {
    log.Fatalln(err)
  }
  var table string
  if req.ENV == "production" {
    table = "pro_event"
  } else {
    table = "stg_event"
  }
  log.Println("create query")
  query, err := db.Prepare(fmt.Sprintf("insert into %s (rid, event_id, referer, created_at) values(\"%s\", \"%d\", \"%s\", \"%s\")", table, req.RID, req.EID, req.REF, time.Now().Format("2006-01-02 15:04:05")))
  if err != nil {
    log.Fatalln(err)
  }
  defer query.Close()
  result, err := query.Exec()
  log.Println("result query")
  if err != nil {
    log.Fatalln(err)
  }
  log.Println("valiable query")
  log.Println(result)
  db.Close()

  c.JSON(200, gin.H{
    "status" : "OK",
  })
}
