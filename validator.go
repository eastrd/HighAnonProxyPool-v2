package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
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

	// Insert the proxy into Hot collection
	if err := hotC.Insert(newProxy); err != nil {
		panic(err)
	}

	removeFromNewCollection(newProxy, db)
}

func removeFromNewCollection(newProxy NewProxy, db *mgo.Database) {
	wg.Add(1)
	defer wg.Done()

	newC := db.C("new")

	// Remove the document from New collection
	if err := newC.Remove(newProxy); err != nil {
		panic(err)
	}
}

func moveToColdStorage(newProxy NewProxy, db *mgo.Database) {
	// Move the given proxy document from New collection to Cold collection
	defer wg.Done()

	coldC := db.C("cold")

	// Insert the proxy into Cold collection
	if err := coldC.Insert(newProxy); err != nil {
		panic(err)
	}

	removeFromNewCollection(newProxy, db)
}

func checkAvailability(newProxy NewProxy, timeoutSeconds time.Duration, db *mgo.Database) {
	// Check the connectibility of the given proxy document
	defer wg.Done()

	// Use the given proxy to make a request to the given server
	req, err := http.NewRequest("GET", "http://149.28.181.216/checkproxy", nil)
	if err != nil {
		return
	}

	proxyURL, err := url.Parse(strings.ToLower(newProxy.Protocol) + "://" + newProxy.IP + ":" + newProxy.Port)
	if err != nil {
		// Likely to be malformed URL
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
		wg.Add(1)
		go moveToColdStorage(newProxy, db)
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if string(body) == newProxy.IP {
		// Proxy's real IP is hidden. Move to Hot storage
		wg.Add(1)
		go moveToHotStorage(newProxy, db)
	} else {
		// Proxy connection successful, but its IP is not high anonymous or its response is blocked by its ISP.
		// Hence delete this proxy from database.
		go removeFromNewCollection(newProxy, db)
	}
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
		for runtime.NumGoroutine() > 20 {
			// Limit the number of concurrent sockets connections
			time.Sleep(time.Second * 1)
		}
		fmt.Println("#", i)
		wg.Add(1)
		go checkAvailability(newProxy, 5, db)
		i++
	}
}

func main() {
	mongoSession := connectDB()
	defer mongoSession.Close()
	db := mongoSession.DB("proxypool")
	useProxies(db)
	wg.Wait()
}
