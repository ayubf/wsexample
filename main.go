package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	port           = "8080"
	templateFolder = "public/"
	webTitle       = "SimpleAPI"
)

func loadTemplate(fName string, args ...interface{}) string {
	body, err := os.ReadFile(templateFolder + fName)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(string(body), args...)
}

func auxMessage(r string) string {
	t := time.Now()
	return fmt.Sprintf("-> %v:%v @ %v ", port, r, t.Format(time.UnixDate))
}

func rootRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println(auxMessage("/"))
	io.WriteString(w, loadTemplate("index.html", webTitle))
}

func createRoomRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		var userName struct {
			Name string
		}

		fmt.Println(string(bod))
		err = json.Unmarshal(bod, &userName)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]interface{}{
			"created": true,
			"code":    "1234ABCD",
		})

	} else {
		io.WriteString(w, loadTemplate("erorr.html", "Forbidden Method", "Method Not Allowed On This Route"))
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsReader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	wsReader(ws)

}

func main() {
	p := 8080
	t := time.Now()
	http.HandleFunc("/", rootRoute)
	http.HandleFunc("/createRoom", createRoomRoute)
	http.HandleFunc("/joinRoom", wsEndpoint)
	// http.HandleFunc("/", joinRoomRoute)

	fmt.Printf("-> Server Running on Port: %v\n-> %v\n", p, t.Format(time.UnixDate))
	err := http.ListenAndServe(":"+port, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Printf("Error while starting server: %v\n", err)
		os.Exit(1)
	}
}
