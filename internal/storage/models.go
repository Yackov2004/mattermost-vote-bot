package storage

// Poll - структура для хранения голосования
type Poll struct {
	ID       uint64   `lua:"id"`
	Question string   `lua:"question"`
	Options  []string `lua:"options"`
	Active   bool     `lua:"active"`
	OwnerID  string   `lua:"owner_id"`
}

// Vote - структура для хранения голоса
type Vote struct {
	PollID uint64
	UserID string
	Option string
}
