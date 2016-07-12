package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ONSdigital/aws-status/assets"
	"github.com/ONSdigital/aws-status/aws"
	"github.com/gorilla/pat"
)

var status *aws.Status

func main() {
	bindAddr := os.Getenv("BIND_ADDR")
	if len(bindAddr) == 0 {
		bindAddr = ":8080"
	}

	var err error
	status, err = aws.NewFromFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		c := time.Tick(30 * time.Second)
		for _ = range c {
			status.Update()
		}
	}()
	status.Update()

	p := pat.New()
	p.Get("/data", awsStatus)
	p.Get("/", func(w http.ResponseWriter, req *http.Request) {
		http.FileServer(assets.AssetFS()).ServeHTTP(w, req)
	})

	if err := http.ListenAndServe(bindAddr, p); err != nil {
		log.Fatal(err)
	}
}

func awsStatus(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(&status)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error:", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(200)
	w.Write(b)
}
