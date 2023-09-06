package gomatchmaker

import (
	"context"
	"fmt"
)

func (h *Hub) RegisterEvent(m *Member) {
	if !rdb.SIsMember(context.Background(), h.roomKey, m.RoomId).Val() {
		rdb.SAdd(context.Background(), h.roomKey, m.RoomId)
		fmt.Println("add room in set: ", m.RoomId)
	}
	fmt.Println("receive register: ", m)
}

func (h *Hub) UnRegisterEvent(m *Member) {
	fmt.Println("receive unregister: ", m)
}
