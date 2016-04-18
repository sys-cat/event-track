package elasticsearch

import (
  "gopkg.in/olivere/elastic.v3"
  "../config"
  "fmt"
)

type EventMaster struct {
  ID    int64
  NAME  string
  ENV   string
}
type ClientMaster struct {
  ID    int64
  NAME  string
}
type GenreMaster struct {
  ID    int64
  NAME  string
}
type Event struct {
  ID    int64
  RID   string
  EID   int64
  REF   string
  GEN   int64
  CLI   int64
  OTH   string
}

var now string = time.Now().Format("2006-01-02 15:04:05")
var con *config.Config
var client *elastic.Client
var url string

func init() {
  con = config.getToml()
  url = fmt.Sprintf("%s:%s", con.ES.URL, con.ES.PORT)
  client, err := elastic.NewClient(
    elastic.SetURL(url)
  )
  if err != nil {
    panic(err)
  }
}

func alive()(res bool) {
  res = true
  info, code, err := client.Ping(url).Do()
  if err != nil {
    res = false
  }
  return res
}
