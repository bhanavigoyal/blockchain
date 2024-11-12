package pkg

import "time"

var (
	PongWait = 10 * time.Second

	PingInterval = (PongWait * 9) / 10
)
