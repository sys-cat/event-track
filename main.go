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

// Config this app
type Config struct {
  User  string `toml:"user"`
  Port  string `toml:"port"`
  Host  string `toml:"host"`
  Pass  string `toml:"pass"`
  DB    string `toml:"db"`
  Env   string `toml:"env"`
}

// Request Paramaters
type Request struct {
  RID string `json:"rid"`
  EID int64 `json:"event"`
  REF string `json:"referer"`
  GEN int64 `json: "genre_id"`
  CLI int64 `json: "client_id"`
  OTH string `json: "other"`
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

type MasterTable struct {
  Event string
  Client string
  Genre string
}

type ListTable struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
}

type GetReport struct {}

var Master MasterTable = MasterTable{"events", "client_master", "genre_master"}
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
  // for version 2.0
  v1 := router.Group("/v1")
  {
    v1.POST("/add", addRecord)
    v1.POST("/report/:all", getReport)
    v1.POST("/master/event", addEvent)
    v1.POST("/master/client", editClient)
    v1.POST("/master/genre", editGenre)
    v1.GET("/master/list/:type", listMaster)
  }
  router.Run(":8080")
}

func selectRecord(table string, wheres []string)(count int) {
  var where string
  for key, value := range wheres {
    if key == 0 {
      where = fmt.Sprintf("%s", value)
    } else {
      where = fmt.Sprintf("%s and %s", where, value)
    }
  }
  var sql string = fmt.Sprintf("select count(*) from %s where %s", table, where)
  rows, err := db.Query(sql)
  if err != nil {
    log.Println(fmt.Sprintf("[Error] can not select %s data."))
  }
  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(&count)
    if err != nil {
      log.Println("[Error] can not get select datas")
    }
  }
  return count
}

func addRecord(c *gin.Context) {
  setAccessHeader(c)
  // Getting Request Paramaters
  var req Request
  c.BindJSON(&req)
  // Set environment
  var table string = fmt.Sprintf("%d_event", req.EID)
  query_body := "insert into %s (rid, event_id, referer, client_id, genre_id, other, created_at) values(\"%s\", %d, \"%s\", %d, %d, \"%s\" ,\"%s\")"
  q := fmt.Sprintf(query_body, table, req.RID, req.EID, req.REF, req.CLI, req.GEN, req.OTH, time.Now().Format("2006-01-02 15:04:05"))
  query, err := db.Prepare(q)
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{"status": "500",})
  }
  defer query.Close()
  result, err := query.Exec()
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{"status": "500",})
  }
  log.Println(fmt.Sprintf("[Info] success Record.detail : %s", result))
  c.JSON(200, gin.H{"status" : "200",})
}

