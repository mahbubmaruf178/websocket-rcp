package main

import (
	"encoding/json"
	"log"
	"net/http"
	"websockrpc/action"
	"websockrpc/rpc"

	"github.com/gorilla/websocket"
)

func main() {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections by default
		},
	}
	fs := http.FileServer(http.Dir(".")) // Serve static files
	http.Handle("/", fs)                 // entry point or return index.html

	//  rpc handler(action) for websocket
	handler := rpc.NewRpchandler()
	handler.AddMethod("add", action.Add) // Example: Adding a method named "add"

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			var req rpc.RPCRequest // requset type
			if err := json.Unmarshal(message, &req); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error": "Invalid request"}`))
				continue
			}
			//  handle rpc that extract from  all action and return result
			result, err := handler.HandleRPC(req.Method, req.Params)
			response := rpc.RPCResponse{
				Result: result,
			}
			if err != nil {
				response.Error = err.Error()
			}

			responseBytes, err := json.Marshal(response)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error": "Internal error"}`))
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error": "Internal error"}`))
				break
			}
		}
	})
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
