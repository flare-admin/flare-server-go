package manager

// Subscribe ， 订阅事件
type Subscribe struct {
	Id        string `json:"id"`
	Name      string `json:"name" `
	Topic     string `json:"topic" `
	Group     string `json:"group" `
	Constants string `json:"constants"`
	Status    int8   `json:"status"`
}
