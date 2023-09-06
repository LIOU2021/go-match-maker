package gomatchmaker

import (
	"context"
	"fmt"
	"log"
)

type Config struct {
	HubName        string
	RegisterBuff   int
	BroadcastBuff  int
	UnRegisterBuff int
	Room           []string // 有哪些房间
}

type Member struct {
	Data   interface{}
	RoomId string
}

type Hub struct {
	register   chan *Member   // 加入撮合
	broadcast  chan []*Member // 撮合成功推播
	unRegister chan *Member   // 退出撮合
	shutDown   chan struct{}  // 关闭服务
	roomKey    string         // 存放在缓存的key名称
}

// new hub instance
func New(config *Config) *Hub {
	if len(config.Room) > 0 {
		err := rdb.SAdd(context.Background(), config.HubName, config.Room).Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Hub{
		register:   make(chan *Member, config.RegisterBuff),
		broadcast:  make(chan []*Member, config.BroadcastBuff),
		unRegister: make(chan *Member, config.UnRegisterBuff),
		shutDown:   make(chan struct{}),
		roomKey:    config.HubName,
	}
}

// execute match maker
func (h *Hub) Run() {
	fmt.Println("run match maker")
	for {
		select {
		case m := <-h.register:
			go h.RegisterEvent(m)
		case m := <-h.unRegister:
			go h.UnRegisterEvent(m)
		case <-h.shutDown:
			close(h.register)
			close(h.broadcast)
			close(h.unRegister)
			close(h.shutDown)
			rdb.Del(context.Background(), h.roomKey)
			fmt.Println("close match maker")
			return
		}
	}
}

func (h *Hub) Close() {
	h.shutDown <- struct{}{}
}

// join match
func (h *Hub) Join(member *Member) {
	h.register <- member
}

// leave match
func (h *Hub) Leave(member *Member) {
	h.unRegister <- member
}

// receive match notification
func (h *Hub) Notification() <-chan []*Member {
	return h.broadcast
}
