package utility

import "github.com/c4pt0r/ini"

var Conf = ini.NewConf("config.ini")

var RedisAddr string
var DBAddr string
var ServerPort string
var Neo4jAddr string
var Neo4jUser string
var Neo4jPassword string
