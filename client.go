package main

import (
	"log"

	"github.com/gorilla/websocket"
)

//map clients
type ClientList map[*Client]bool

//define client as a visitor in frontend
type Client struct {
	connection *websocket.Conn

	manager *Manager

	//egres to avoid concurrent writes
	egress chan []byte
}

//initialize new client
func NewClient(conn *websocket.Conn, manager *Manager) *Client{
	return &Client{
		connection: conn,
		manager: manager,
		egress: make(chan []byte),
	}
}

//handle reading
func (c *Client) readMessages(){
	defer func(){
		//Gracefully close connection once done
		c.manager.removeClient(c)
	}()
	//while true loop
	for{
		//read queue
		messageType,payload,err := c.connection.ReadMessage()

		if err != nil{
			//exception
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				log.Printf("erro lendo mensagem: %v", err)
			}
			break
		}
		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))

		//for testing only remove later
		for wsclient := range c.manager.clients {
			wsclient.egress <- payload
		}
	}
}

//handle writing

func (c *Client) writeMessages(){
	defer func(){
		//close if reading trigger
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Println("connection closed: ", err)
				}
				// Return to close the goroutine
				return
			}
			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
			}
			log.Println("sent message")
		}

	}
}
