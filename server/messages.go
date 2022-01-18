package server

func newMessage() Message {
	return Message{Messages: make(map[string][]byte)}
}
func (m *Message) addMessage(name string, body []byte) {
	m.Messages[name] = body
}
func (d *deviceInfo) getMeStatus() Message {
	m := newMessage()
	m.addMessage("status", make([]byte, 0))
	return m
}
