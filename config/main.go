package config

import (
  "github.com/BurntSushi/toml"
)

type Config struct {
  ES    Elasticsearch
  ENV   Environment
  BLOW  Blowfish
}

type Elasticsearch struct {
  URL   string
  PORT  string
}

type Environment struct {
  VAL   string
}

type Blowfish struct {
  SEC   string
  SALT  string
  PASS  string
}

func getToml()(res *Config) {
  _, err := toml.DecodeFile("./setting.toml", &res)
  if err != nil {
    panic(err)
  }
  return res
}
