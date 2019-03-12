package main
import(
	"log"
)


//待機室
type Entrance struct{
	clients map[*Client]Client
	rooms map[*Room]Room
	register chan *Room
	enter chan *Client
	leave chan *Client
}

func newEntrance() *Entrance{
	return &Entrance{
		clients: make(map[*Client]Client),
		rooms : make(map[*Room]Room),
		register: make(chan *Room) ,
		enter: make(chan *Client),
		leave: make(chan *Client),
	}
}

func( entrance *Entrance) run(){
	for{
		select{
		//待機所への登録が行われた時
		case client := <-entrance.enter :
			entrance.clients[client] = *client
			//クライアントが選択したルームIDと一致するルームを探し
			//クライアントを登録する
			for room := range(entrance.rooms){
				if room.id == client.room_id {
					client.room = room
					room.register <- client
				}
			}

		case client := <-entrance.leave :
			if _ , ok := entrance.clients[client]; ok{
				delete(entrance.clients , client)
			}
		case room := <-entrance.register:
			log.Println(room)
			if _ , ok := entrance.rooms[room]; !ok{
				entrance.rooms[room] = *room
			}
		default:
		}
	}
}
