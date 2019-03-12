package main

import(
	"log"
)

type Room struct{
	id int
	entrance *Entrance
	clients map[*Client]Client
	broadcast chan SendData
	register chan *Client
	unregister chan *Client
	title string
}

type RoomSummary struct{
	Id int `json:"id"`
	Title string `json:"title"`
	Count int `json:"count"`
}

func newRoom(e *Entrance , id int , title string) *Room{
	return &Room{
		id: id ,
		entrance: e ,
		broadcast: make(chan SendData , 1),
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make( map[*Client]Client),
		title: title ,
	}
}

func refleshUsers(room *Room , urclient *Client , isRegister bool ){
			//同じルームのユーザーを表示する
			var sendData SendData
			var users []string
			for client := range(room.clients){
				users = append( users , client.name)
			}

			if isRegister {
				sendData.Type = "register"
			}else{
				sendData.Type = "unregister"
			}
			sendData.Message = urclient.name
			sendData.Users = users

			room.broadcast <- sendData
}

func( room *Room) run(){
	for{
		select{
		case client := <-room.register :
			log.Println("client registerd to room")
			room.clients[client] = *client

			refleshUsers(room , client , true )

		case client := <-room.unregister :
			if _,ok := room.clients[client]; ok{
				delete(room.clients , client)
				close(client.send)

				refleshUsers(room , client , false)
			}
		case message := <-room.broadcast:
			for client := range room.clients{
				select{
				case client.send <-message:
				default:
					close(client.send)
					delete(room.clients , client)
				}
			}
		}
	}
}
