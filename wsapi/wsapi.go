package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/contracts"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/web"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const httpApi = "http://127.0.0.1:8080"
const listenInterface = "0.0.0.0:8080"

type EventList struct {
	Count uint64         `json:"count"`
	Items []*ptnet.Event `json:"items"`
}

type EventCount struct {
	Count uint64 `json:"count"`
}

func getEventList(schema string, oid string, count_only bool) *EventList {
	txn := ptnet.Txn(schema, false)
	events, _ := txn.Get(ptnet.EventTable, "Oid", oid)
	result := new(EventList)
	result.Count = 0

	for x := events.Next(); x != nil; x = events.Next() {
		result.Count++
		if false == count_only {
			result.Items = append(result.Items, x.(*ptnet.Event))
		}
	}
	return result
}

// publish a new event
func dispatch(ctx *web.Context, schema string, oid string, action string, value string) {
	amount, _ := strconv.ParseInt(value, 10, 64)
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	event, _ := ptnet.Commit(schema, oid, action, uint64(amount), body)
	data, _ := json.Marshal(event)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.Write(data)
}

// query state
func state(ctx *web.Context, schema string, Oid string) {
	txn := ptnet.Txn(schema, false)
	raw, _ := txn.First(ptnet.StateTable, "id", Oid)

	ctx.Header().Set("Content-Type", "application/json")
	if raw != nil {
		data, _ := json.Marshal(raw.(ptnet.State))
		ctx.ResponseWriter.Write(data)
	}
}

// get state machine definition
func machine(ctx *web.Context, schema string) {
	data, _ := json.Marshal(ptnet.StateMachines[schema])
	ctx.ResponseWriter.Write(data)
}

// number of events in memory
func count(ctx *web.Context, schema string, oid string) {
	result := getEventList(schema, oid, true)
	data, _ := json.Marshal(EventCount{result.Count})
	ctx.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.Write(data)
}

// return full event stream
func stream(ctx *web.Context, schema string, oid string) {
	result := getEventList(schema, oid, false)
	data, _ := json.Marshal(result)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.Write(data)
}

// get contract definition
func contractMachine(ctx *web.Context, schema string) {
	data, _ := json.Marshal(contracts.Contracts[schema])
	ctx.ResponseWriter.Write(data)
}

// get contract definition
func contractState(ctx *web.Context, schema string) {
	// FIXME
	// ctx.ResponseWriter.Write(data)
}

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

// configure web api
func addRoutes() {
	web.Post("/dispatch/(.+)/(.+)/(.+)/(.+)", dispatch)
	web.Get("/state/(.+)/(.+)", state)
	web.Get("/machine/(.+)", machine)
	web.Get("/stream/(.+)/(.+)", stream)
	web.Get("/count/(.+)/(.+)", count)
	web.Get("/contract/machine/(.+)", contractMachine)
	//web.Get("/contract/state/(.+)/(.+)", contractState)
}

func Run() {
	addRoutes()
	web.Run(listenInterface)
}

// start polling test server
func main() {
	go tickerTest()
	Run()
}
