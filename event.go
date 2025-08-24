
package main

import "encoding/json"

// Event is the Messages sent over the websocket
type Event struct {
	// Type is the message type sent
	Type string `json:"type"`

	// Payload is the data Based on the Type aka message in this case
	Payload json.RawMessage `json:"payload"`
}

//deal with message type to sockets

type EventHandler func(event Event, c *Client) error

const(
	//evetn for message sent
	EventSendMessage= "send_message"
)

//payload in send_message event
type SendMessageEvent struct{
	Message string `json:"message"`
	From    string `json:"from"`
}
