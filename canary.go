package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
)

func unixtime() int32 {
	return int32(time.Now().Unix())
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", unixtime())
}

// TODO: terminate stream after disconnet / fixed time
func stream(w http.ResponseWriter, r *http.Request) {
	flusher, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	for {
		fmt.Fprintf(w, "data: %v\n\n", unixtime())
		flusher.Flush()
		time.Sleep(time.Second)
	}
}

func checker() {
	fmt.Println("starting checker ...")
	for {
		resp, _ := http.Get("http://localhost:8000")
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(body)
		time.Sleep(time.Second)		
	}
}

func main() {
	go checker()
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/stream", stream)
	http.ListenAndServe(":8000", nil)
}
