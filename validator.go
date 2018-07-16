package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

var wg sync.WaitGroup

// NewProxy : data tructure for holding the proxy document in MongoDB
type NewProxy struct {
	IP       string `bson:"ip"`
	Port     string `bson:"port"`
	Country  string `bson:"country"`
	Protocol string `bson:"protocol"`
}

const (
	mongoDBHosts = "149.28.195.244"
	authDatabase = "admin"
	authUserName = "xiaoyu"
	authPassword = "Jianyuhao123$"
)

func connectDB() *mgo.Session {
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

	return mongoSession
}

func displayNumGoroutines() {
	for {
		fmt.Println("Current # goroutines: ", runtime.NumGoroutine())
		time.Sleep(time.Second * 1)
	}
}

func moveToHotStorage(newProxy NewProxy, db *mgo.Database) {
	// Move the given proxy document from New collection to Hot collection
	defer wg.Done()

	hotC := db.C("hot")
	newC := db.C("new")

	// Insert the proxy into Hot collection
	if err := hotC.Insert(newProxy); err != nil {
		panic(err)
	}

	// Remove the document from New collection
	if err := newC.Remove(newProxy); err != nil {
		panic(err)
	}
}

func moveToColdStorage(newProxy NewProxy, db *mgo.Database) {
	// Move the given proxy document from New collection to Cold collection
	defer wg.Done()

	coldC := db.C("cold")
	newC := db.C("new")

	// Insert the proxy into Cold collection
	if err := coldC.Insert(newProxy); err != nil {
		panic(err)
	}

	// Remove the document from New collection
	if err := newC.Remove(newProxy); err != nil {
		panic(err)
	}
}

func checkAvailability(newProxy NewProxy, timeoutSeconds time.Duration, db *mgo.Database) {
	// Check the connectibility of the given proxy document
	defer wg.Done()

	// Use the given proxy to make a request to the given server
	req, err := http.NewRequest("GET", "https://api.ipify.org/", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	proxyURL, err := url.Parse(strings.ToLower(newProxy.Protocol) + "://" + newProxy.IP + ":" + newProxy.Port)
	if err != nil {
		// Likely to be malformed URL
		fmt.Println(err)
		return
	}

	// Timeout settings for different stages of the connection
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds * time.Second),
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Second*timeoutSeconds)
			},
			TLSHandshakeTimeout:   timeoutSeconds * time.Second,
			ResponseHeaderTimeout: timeoutSeconds * time.Second,
			IdleConnTimeout:       timeoutSeconds * time.Second,
			ExpectContinueTimeout: timeoutSeconds * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			Proxy:                 http.ProxyURL(proxyURL)}}

	res, err := client.Do(req)
	if err != nil {
		// Proxy is somewhat not available, move to cold storage
		fmt.Println(err)
		wg.Add(1)
		fmt.Println("Cold Storage:", newProxy.IP)
		go moveToColdStorage(newProxy, db)
		return
	}

	defer res.Body.Close()

	// body, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))
	fmt.Println("Okay", newProxy.IP)
}

func useProxies(db *mgo.Database) {
	/*
		Iterate through all of the new proxies and use them to send request to the Go server.
		If unavailable, remove them from the new collection.
	*/
	newC := db.C("new")

	var newProxy NewProxy

	iter := newC.Find(nil).Iter()

	i := 1
	for iter.Next(&newProxy) {
		for runtime.NumGoroutine() > 60 {
			// Limit the number of concurrent sockets connections
			time.Sleep(time.Second * 1)
		}
		fmt.Println("#", i)
		wg.Add(1)
		go checkAvailability(newProxy, 4, db)
		i++
	}
}

func main() {
	// Start the server for profiling
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	mongoSession := connectDB()
	defer mongoSession.Close()
	db := mongoSession.DB("proxypool")
	useProxies(db)
	wg.Wait()
}
