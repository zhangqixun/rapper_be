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
var Neo4jAddress = utility.Conf.String("database", "neo4jAddress", "Neo4j Address")
var Neo4jUser = utility.Conf.String("database", "neo4jUser", "Neo4j User")
var Neo4jPassword = utility.Conf.String("database", "neo4jPassword", "Neo4j Password")

func main() {
	err := utility.Conf.Parse()
	if err != nil {
		log.Fatal(err)
	}
	utility.RedisAddr = *RedisString
	utility.DBAddr = *DBString
	utility.ServerPort = *ServerPortString
	utility.Neo4jAddr = *Neo4jAddress
	utility.Neo4jUser = *Neo4jUser
	utility.Neo4jPassword = *Neo4jPassword

	err = utility.InitNeo4j()
	if err != nil {
		log.Fatal(err)
	}

	//_, err = utility.QueryTopSimilarities("1")
	//if err != nil {
	//	log.Fatal(err)
	//}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/quickstart", controllers.QuickStart)
	rtr.HandleFunc("/register", controllers.UserRegister).Methods("POST")
	rtr.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	rtr.HandleFunc("/logout", controllers.UserLogout).Methods("POST")
	rtr.HandleFunc("/editor", controllers.UserEditor).Methods("POST")
	rtr.HandleFunc("/userquery", controllers.UserQuery).Methods("POST")
	rtr.HandleFunc("/passwdmodi", controllers.UserPasswordModification).Methods("POST")
	rtr.HandleFunc("/movie/id", controllers.MovieIDQuery).Methods("GET")
	rtr.HandleFunc("/movie/type", controllers.MovieTypeQuery).Methods("GET")
	rtr.HandleFunc("/movie/keyword", controllers.MovieKeywordQuery).Methods("GET")
	rtr.HandleFunc("/movie/similarity", controllers.MovieSimilarityQuery).Methods("GET")
	rtr.HandleFunc("/user/footprint/new", controllers.UserBrowse).Methods("POST")
	rtr.HandleFunc("/user/footprints", controllers.UserBrowseQuery).Methods("POST")
	rtr.HandleFunc("/user/footprints/ID", controllers.UserBrowseQuery).Methods("POST")
	http.Handle("/", rtr)

	log.Println("Server running.")
	err = http.ListenAndServe(utility.ServerPort, nil)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Ping!")
	err = utility.DestroyNeo4j()
	if err != nil {
		log.Fatal(err)
	}
}
