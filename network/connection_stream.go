package network

/* Handles read/writes between server and client */
type connectionStream struct {
}

func NewConnectionStream() *connectionStream {
	return &connectionStream{}
}

func (cs *connectionStream) BeginRead() {

}
