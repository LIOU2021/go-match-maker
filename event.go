package gomatchmaker

import "fmt"

func (h *Hub) RegisterEvent(m *Member) {
	fmt.Println("receive register: ", m)
}

func (h *Hub) UnRegisterEvent(m *Member) {
	fmt.Println("receive unregister: ", m)
}
