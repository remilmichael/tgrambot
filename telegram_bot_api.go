package main

import(
	"net/http"
	"fmt"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"
)

type Data struct {
	Ok bool `json:"ok"`
	Result []map[string]interface{} `json:"result"`
}

type Error struct {
	Ok bool `json:"ok"`
	Description string `json:"description"`
}

var token string
var response Data
var errc Error

func main(){
	fmt.Println("\n==Telegram bot==\nStarted...")
	token = "640784566:AAEj_ak_EjGeOR_AbmrJcdj2MhhSWUJbeHk"
	for true{
		go receive()
		time.Sleep(5 * time.Second)
	}
}

func receive(){
	resp, err := http.PostForm("https://api.telegram.org/bot"+token+"/getUpdates", url.Values{"limit":{"1"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&response)
	if len(response.Result) > 0 {
		up_id := int(response.Result[0]["update_id"].(float64))
		textv := response.Result[0]["message"].(map[string]interface{})["text"]
		chat_id := strconv.Itoa(int(response.Result[0]["message"].(map[string]interface{})["chat"].(map[string]interface{})["id"].(float64)))
		_ , err := http.PostForm("https://api.telegram.org/bot"+token+"/getUpdates", url.Values{"offset":{strconv.Itoa(up_id+1)}})
		if err != nil {
			panic(err)
		}
		if textv == "/idme" {
			cresp, err := http.PostForm("https://icanhazip.com", url.Values{})
			if err != nil {
				panic(err)
			}
			defer cresp.Body.Close()
			body, err := ioutil.ReadAll(cresp.Body)
			if err != nil {
				panic(err)
			}
			ip_addr := string(body[:len(body)-1])
			go send(ip_addr, chat_id)
		}
	}
}

func send(ip_addr string, chat_id string){
	sresp, err := http.PostForm("https://api.telegram.org/bot"+token+"/sendMessage", url.Values{"chat_id":{chat_id}, "text":{ip_addr}})
	if err != nil {
		panic(err)
	}
	defer sresp.Body.Close()
	err = json.NewDecoder(sresp.Body).Decode(&errc)
	if err != nil {
		panic(err)
	}
	if errc.Ok != true {
		fmt.Println("Error")
	}
}
