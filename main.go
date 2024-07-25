package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Ammowar/goTimeSecret/handlers"
	"github.com/Ammowar/goTimeSecret/notifications"
	"github.com/Ammowar/goTimeSecret/storages"
)

var templates *template.Template
var secret_storage *storages.MemoryStorage
var ch *handlers.CreateHandler
var ih *handlers.IndexHandler
var vh *handlers.ViewHandler
var wsh *handlers.WebSocketHandler

var wsConnManager *notifications.WsConnectionManager

func start() {
	wsConnManager = notifications.NewWsConnectionManager()
	templates = template.Must(template.ParseGlob("templates/*.html"))
	secret_storage = storages.NewMemoryStorage()
	ih = handlers.NewIndexHandler(templates)
	ch = handlers.NewCreateHandler(secret_storage, templates)
	vh = handlers.NewViewHandler(secret_storage, templates, wsConnManager)
	wsh = handlers.NewWebSocketHandler(wsConnManager)
}

func main() {
	start()
	http.HandleFunc("/", ih.Handler)
	http.HandleFunc("/create", ch.Handler)
	http.HandleFunc("/secret/", vh.Handler)
	http.HandleFunc("/ws", wsh.Handler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
