package main

import (
  "os"
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
  Env   string `toml:"env"`
}

type Request struct {
  RID string `json:"rid"`
  EID int64 `json:"event"`
  REF string `json:"referer"`
  ENV string `json:"environment"`
}

type ReqEvent struct {
  ID int64 `json:"id"`
  NAME string `json:"name"`
}

type ReqReport struct {
  ID int64 `json:"id"`
  ENV string `json:"environment"`
}

type ReqEditClient struct {
  ID int64 `json:"id"`
  NAME string `json:"name"`
}

type ReqEditGenre struct {
  ID int64 `json:"id"`
  NAME string `json:"name"`
}

var con string
var db *sql.DB

func init() {
  var conf Config
  _, err := toml.DecodeFile("./config.toml", &conf)
  if err != nil {
    log.Fatalln("[Error] can not open `config.toml`")
  }
  con = fmt.Sprintf("%s:%s@%s/%s", conf.User, conf.Pass, conf.Host, conf.DB)
  db, err = sql.Open("mysql", con)
  if err != nil {
    log.Fatalln(fmt.Sprintf("[Error] can not connect MySQL. DSN is %s", con))
  }
  if conf.Env == "production" {
    gin.SetMode(gin.ReleaseMode)
  }
}

func setAccessHeader(c *gin.Context) {
  c.Header("Access-Control-Allow-Origin", "*")
  c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
  c.Header("Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS")
}

func main() {
  router := gin.Default()
  f, err := os.OpenFile("./event_track.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
  if err != nil {
    panic(err)
  }
  gin.DefaultWriter = f
  router.Use(gin.Logger())
  router.POST("/", addRecord)
  router.POST("/event/add", addEvent)
  router.POST("/event/report", getReport)
  router.POST("/master/client/edit", editClient)
  router.POST("/master/genre/edit", editGenre)
  router.GET("/call/status/check.json", func(c *gin.Context){
    c.JSON(200,gin.H{
      "status":"200",
      "result":"OK",
    })
  })
  router.Run(":8080")
}

func addRecord(c *gin.Context) {
  setAccessHeader(c)
  // Getting Request Paramaters
  var req Request
  c.BindJSON(&req)
  // Set environment
  var table string
  if req.ENV == "production" {
    table = "pro_event"
  } else {
    table = "stg_event"
  }
  q := fmt.Sprintf("insert into %s (rid, event_id, referer, created_at) values(\"%s\", \"%d\", \"%s\", \"%s\")", table, req.RID, req.EID, req.REF, time.Now().Format("2006-01-02 15:04:05"))
  query, err := db.Prepare(q)
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status": "500",
      "error": "cant add record.",
    })
  }
  defer query.Close()
  result, err := query.Exec()
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status": "500",
      "error": "cant add record.",
    })
  }
  log.Println(fmt.Sprintf("[Info] success Record.detail : %s", result))
  c.JSON(200, gin.H{
    "status" : "200",
    "message" : "add record done.",
  })
}

func addEvent(c *gin.Context) {
  setAccessHeader(c)
  // Getting Request Paramaters
  var req ReqEvent
  c.BindJSON(&req)
  fmt.Println(fmt.Sprintf("paramater is %s", req))
  // create add event query
  var q string
  if req.ID == 0 {
    q = fmt.Sprintf("insert into events (name, created_at, updated_at) values(\"%s\", \"%s\", \"%s\")", req.NAME, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
  } else {
    q = fmt.Sprintf("insert into events (id, name, created_at, updated_at) values(%d, \"%s\", \"%s\", \"%s\")", req.ID, req.NAME, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
  }
  query, err := db.Prepare(q)
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status":"500",
      "error":"cant create query",
    })
  }
  defer query.Close()
  result, err := query.Exec()
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status":"500",
      "error":"cant create query",
    })
  }
  log.Println(fmt.Sprintf("[Info] success Event.detail : %s", result))
  c.JSON(200, gin.H{
    "status":"200",
    "message":"add event done.",
  })
}

func getReport(c *gin.Context) {
  setAccessHeader(c)
  // Getting Request Paramaters
  var req ReqReport
  c.BindJSON(&req)
  var table string
  if req.ENV == "production" {
    table = "pro_event"
  } else {
    table = "stg_event"
  }
  // create get event report query
  q := fmt.Sprintf("select date_format(%s.created_at, '%%m') as month, date_format(%s.created_at, '%%d') as day, events.name, count(%s.id) as id from events left join %s on events.id = %s.event_id where events.id = %d group by date_format(%s.created_at, '%%Y%%m%%d')", table, table, table, table, table, req.ID, table)
  log.Println(fmt.Sprintf("[Info] SQL %s", q))
  rows, err := db.Query(q)
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status":"500",
      "error":"cant create query",
    })
  }
  defer rows.Close()
  // create result slice.
  var name string
  var res [][3]string
  for rows.Next() {
    var month string
    var day string
    var id string
    if err := rows.Scan(&month, &day, &name, &id); err != nil {
      log.Println(err)
    }
    var row [3]string = [3]string{month, day, id}
    res = append(res, row)
  }
  if err := rows.Err(); err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status":"500",
      "error":"cant create query",
    })
  }
  log.Println(fmt.Sprintf("[Info] success Get.detail : %s", res))
  c.JSON(200, gin.H{
    "status":"200",
    "name":name,
    "value":res,
  })
}

func editClient(c *gin.Context) {
  setAccessHeader(c)
  var req ReqEditClient
  c.BindJSON(&req)
  var sql string
  if req.ID != 0 {
    sql = fmt.Sprintf("update client_master set name = \"%s\" where id = %d", req.NAME, req.ID)
  } else {
    sql = fmt.Sprintf("insert into client_master(name) values(\"%s\")", req.NAME)
  }
  query, err := db.Prepare(q)
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status":"500",
      "error":"cant create query",
    })
  }
  defer query.Close()
  result, err := query.Exec()
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status":"500",
      "error":"cant create query",
    })
  }
  log.Println(fmt.Sprintf("[Info] success Event.detail : %s", result))
  c.JSON(200, gin.H{
    "status":"200",
    "value":result,
  })
}

func editGenre(c *gin.Context) {
  setAccessHeader(c)
  var req ReqEditGenre
  c.BindJSON(&req)
  var sql string
  if req.ID != 0 {
    sql = fmt.Sprintf("update genre_master set name = \"%s\" where id = %d", req.NAME, req.ID)
  } else {
    sql = fmt.Sprintf("insert into genre_master(name) values(\"%s\")", req.NAME)
  }
  query, err := db.Prepare(q)
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status":"500",
      "error":"cant create query",
    })
  }
  defer query.Close()
  result, err := query.Exec()
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{
      "status":"500",
      "error":"cant create query",
    })
  }
  log.Println(fmt.Sprintf("[Info] success Event.detail : %s", result))
  c.JSON(200, gin.H{
    "status":"200",
    "value":result,
  })
}
