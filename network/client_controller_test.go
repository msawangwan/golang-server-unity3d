package network

import (
	"github.com/msawangwan/unity-server/util"
	"net"
	"testing"
)

const mockHostAddr = ":9082"

func TestASingleClientConnection(t *testing.T) {
	t.Log("\tCan the client controller respond to a single client connection?")
	mockServer := NewServerInstance(mockHostAddr)
	mockServer.Start()

	go func() {
		conn, err := net.Dial("tcp", mockHostAddr)
		if err != nil {
			t.Fatal("\t\tMock client failed to establish a connection to the server. [", util.FailMark, "]")
		}
		conn.Close()
		return
	}()

	t.Log("\t\tClient Controller was able to handle a single client connection. [", util.PassMark, "]")
}
