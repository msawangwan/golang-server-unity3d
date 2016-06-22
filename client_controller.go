package main

var activeConnections map[int]ClientHandler

func init() {
	activeConnections = make(map[int]ClientHandler)
}
