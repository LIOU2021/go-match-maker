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
	}

	myHub = gomatchmakek.New(&config)

	go testJoin()
	go testNotification()
	time.AfterFunc(1*time.Second, myHub.Close)
	myHub.Run()
}

func testJoin() {
	for i := 0; i < 10; i++ {
		m := &gomatchmakek.Member{
			Data: i,
		}

		myHub.Join(m)
	}
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
