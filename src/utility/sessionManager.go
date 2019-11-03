package utility

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"strconv"
)

var RedisConn redis.Conn

func establishRedis() {
	var err error
	RedisConn, err = redis.Dial("tcp", RedisAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func PreprocessXHR(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

func NewSession(uid int) string {
	establishRedis()
	u, _ := uuid.NewV4()
	fmt.Println(u.String() + " " + strconv.Itoa(uid))
	_, err := RedisConn.Do("SET", u.String(), strconv.Itoa(uid))
	if err != nil {
		log.Fatal(err)
	}
	err = RedisConn.Close()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

func CheckSession(token string) int {
	establishRedis()
	uid, err := redis.String(RedisConn.Do("GET", token))
	if err != nil {
		return -1
	}
	id, _ := strconv.Atoi(uid)
	err = RedisConn.Close()
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func DropSession(token string) int {
	establishRedis()
	uid, err := redis.String(RedisConn.Do("GET", token))
	if err != nil {
		return -1
	}
	id, err := strconv.Atoi(uid)
	if err != nil {
		return -1
	}
	_, err = RedisConn.Do("DEL", token)
	if err != nil {
		log.Fatal(err)
	}
	err = RedisConn.Close()
	if err != nil {
		log.Fatal(err)
	}
	return id
}
