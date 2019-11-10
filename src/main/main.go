package main

import (
	"controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"utility"
)

var DBString = utility.Conf.String("database", "connectString", "MySQL connectString")
var RedisString = utility.Conf.String("database", "redisAddress", "Redis Address")
var ServerPortString = utility.Conf.String("server", "serverPort", "Server Port")

func main() {
	err := utility.Conf.Parse()
	if err != nil {
		log.Fatal(err)
		return
	}
	utility.RedisAddr = *RedisString
	utility.DBAddr = *DBString
	utility.ServerPort = *ServerPortString

	rtr := mux.NewRouter()
	rtr.HandleFunc("/quickstart", controllers.QuickStart)
	rtr.HandleFunc("/register", controllers.UserRegister).Methods("POST")
	rtr.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	rtr.HandleFunc("/logout", controllers.UserLogout).Methods("POST")
	rtr.HandleFunc("/editor", controllers.UserEditor).Methods("POST")
	rtr.HandleFunc("/userquery", controllers.UserQuery).Methods("POST")
	rtr.HandleFunc("/passwdmodi", controllers.UserPasswordModification).Methods("POST")
	rtr.HandleFunc("/movietypequery", controllers.MovieTypeQuery).Methods("GET")
	http.Handle("/", rtr)
	http.ListenAndServe(utility.ServerPort, nil)

	//fmt.Println("Ping!")
}
