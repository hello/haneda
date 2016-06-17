package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/hello/haneda"
	"net/http"
)

func proxy(w http.ResponseWriter, r *http.Request) {

}

type PublishHanlder struct {
	pool  *redis.Pool
	topic string
}

func (h *PublishHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := h.pool.Get()
	msg := `{"name": "Sense1", "id": 1234567890}`
	conn.Do("PUBLISH", h.topic, msg)
	fmt.Fprintf(w, "%s", msg)
}

func webserver(topic string, pool *redis.Pool) {

	ph := &PublishHanlder{
		topic: topic,
		pool:  pool,
	}

	http.Handle("/publish", ph)
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func main() {
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", ":6379")

		if err != nil {
			return nil, err
		}

		return c, err
	}, 10)

	topic := "example"
	defer redisPool.Close()
	done := make(chan bool, 0)
	vent := haneda.NewVentilator(done, topic, redisPool)

	// go vent.Publish() // only required for testing
	go vent.Run()
	go vent.Listen()

	server := haneda.NewHelloServer(vent)
	go server.Run()

	go webserver(topic, redisPool)
	wsHandler := haneda.NewWsHandler(vent, server)
	http.HandleFunc("/health", haneda.HealthHandler)
	http.Handle("/echo", wsHandler)
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
