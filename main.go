package main

import (
	"bytes"
	"fmt"
	"github.com/FactomProject/web"
	"github.com/FactomProject/ptnet-eventstore/wsapi"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const httpApi = "http://127.0.0.1:8080"
const listenInterface = "0.0.0.0:8080"

// post a new event
func doPost() {
	actions := []string{"INC_0", "DEC_0"}
	//action := Actions[random.RandInt()%len(Actions)]
	action := actions[0]
	amount := "2"

	url := httpApi + "/dispatch/counter/foo/" + action + "/" + amount

	var jsonStr = []byte(`{"Hello":"World"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	_ = body
	fmt.Println("post response Body:", string(body))
}

// get test stream
func doGet(uri string) {
	url := httpApi + uri
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		//fmt.Println("response Status:", resp.Status)
		//fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(response.Body)
		_ = body
		fmt.Println("get response Body:", string(body))
	}
}

// test using periodic web requests
func tickerTest() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		//doGet("/stream/counter/foo")
		//doGet("/state/counter/foo")
		//doGet("/count/counter/foo")
		//doGet("/machine/counter")
		doGet("/contract/machine/counter")
		doPost()
	}
}


// start api + test ticker
// this tests only state machine and eventstore operations
func main() {
	go tickerTest()
	wsapi.AddRoutes()
	web.Run(listenInterface)
}
