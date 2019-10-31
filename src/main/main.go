package main

import (
	"fmt"
	"log"
	"utility"
)

var DBString = utility.Conf.String("database", "connectString", "MySQL connectString")
var RedisString = utility.Conf.String("database", "redisAddress", "Redis Address")

func main() {
	err := utility.Conf.Parse()
	if err != nil {
		log.Fatal(err)
		return
	}
	utility.RedisAddr = *RedisString

	fmt.Println("Ping!")

}
