package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	mongoDBHosts = "149.28.195.244"
	authDatabase = "admin"
	authUserName = "xiaoyu"
	authPassword = "Jianyuhao123$"
)

func connectDB() *mgo.Database {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoDBHosts},
		Timeout:  60 * time.Second,
		Database: authDatabase,
		Username: authUserName,
		Password: authPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	mongoSession.SetMode(mgo.Monotonic, true)
	db := mongoSession.DB("proxypool")
	return db
}

func findMulit() {
	db := connectDB()

	collection := db.C("new")

	// 使用All方法，一次性消耗较多内存，如果数据较多，可以考虑使用迭代器
	type NewProxy struct {
		IP       string `bson:"ip"`
		Port     string `bson:"port"`
		Country  string `bson:"country"`
		Protocol string `bson:"protocol"`
	}

	var newProxy NewProxy

	iter := collection.Find(nil).Iter()

	for iter.Next(&newProxy) {
		/*
			Loop through each new proxy:
		*/
		fmt.Println(newProxy)
	}
}

func main() {
	findMulit()
}
