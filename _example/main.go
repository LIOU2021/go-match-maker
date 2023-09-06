package main

import (
	"fmt"
	"time"

	gomatchmakek "github.com/LIOU2021/go-match-maker"
)

var myHub *gomatchmakek.Hub

func main() {
	config := gomatchmakek.Config{
		RegisterBuff:   200,
		BroadcastBuff:  200,
		UnRegisterBuff: 200,
		Room:           []string{"a", "b", "c", "d"},
		HubName:        "go-match-maker",
	}

	myHub = gomatchmakek.New(&config)

	go testJoin()
	go testLeave()
	go testNotification()

	time.AfterFunc(1*time.Second, myHub.Close)

	myHub.Run()
}

func testJoin() {
	for i := 0; i < 10; i++ {
		m := &gomatchmakek.Member{
			Data:   i,
			RoomId: testGetRoomId(i),
		}

		myHub.Join(m)
	}
}

func testLeave() {
	for i := 10; i > 0; i-- {
		m := &gomatchmakek.Member{
			Data:   i,
			RoomId: testGetRoomId(i),
		}

		myHub.Leave(m)
	}
}

func testGetRoomId(i int) (roomId string) {
	if i%4 == 0 {
		roomId = "a"
	} else if i%3 == 0 {
		roomId = "b"
	} else if i%2 == 0 {
		roomId = "c"
	} else {
		roomId = "d"
	}
	return roomId

}

func testNotification() {
	for ms := range myHub.Notification() {
		fmt.Print("receive notification: ")
		for _, v := range ms {
			fmt.Print(v.Data, ", ")
		}
		fmt.Println("")
	}
}
