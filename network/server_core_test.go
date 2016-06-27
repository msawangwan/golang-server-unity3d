package network

import (
	"github.com/msawangwan/unity-server/util"
	"testing"
	"time"
)

/* Does the server start and shutdown without issues? */
func TestServerStartAndShutdown(t *testing.T) {
	t.Log("\tDoes the server start and shutdown when signaled?")
	mockHostAddr := ":9081"
	mockServer := NewServerInstance(mockHostAddr)
	mockServer.Start()
	mockServer.Shutdown(true)
	t.Log("\t\tServer started then immediately sent a shutdown signal. [", util.PassMark, "]")
}

/* Can the server enter the listen loop and exit cleanly? */
func TestServerStartAndRunAndShutdown(t *testing.T) {
	t.Log("\tDoes the server start, listen and shutdown when signaled?")
	mockHostAddr := ":9081"
	mockServer := NewServerInstance(mockHostAddr)
	mockServer.Start()
	for i := 0; i < 5; i++ {
		time.Sleep(1000 * time.Millisecond)
	}
	mockServer.Shutdown(true)
	t.Log("\t\tServer started, ran through listen loop then was sent a kill signal. [", util.PassMark, "]")
}

// TODO: test fail cases
func TestServerAttemptBindWrongIP(t *testing.T) {
	t.Log("\t[Test case not implemented]")
}
