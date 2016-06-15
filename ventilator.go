package haneda

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/api"
	"sync"
	"time"
)

type Ventilator struct {
	pairs map[string]*websocket.Conn
	sync.Mutex
	topic string
	pool  *redis.Pool
}

func NewVentilator(topic string, pool *redis.Pool) *Ventilator {
	return &Ventilator{
		topic: topic,
		pairs: make(map[string]*websocket.Conn),
		pool:  pool,
	}
}

func (v *Ventilator) register(senseId string, conn *websocket.Conn) {
	v.Lock()
	defer v.Unlock()
	v.pairs[senseId] = conn
	fmt.Printf("Registered: %s\n", senseId)
}

func (v *Ventilator) deregister(senseId string) {
	v.Lock()
	defer v.Unlock()
	delete(v.pairs, senseId)
	fmt.Printf("Deregistered: %s\n", senseId)
}

func (v *Ventilator) Listen() {

	for {
		// Get a connection from a pool
		c := v.pool.Get()
		psc := redis.PubSubConn{c}

		// Set up subscriptions
		psc.Subscribe("example")

		// While not a permanent error on the connection.
		for c.Err() == nil {
			switch val := psc.Receive().(type) {
			case redis.Message:
				jm := &api.JsonMessage{}
				json.Unmarshal(val.Data, jm)
				fmt.Printf("%s: message: %s\n", val.Channel, val.Data)
				fmt.Printf("name: %s, id:%d\n", jm.Name, jm.Id)
				v.Lock()
				conn := v.pairs[jm.Name]
				v.Unlock()
				m := serialize(jm.Id)
				if err := conn.WriteMessage(websocket.BinaryMessage, m); err != nil {
					v.deregister(jm.Name)

				}
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", val.Channel, val.Kind, val.Count)
			case error:
				fmt.Printf("%v", val)
			}
		}
		c.Close()
	}
}

// simulate messeji requests
func (v *Ventilator) Publish() {
	t := time.NewTicker(1 * time.Second)
	c := v.pool.Get()
	for {
		select {
		case <-t.C:
			names := make([]string, 0)
			v.Lock()
			for k, _ := range v.pairs {
				names = append(names, k)
			}
			v.Unlock()
			for _, name := range names {
				jm := &api.JsonMessage{Name: name, Id: uint64(time.Now().Unix())}
				buff, _ := json.Marshal(jm)
				c.Do("PUBLISH", v.topic, buff)
			}
		}
	}
}
