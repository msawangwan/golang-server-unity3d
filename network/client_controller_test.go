package network

import (
	"github.com/msawangwan/unity-server/util"
	"net"
	"testing"
	"time"
)

/* Test the handling of a single client connection attempt. */
func TestASingleClientConnection(t *testing.T) {
	t.Log("\tCan the client controller respond to a single client connection?")
	mockHostAddr := ":9082"
	mockServer := NewServerInstance(mockHostAddr)
	mockServer.Start()

	go func() {
		conn, err := net.Dial("tcp", mockHostAddr)
		if err != nil {
			t.Fatal("\t\tMock client failed to establish a connection to the server. [ ", util.FailMark, " ]")
		}
		//t.Log("\t\t\tClient connection handled ...")
		time.Sleep(3000 * time.Millisecond)
		conn.Close()
		//t.Log("\t\t\tClosing client connection. ")
		return
	}()

	time.Sleep(5000 * time.Millisecond)
	t.Log("\t\tClient Controller was able to handle a single client connection. [ ", util.PassMark, " ]")
}
