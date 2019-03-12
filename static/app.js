const NOT_ENTER = -1
let socket = null

let vm = new Vue({
  data() {
    return {
      inputName: "user",
      isLogin: false,
      isRoom: false,
      currentRoomId: NOT_ENTER,
      currentRoom: null,
      rooms_summary: [],
      chatMessage: "",
      chatAllPosts: []
    }
  },
  created() {
    //NONE
  },
  methods: {
    login() {
      if (this.inputName.length <= 0) {
        return
      }
      this.getCurrentSummary();
    },
    logout() {
      if (socket) {
        socket.close()
      }

      this.isLogin = false
      this.inputName = ""
    },
    enterRoom(room) {
      console.log("Enter Room SocketState: ", socket)
      console.log(`Enter to ${room.id}`)
      if (!socket) {
        socket = new WebSocket("ws://localhost:8080/ws")
        socket.onopen = (ev) => {
          console.log("Websocket Connect")
          socket.send(JSON.stringify({
            type: "enter",
            name: this.inputName,
            room_id: room.id,
            body: ""
          }))
        }

        socket.onmessage = (ev) => {
          console.log(ev)
          this.chatAllPosts.push(ev.data)
        }
      }

      this.isRoom = true
      this.currentRoom = room
    },
    exitRoom() {
      socket.send(JSON.stringify({
        type: "leave",
        name: this.inputName,
        room_id: this.currentRoomId,
        body: ""
      }))

      this.getCurrentSummary()
      this.isRoom = false
      socket.close()
      socket = null
      this.chatAllPosts = []
      this.chatMessage = ""
    },
    getCurrentSummary() {
      axios.get("/api/rooms/show").then(res => {
        this.rooms_summary = res.data
        this.isLogin = true
      })
    },
    postChatMessage() {
      if (socket) {
        socket.send(JSON.stringify({
          type: "post",
          name: this.inputName,
          room_id: Number(this.currentRoom.id),
          body: this.chatMessage
        }))

        this.chatMessage = ""
      }
    }
  },
  computed: {
    isSelectRooms() {
      return this.isLogin && !this.isRoom
    },
    isEnteredRooms() {
      return this.isLogin && this.isRoom
    }
  },
  el: '#app'
})
