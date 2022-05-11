package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

type ConnectionDetails struct {
	RpcUser    string
	RpcPass    string
	RpcHost    string
	RpcPort    string
	RemoteHost string
	P2PPort    string
}

type Page struct {
	Title             string
	ConnectionDetails *ConnectionDetails
	BlockchainInfo    *GetBlockchainInfoResponse
}

func main() {

	details := &ConnectionDetails{
		RpcUser:    os.Getenv("RPC_USER"),
		RpcPass:    os.Getenv("RPC_PASS"),
		RpcHost:    os.Getenv("RPC_HOST"),
		RpcPort:    os.Getenv("RPC_PORT"),
		P2PPort:    os.Getenv("P2P_PORT"),
		RemoteHost: os.Getenv("REMOTE_HOST"),
	}

	client, err := NewClient(details.RpcHost, details.RpcPort, details.RpcUser, details.RpcPass, false, 30)
	if err != nil {
		log.Fatalf("jsonrpc: %v", err)
	}
	service := &Elements{client}
	response, err := service.getBlockchainInfo()
	if err != nil {
		log.Fatalf("jsonrpc: %v", err)
	}

	tmpl := template.Must(template.ParseFiles("layout.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Page{
			Title:             "Elements Core",
			ConnectionDetails: details,
			BlockchainInfo:    response,
		}
		tmpl.Execute(w, data)
	})
	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("error starting http server: %w", err)
	}
}
