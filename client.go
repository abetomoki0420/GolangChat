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
	send chan SendData
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
			//ルームへの参加
			c.name = user.Name
			c.room_id = user.Room_id
			entrance.enter <- c
		}else if user.Type == "leave" {
			//ルームから離脱
			c.room.unregister <- c
		}else if user.Type == "post" {
			var sendData SendData
			sendData.Type = "post"
			sendData.Message = user.Body
			var name []string
			name = append(name , c.name )
			sendData.Users = name
			c.room.broadcast <- sendData
		}

	}
}

func (c *Client) writePump(){
	defer func(){
		c.conn.Close()
	}()

	for{
		select{
		case send_data , ok := <-c.send:
			if !ok{
				c.conn.WriteMessage( websocket.CloseMessage , []byte{})
				return
			}
			w,  err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			message , err := json.Marshal( send_data )
			if err != nil{
				log.Fatalf("MarshalError err:%v " , err )
				return
			}
			w.Write(message)

			w.Close()
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024 ,
	WriteBufferSize: 1024 ,
}
