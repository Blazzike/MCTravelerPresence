package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"sync"
)

var addr = flag.String("addr", "play.mctraveler.eu:1337", "MCTraveler Rich Presence")

func webSocketConnect(uuid string, connectedHandler func(func()), jsonHandler func(map[string]interface{})) {
	var waitGroup sync.WaitGroup

	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	query := u.Query()
	query.Set("api-version", "1")
	u.RawQuery = query.Encode()

	webSocketClient, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	connectedHandler(func() {
		waitGroup.Done()
	})

	defer webSocketClient.Close()

	done := make(chan struct{})

	waitGroup.Add(1)
	go func() {
		defer close(done)

		for {
			_, message, err := webSocketClient.ReadMessage()
			if err != nil {
				log.Println(err)

				return
			}

			payload := map[string]interface{}{}
			err = json.Unmarshal(message, &payload)
			if err != nil {
				panic(err)
			}

			jsonHandler(payload)
		}
	}()

	uuidPayload := struct {
		PayloadType string `json:"type"`
		Uuid        string `json:"uuid"`
	}{
		PayloadType: "uuid",
		Uuid:        uuid,
	}

	payload, err := json.Marshal(uuidPayload)
	if err != nil {
		panic(err)
	}

	err = webSocketClient.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		log.Println("Error sending UUID payload:", err)
	}

	waitGroup.Wait()
}
