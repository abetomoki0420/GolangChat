package main

import(
	"fmt"
	"net/http"
	"log"
	"encoding/json"
)

func serveHome( w http.ResponseWriter , r *http.Request ){
	log.Println(r.URL)
	if r.URL.Path != "/"{
		http.Error(w, "Not found", http.StatusNotFound )
		return
	}

	if r.Method != "GET" {
		http.Error( w, "Method not allowed" , http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile( w,  r, "index.html")
}

type User struct{
	Type string `json:"type"`
	Name string `json:"name"`
	Room_id int `json:"room_id"`
	Body string `json:"body"`
}

func roomsHandler( w http.ResponseWriter , r *http.Request){
	rooms := []RoomSummary{}
	for _,room := range(entrance.rooms){
		rooms = append( rooms , RoomSummary{Title: room.title , Id: room.id , Count: len(room.clients)})
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	res ,_ := json.Marshal(rooms)
	w.Write(res)
}

var entrance *Entrance
func main(){
	entrance = newEntrance()
	go entrance.run()
	for i := 0 ; i < 3 ; i++ {
		room := newRoom(entrance , i , fmt.Sprintf("Room_%d" , i+1) )
		go room.run()
		entrance.register <- room
	}

	http.HandleFunc("/" , serveHome)
	http.HandleFunc("/api/rooms/show", roomsHandler)
	// http.HandleFunc("/api/rooms/create", roomsHandler)

	//ルームへのアクセス
	http.HandleFunc("/ws" , func( w http.ResponseWriter , r *http.Request ){
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn , err := upgrader.Upgrade( w, r , nil )
		if err != nil{
			log.Println(err)
			return
		}

		client := &Client{
			room: nil ,
			room_id: -1,
			conn: conn ,
			send: make(chan []byte , 256),
			name: "noname"  ,
		}

		go client.readPump()
		go client.writePump()

	})

	//静的ファイルのアクセス
	http.Handle("/static/" , http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(":8080" , nil )
	if err != nil {
		log.Fatal("ListenAndServe: " , err  )
	}
}
