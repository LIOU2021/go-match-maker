package gomatchmaker

import (
	"context"
	"errors"
	"fmt"
	"log"
)

var ErrMemberExists = errors.New("member already exists")
var ErrMemberNotExists = errors.New("member Not exists")

func (h *Hub) RegisterEvent(m *Member) (err error) {
	h.Lock()
	defer h.Unlock()

	if !rdb.SIsMember(context.Background(), h.roomKey, m.RoomId).Val() {
		rdb.SAdd(context.Background(), h.roomKey, m.RoomId)
		h.DebugLog("add room in set: %s\n", m.RoomId)
	}

	memberKey := fmt.Sprintf("%s:member:%s", h.roomKey, m.RoomId)
	if err := rdb.SAdd(context.Background(), memberKey, m.Id).Err(); err != nil {
		log.Fatal(err)
	}

	if _, ok := h.members[m.Id]; ok {
		return ErrMemberExists
	}

	h.members[m.Id] = m
	fmt.Printf("receive register - roomId: %s, userId: %s \n", m.RoomId, m.Id)
	return
}

func (h *Hub) UnRegisterEvent(m *Member) (err error) {
	h.Lock()
	defer h.Unlock()
	delete(h.members, m.Id)

	memberKey := fmt.Sprintf("%s:member:%s", h.roomKey, m.RoomId)

	if !rdb.SIsMember(context.Background(), memberKey, m.Id).Val() {
		return ErrMemberNotExists
	}

	if err := rdb.SRem(context.Background(), memberKey, m.Id).Err(); err != nil {
		log.Fatal(err)
	}

	if rdb.SCard(context.Background(), memberKey).Val() == 0 { // 该房间内没人了
		rdb.SRem(context.Background(), h.roomKey, m.RoomId) // 移除房间
	}

	fmt.Printf("receive unregister - roomId: %s, userId: %s \n", m.RoomId, m.Id)
	return
}
