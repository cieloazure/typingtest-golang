package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"strings"
)

const MaxPlayersPerRoom = 3

var nextPlayerId = 1
var nextRoomId = 1

var exampleTest = `
Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
`

type Player struct {
	id     int
	socket net.Conn
	data   chan []byte
}

type Room struct {
	id       int
	players  []*Player
	capacity int
}

func (room *Room) broadcastInRoom(message string) {
	for _, player := range room.players {
		player.data <- []byte(message)
	}
}

type PlayerMessageTuple struct {
	player  *Player
	message []byte
}

type RoomManager struct {
	rooms        map[int]*Room
	playerToRoom map[*Player]int
	register     chan *Player
	allMessages  chan *PlayerMessageTuple
}

func (manager *RoomManager) findRoom() (int, bool) {
	fmt.Println("In findRoom")
	for id, room := range manager.rooms {
		fmt.Printf("Room %d has capacity %d\n", id, room.capacity)
		if room.capacity < MaxPlayersPerRoom {
			return id, true
		}
	}
	return -1, false
}

func (manager *RoomManager) createRoom() int {
	fmt.Println("In createRoom")
	room := &Room{players: make([]*Player, 0, MaxPlayersPerRoom)}
	manager.rooms[nextRoomId] = room
	prevRoomId := nextRoomId
	nextRoomId++
	return prevRoomId
}

func (manager *RoomManager) receive(player *Player) {
	for {
		message := make([]byte, 4096)
		length, err := player.socket.Read(message)
		if err != nil {
			//manager.unregister <- player
			player.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED:" + string(message))
			message := &PlayerMessageTuple{player: player, message: message}
			manager.allMessages <- message
		}
	}
}

func (manager *RoomManager) send(player *Player) {
	defer player.socket.Close()
	for {
		select {
		case message, ok := <-player.data:
			if !ok {
				return
			}
			player.socket.Write(message)
		}
	}
}

func (manager *RoomManager) addPlayerToRoom(player *Player, roomid int) bool {
	fmt.Printf("%+v\n", player)
	manager.rooms[roomid].players = append(manager.rooms[roomid].players, player)
	manager.rooms[roomid].capacity++
	manager.playerToRoom[player] = roomid
	return manager.rooms[roomid].capacity == MaxPlayersPerRoom
}

func (manager *RoomManager) start() {
	for {
		select {
		case player := <-manager.register:
			whichRoom := -1
			if roomId, ok := manager.findRoom(); ok {
				whichRoom = roomId
			} else {
				whichRoom = manager.createRoom()
			}
			full := manager.addPlayerToRoom(player, whichRoom)
			fmt.Printf("Assigned room id %d to player %d\n", whichRoom, player.id)
			if full {
				fmt.Printf("Sending start messages to all players in room %d\n", whichRoom)
				room := manager.rooms[whichRoom]
				var buffer bytes.Buffer
				buffer.WriteString("START")
				buffer.WriteString(exampleTest)
				fmt.Println(buffer.String())
				room.broadcastInRoom(buffer.String())
			}

			//case tuple := <-manager.allMessages:
			//roomId := manager.playerToRoom[tuple.player]
			//room := manager.rooms[roomId]
			// TODO: Do something with the received message from
			// player of this room
		}
	}

}

func startServerMode() {
	fmt.Println("Starting server...")
	listener, error := net.Listen("tcp", ":12345")
	if error != nil {
		fmt.Println(error)
	}

	manager := RoomManager{
		rooms:        make(map[int]*Room),
		playerToRoom: make(map[*Player]int),
		register:     make(chan *Player),
	}

	go manager.start()

	for {
		connection, error := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		player := &Player{id: nextPlayerId, socket: connection, data: make(chan []byte)}
		nextPlayerId++

		manager.register <- player
		go manager.receive(player)
		go manager.send(player)
	}
}

type PlayerConnection struct {
	socket net.Conn
	data   chan []byte
}

func (player *PlayerConnection) receive() {
	for {
		message := make([]byte, 4096)
		length, err := player.socket.Read(message)
		if err != nil {
			player.socket.Close()
			break
		}
		if length > 0 {
			message := string(message)
			//fmt.Println("RECEIVED: " + message)
			if strings.HasPrefix(message, "START") {
				fmt.Println("Example test")
				fmt.Println(strings.TrimPrefix(message, "START"))
			}
		}
	}
}

func startClientMode() {
	fmt.Println("Starting player...")
	connection, error := net.Dial("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}
	player := &PlayerConnection{socket: connection}
	go player.receive()
	for {
	}
}

func main() {
	flagMode := flag.String("mode", "server", "start in client or server mode")
	flag.Parse()
	if strings.ToLower(*flagMode) == "server" {
		startServerMode()
	} else {
		startClientMode()
	}
}
