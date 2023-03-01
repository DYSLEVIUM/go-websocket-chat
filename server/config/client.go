package config

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager

	// egress is a channel, used to avoid concurrent writes on the websocket connection
	egress chan []byte
}

// factory for newclient
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func() {
		// cleanup connection
		c.manager.removeClient(c)
	}()
	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error Reading message: %v\n", err)
			}
			break
		}

		for wsclient := range c.manager.clients {
			// if c == wsclient {
			// 	continue
			// }
			wsclient.egress <- payload
		}

		log.Println("MessageType: ", messageType)
		log.Println("Payload:", string(payload))
	}
}

// we are making a channel to write the message as we don't want to write a lot of messages at the same time as the connection can't handle that
func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	/*
		The infinite for loop in the writeMessages method is used to continuously listen for messages on the egress channel and write them to the WebSocket connection.

		The select statement inside the for loop waits for either a message to be received on the egress channel or for the loop to be interrupted (e.g. by a call to return). If a message is received on the egress channel, the message is written to the WebSocket connection using the WriteMessage method. If the egress channel is closed (i.e. !ok), then the connection is closed and the method returns.

		Without the for loop, the method would only be able to write a single message to the WebSocket connection and then return, even if there are more messages waiting to be sent on the egress channel. The loop ensures that the method continues to write messages as long as the egress channel is open and there are messages to be sent.
	*/
	for {
		/*
			select statement:

			Used to wait on multiple channel operations simultaneously and execute the code block associated with the first channel that is ready to communicate.
			Can only be used with channel operations (send, receive, and close).
			Can have an optional default case which is executed if none of the channels are ready.

		*/
		select {
		case message, ok := <-c.egress:
			if !ok {
				// check if the channel is still up, but if it is closed, we send close connection to the client
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// if we have problem sending, we probably have closed the connection
					log.Println("Connection Closed")
				}
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			log.Println("Message Sent")
		}
	}
}
