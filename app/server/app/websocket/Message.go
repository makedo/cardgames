package websocket

const (
    MESSAGE_TYPE_CONNECTED string = "connected"
    MESSAGE_TYPE_SELF_CONNECTED string = "self_connected"
    MESSAGE_TYPE_DISCONNECTED string = "disconnected"
    
    MESSAGE_TYPE_READY string = "ready"
    MESSAGE_TYPE_STATE string = "state"
    MESSAGE_TYPE_ERROR string = "error"
)

type Message struct {
    Type string `json:"type"`
    Data interface{} `json:"data"`
}

func (m *Message) GetData() map[string]interface{} {
    return m.Data.(map[string]interface{})
}
