package main

import (
	"fmt"
	"time"

	gomatchmakek "github.com/LIOU2021/go-match-maker"
	"github.com/google/uuid"
)

var myHub *gomatchmakek.Hub

func main() {
	config := gomatchmakek.Config{
		RegisterBuff:   200,
		BroadcastBuff:  200,
		UnRegisterBuff: 200,
		Room:           []string{"a", "b", "c", "d"},
		HubName:        "go-match-maker",
		// Mode:           gomatchmakek.Debug,
		Mode: gomatchmakek.Release,
	}

	myHub = gomatchmakek.New(&config)

	testJoin()
	testLeave()
	go testNotification()

	go testNewData()

	time.AfterFunc(2*time.Second, myHub.Close)

	myHub.Run()
}

var testData []*gomatchmakek.Member

func testNewData() {
	testNewData := &gomatchmakek.Member{ // 增加初始化不存在的room做测试
		Data:   99,
		RoomId: "e",
		Id:     uuid.New().String(),
	}
	myHub.Join(testNewData)
	myHub.Leave(testNewData)
}

func testJoin() {
	for i := 0; i < 6; i++ {
		m := &gomatchmakek.Member{
			Data:   i,
			Id:     uuid.New().String(),
			RoomId: testGetRoomId(i),
		}
		testData = append(testData, m)
		myHub.Join(m)
	}
}

func testLeave() {
	for i := 0; i < 3; i++ {
		m := testData[i]

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
