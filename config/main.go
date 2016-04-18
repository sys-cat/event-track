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
  URL   string `toml:"url"`
  PORT  string `toml:"port"`
}

type Environment struct {
  VAL   string `toml:"env"`
}

type Blowfish struct {
  SEC   string `toml:"secret"`
  SALT  string `toml:"salt"`
  PASS  string `toml:"pass"`
}

func getToml()(res *Config) {
  _, err := toml.DecodeFile("./setting.toml", &res)
  if err != nil {
    panic(err)
  }
  return res
}
