package main

import "github.com/philippecarle/moood/api/internal/bus"

func main() {
	conn := bus.Connection{}
	conn.Init()
}
