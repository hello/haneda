package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/hello/haneda"
	"github.com/hello/haneda/api"
	"net/http"
)

func main() {
	stats := make(chan api.Stat, 0)
	go haneda.DisplayStats(stats)

	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", ":6379")

		if err != nil {
			return nil, err
		}

		return c, err
	}, 10)

	defer redisPool.Close()
	vent := haneda.NewVentilator("example", redisPool)

	go vent.Publish()
	go vent.Listen()
	wsHandler := haneda.NewWsHandler(stats, vent)
	http.HandleFunc("/health", haneda.HealthHandler)
	http.Handle("/echo", wsHandler)
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
