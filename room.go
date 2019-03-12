package main

import(
	"log"
)

type Room struct{
	id int
	entrance *Entrance
	clients map[*Client]Client
	broadcast chan []byte
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
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make( map[*Client]Client),
		title: title ,
	}
}

func( room *Room) run(){
	for{
		select{
		case client := <-room.register :
			log.Println("client registerd to room")
			room.clients[client] = *client
		case client := <-room.unregister :
			if _,ok := room.clients[client]; ok{
				delete(room.clients , client)
				close(client.send)
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
