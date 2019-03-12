package main

import(
	"log"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Client struct{
	room *Room
	room_id int
	conn *websocket.Conn
	send chan []byte
	name string
}

func ( c *Client) readPump(){
	defer func(){
		c.room.unregister <- c
		c.conn.Close()
	}()
	for{
		_ , mes , err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError( err , websocket.CloseGoingAway , websocket.CloseAbnormalClosure){
				log.Printf("error: %v" , err )
			}
			break
		}

		var user User
		if err = json.Unmarshal( mes , &user); err != nil {
			log.Fatalf("UnmarshalError :%v" ,err)
			return
		}

		log.Println(user.Type)
		if user.Type == "enter" {
			c.name = user.Name
			c.room_id = user.Room_id
			entrance.enter <- c
		}else if user.Type == "leave" {
			c.room.unregister <- c
		}else if user.Type == "post" {
			c.room.broadcast <- []byte(user.Body)
		}else{
			c.room.broadcast <- mes
		}

	}
}

func (c *Client) writePump(){
	defer func(){
		c.conn.Close()
	}()

	for{
		select{
		case message , ok := <-c.send:
			if !ok{
				c.conn.WriteMessage( websocket.CloseMessage , []byte{})
				return
			}
			w,  err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i:= 0 ; i < n ; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			w.Close()
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024 ,
	WriteBufferSize: 1024 ,
}
