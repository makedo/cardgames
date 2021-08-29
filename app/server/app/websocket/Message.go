package websocket

type MessageType string

const (
    CONNECTED = "connected"
    SELF_CONNECTED = "self_connected"
    DISCONNECTED = "disconnected"
    READY = "ready"    
)

type Message struct {
    Type string `json:"type"`
    Data interface{} `json:"data"`
}
