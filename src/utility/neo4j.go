package utility

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
)

var Neo4jDriver neo4j.Driver

func InitNeo4j() (err error) {
	Neo4jDriver, err = neo4j.NewDriver(Neo4jAddr, neo4j.BasicAuth(Neo4jUser, Neo4jPassword, ""))
	if err == nil {
		log.Println("Neo4j connected.")
	}
	return
}

func DestroyNeo4j() (err error) {
	err = Neo4jDriver.Close()
	return
}

func ImportData(url string) (err error) {
	var (
		session neo4j.Session
		result  neo4j.Result
	)
	session, err = Neo4jDriver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()
	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		log.Println("Starting query.")
		result, err = transaction.Run("LOAD CSV FROM '"+url+"' AS row\n"+
			"MERGE (a:Movie {MovieID: row[0]})\n"+
			"MERGE (b:Movie {MovieID: row[1]})\n"+
			"MERGE (a)-[:Similar {Similarity: toFloat(row[2])}]-(b)", map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		log.Println("Result acquired.")
		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}
		return nil, result.Err()
	})
	return err
}
