package main

import (
	"WB_L0_Task/app"
)

func main() {
	//go app.Test()
	app.Start()
}

//docker run -p 4222:4222 -p 8222:8222 nats-streaming -p 4222 -m 8222
