package main

func main() {
	server := NewServerInstance(":9081")
	server.Start()
	server.Shutdown()
}
