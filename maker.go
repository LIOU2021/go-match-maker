package gomatchmaker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type mode string

var Debug mode = "debug"
var Release mode = "release"
var closeSignal = make(chan struct{}, 1)

type Config struct {
	HubName        string
	RegisterBuff   int
	BroadcastBuff  int
	UnRegisterBuff int
	Room           []string // 有哪些房间
	Mode           mode
}

type Member struct {
	Data   interface{}
	Id     string // user id
	RoomId string
}

type Hub struct {
	register   chan *Member       // 加入撮合
	broadcast  chan []*Member     // 撮合成功推播
	unRegister chan *Member       // 退出撮合
	shutDown   chan struct{}      // 关闭服务
	roomKey    string             // 存放在缓存的key名称
	members    map[string]*Member // 存总用户。key为user id
	sync.Mutex
	mode mode // 模式
}

func (h *Hub) GetMembers() map[string]*Member {
	return h.members
}

// new hub instance
func New(config *Config) *Hub {
	if len(config.Room) > 0 {
		err := rdb.SAdd(context.Background(), config.HubName, config.Room).Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	var mode mode
	if config.Mode == Release {
		mode = Release
	} else {
		mode = Debug
	}

	return &Hub{
		register:   make(chan *Member, config.RegisterBuff),
		broadcast:  make(chan []*Member, config.BroadcastBuff),
		unRegister: make(chan *Member, config.UnRegisterBuff),
		shutDown:   make(chan struct{}),
		roomKey:    config.HubName,
		members:    make(map[string]*Member),
		mode:       mode,
	}
}

// execute match maker
func (h *Hub) Run() {
	fmt.Println("run match maker")
	go h.executeMatchRunner()

	for {
		select {
		case m := <-h.register:
			go func() {
				if err := h.RegisterEvent(m); err != nil {
					fmt.Printf("[registerEvent Err] id: %s - %v\n", m.Id, err)
				}
			}()
		case m := <-h.unRegister:
			go func() {
				if err := h.UnRegisterEvent(m); err != nil {
					fmt.Printf("[unregisterEvent Err] id: %s - %v\n", m.Id, err)
				}
			}()
		case <-h.shutDown:
			close(h.register)
			close(h.broadcast)
			close(h.unRegister)
			close(h.shutDown)

			if h.mode == Release {
				for {
					var keys []string
					var err error
					var cursor uint64
					keys, cursor, err = rdb.SScan(context.Background(), h.roomKey, cursor, "*", 10).Result()
					if err != nil {
						log.Fatal(err)
					}

					for _, roomId := range keys {
						memberKey := fmt.Sprintf("%s:member:%s", h.roomKey, roomId)
						rdb.Del(context.Background(), memberKey)
					}

					// 没有更多key了
					if cursor == 0 {
						break
					}
				}

				rdb.Del(context.Background(), h.roomKey)

			}

			fmt.Println("close match maker")
			closeSignal <- struct{}{}
			return
		}
	}
}

// 这个方法会堵塞，直到正常关闭hub
func (h *Hub) Close() {
	h.Lock()
	defer h.Unlock()

	go func() {
		h.shutDown <- struct{}{}
	}()
	<-closeSignal
	close(closeSignal)
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

// 撮合配对
func (h *Hub) executeMatchRunner() {
	ctx := context.Background()
	for {
		if len := rdb.SCard(ctx, h.roomKey).Val(); len < 2 {
			time.Sleep(200 * time.Millisecond)
			fmt.Println("群组数量不足，等待中...")
			continue
		}

		rooms := rdb.SRandMemberN(ctx, h.roomKey, 2).Val()
		r1, r2 := rooms[0], rooms[1]
		fmt.Printf("r1: %s, r2: %s \n", r1, r2)
		memberKey1 := fmt.Sprintf("%s:member:%s", h.roomKey, r1)
		memberKey2 := fmt.Sprintf("%s:member:%s", h.roomKey, r2)

		h.Lock()

		uid1 := rdb.SPop(ctx, memberKey1).Val()
		uid2 := rdb.SPop(ctx, memberKey2).Val()

		h.broadcast <- []*Member{
			h.members[uid1],
			h.members[uid2],
		}
		delete(h.members, uid1)
		delete(h.members, uid2)
		h.Unlock()
	}
}