func addEvent(c *gin.Context) {
  setAccessHeader(c)
  // Getting Request Paramaters
  var req ReqEvent
  c.BindJSON(&req)
  var where []string
  if req.ID != 0 {
    where = append(where, fmt.Sprintf("id=%d", req.ID))
  }
  // create add event query
  var sql string
  if selectRecord(Master.Event, where) > 0 {
    sql = fmt.Sprintf("update events set name = \"%s\", created_at = \"%s\", updated_at = \"%s\" where id = %d", req.NAME, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), req.ID)
  } else {
    if req.ID != 0 {
      sql = fmt.Sprintf("insert into events (id, name, created_at, updated_at) values(%d, \"%s\", \"%s\", \"%s\")", req.ID, req.NAME, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
    } else {
      sql = fmt.Sprintf("insert into events (name, created_at, updated_at) values(\"%s\", \"%s\", \"%s\")", req.ID, req.NAME, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
    }
  }
  query1, err := db.Prepare(sql)
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{"status":"500",})
  }
  defer query1.Close()
  result, err := query1.Exec()
  if err != nil {
    log.Println(err)
    c.JSON(500, gin.H{"status":"500",})
  } else {
    log.Println(fmt.Sprintf("[Info] success Event.detail : %s", result))
    sql = fmt.Sprintf("select id from events where name=\"%s\" limit 1", req.NAME)
    query2, err := db.Query(sql)
    if err != nil {log.Println(err)}
    defer query2.Close()
    var id int64
    for query2.Next() {
      if err := query2.Scan(&id); err != nil {log.Println(fmt.Sprintf("[Error] can not select event id.detail: %s", err))}
    }
    sql = `CREATE TABLE IF NOT EXISTS %s (
      id int(11) unsigned NOT NULL AUTO_INCREMENT,
      rid varchar(256) NOT NULL DEFAULT '',
      event_id int(11) NOT NULL DEFAULT 0,
      referer varchar(1024) DEFAULT NULL,
      client_id int(11) NOT NULL DEFAULT 0,
      genre_id int(11) NOT NULL DEFAULT 0,
      other text,
      created_at datetime NOT NULL,
      PRIMARY KEY (id),
      INDEX(id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
    table := fmt.Sprintf("%d_event", id)
    sql = fmt.Sprintf(sql, table)
    query3, err := db.Prepare(sql)
    if err != nil {log.Println(fmt.Sprintf("[Error] can not prepare sql.detail: %s", err))}
    defer query3.Close()
    result, err := query3.Exec()
    log.Println(fmt.Sprintf("[Info] create event_id table.detail: %s", result))
    if err != nil {
      c.JSON(500, gin.H{"status":"500",})
    } else {
      c.JSON(200, gin.H{"status":"200",})
    }
  }
}

func getReport(c *gin.Context) {
  setAccessHeader(c)
  // Getting Request Paramaters
  var req ReqReport
  c.BindJSON(&req)
  var table string = fmt.Sprintf("%d_event", req.ID)
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
  // Set Request Paramaters
  var req ReqEditClient
  c.BindJSON(&req)
  var where []string
  if req.ID != 0 {
    where = append(where, fmt.Sprintf("id=%d", req.ID))
  }
  var sql string
  if selectRecord(Master.Client, where) > 0 {
    sql = fmt.Sprintf("update client_master set name = \"%s\" where id = %d", req.NAME, req.ID)
  } else {
    if req.ID != 0 {
      sql = fmt.Sprintf("insert into client_master(id, name) values(%d, \"%s\")", req.ID, req.NAME)
    } else {
      sql = fmt.Sprintf("insert into client_master(name) values(\"%s\")", req.NAME)
    }
  }
  query, err := db.Prepare(sql)
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
  })
}

func editGenre(c *gin.Context) {
  setAccessHeader(c)
  var req ReqEditGenre
  c.BindJSON(&req)
  var where []string
  if req.ID != 0 {
    where = append(where, fmt.Sprintf("id=%d", req.ID))
  }
  var sql string
  if selectRecord(Master.Genre, where) > 0 {
    sql = fmt.Sprintf("update genre_master set name = \"%s\" where id = %d", req.NAME, req.ID)
  } else {
    if req.ID != 0 {
      sql = fmt.Sprintf("insert into genre_master(id, name) values(%d, \"%s\")", req.ID, req.NAME)
    } else {
      sql = fmt.Sprintf("insert into genre_master(name) values(\"%s\")", req.NAME)
    }
  }
  query, err := db.Prepare(sql)
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
  })
}

func listMaster(c *gin.Context) {
  setAccessHeader(c)
  // Getting Request Paramaters
  get_type := c.Param("type")
  var sql string
  switch get_type {
  case Master.Event:
    sql = fmt.Sprintf("select id, name from %s", Master.Event)
  case Master.Client:
    sql = fmt.Sprintf("select id, name from %s", Master.Client)
  case Master.Genre:
    sql = fmt.Sprintf("select id, name from %s", Master.Genre)
  default:
    c.JSON(500, gin.H{
      "status":"500",
      "error":"Can not use table Type.",
    })
  }
  rows, err := db.Query(sql)
  if err != nil {
    log.Println(fmt.Sprintf("[Error] can not get data.detail: %s", err))
  }
  defer rows.Close()
  var res []ListTable
  for rows.Next() {
    var id int64
    var name string
    if err:= rows.Scan(&id, &name); err != nil {
      log.Println(fmt.Sprintf("[Error] can not scan query result.detail: %s", err))
    }
    row := ListTable{Id: id, Name: name}
    res = append(res, row)
  }
  if err := rows.Err();err != nil {
    log.Println(fmt.Sprintf("[Error] unknown query error.detail: %s", err))
  }
  c.JSON(200, gin.H{
    "status":"200",
    "value":res,
  })
}
