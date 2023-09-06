package gomatchmaker

import (
	"context"
	"errors"
	"fmt"
	"log"
)

var ErrMemberExists = errors.New("member already exists")

func (h *Hub) RegisterEvent(m *Member) (err error) {
	h.Lock()
	defer h.Unlock()

	if !rdb.SIsMember(context.Background(), h.roomKey, m.RoomId).Val() {
		rdb.SAdd(context.Background(), h.roomKey, m.RoomId)
		fmt.Println("add room in set: ", m.RoomId)
	}

	memberKey := fmt.Sprintf("%s:member", h.roomKey)
	if err := rdb.SAdd(context.Background(), memberKey, m.Id).Err(); err != nil {
		log.Fatal(err)
	}

	if _, ok := h.members[m.Id]; ok {
		err = ErrMemberExists
		return
	}

	h.members[m.Id] = m
	fmt.Println("receive register: ", m)
	return
}

func (h *Hub) UnRegisterEvent(m *Member) {
	fmt.Println("receive unregister: ", m)
}
