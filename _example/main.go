package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	gomatchmaker "github.com/LIOU2021/go-match-maker"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var myHub *gomatchmaker.Hub

var count = 1000 // 模擬幾個參與者併發
var matchCount = 0
var matchCheck = map[string]int{}
var unMatchCount = 0
var unRegisterCount = 0
var unRegisterCountMX = sync.Mutex{}
var alreadyClose = false

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func init() {
	gomatchmaker.RegisterRedisClient(rdb)
}

func main() {
	config := gomatchmaker.Config{
		RegisterBuff:   200,
		BroadcastBuff:  200,
		UnRegisterBuff: 200,
		Room:           []string{"a", "b", "c", "d"},
		HubName:        "go-match-maker",
		// Mode:           gomatchmaker.Debug,
		Mode:     gomatchmaker.Release,
		Interval: time.Millisecond * 200,
		RegisterFunc: func(m gomatchmaker.Member) {
		},
		UnRegisterFunc: func(m gomatchmaker.Member) {
			unRegisterCountMX.Lock()
			defer unRegisterCountMX.Unlock()
			unRegisterCount++
		},
	}

	myHub = gomatchmaker.New(&config)

	go myHub.Run()

	testJoin()
	testLeave()
	fmt.Println("leave 結束======")
	go testNotification()

	go testNewData()

	time.AfterFunc(2*time.Second, func() {
		testClose(true)
		alreadyClose = true
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	testClose(false)
}

func testClose(mustClose bool) {
	if alreadyClose {
		return
	}

	for _, m := range myHub.GetMembers() {
		unMatchCount++
		fmt.Printf("剩餘roomId: %s, Id: %s\n", m.RoomId, m.Id)
	}
	fmt.Println("unRegisterCount: ", unRegisterCount)
	fmt.Println("unMatchCount: ", unMatchCount)
	myHub.Close()
	if mustClose {
		os.Exit(0)
	}
}

var testData = struct {
	sync.RWMutex
	list []*gomatchmaker.Member
}{}

func testNewData() {
	testNewData := &gomatchmaker.Member{ // 增加初始化不存在的room做测试
		Data:   99,
		RoomId: "e",
		Id:     uuid.New().String(),
	}
	myHub.Join(testNewData)
	myHub.Leave(testNewData)
}

func testJoin() {
	for i := 0; i < count; i++ {
		go func(index int) {
			m := &gomatchmaker.Member{
				Data:   index,
				Id:     uuid.New().String(),
				RoomId: testGetRoomId(index),
			}
			testData.Lock()
			testData.list = append(testData.list, m)
			testData.Unlock()
			myHub.Join(m)
		}(i)
	}
}

func testLeave() {
	i := 0
	for i < count {
		testData.RLock()
		fmt.Println("d83xjisgjipofdsgji")
		if i > len(testData.list)-1 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		myHub.Leave(testData.list[i])
		i = i + 2
		testData.RUnlock()
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
		msg := ""
		for _, v := range ms {
			msg += fmt.Sprint("[", v.RoomId, "] : ", v.Id, ", ")
			matchCount++

			if _, exists := matchCheck[v.Id]; exists {
				log.Fatalf("duplicate match - roomId: %s, id: %s", v.RoomId, v.Id)
			}
			matchCheck[v.Id]++
		}
		fmt.Println(msg)
		msg = ""
	}
	fmt.Println("matchCount: ", matchCount)
}
