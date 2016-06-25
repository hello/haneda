package core

// import (
// 	"fmt"
// 	"github.com/garyburd/redigo/redis"
// 	"github.com/hello/haneda/api"
// 	"log"
// )

// type Ventilator struct {
// 	store map[string]*SenseConn
// 	in    chan *SenseConn
// 	out   chan string
// 	msg   chan *AckMessageWrapper
// 	raw   chan *api.MessageParts
// 	done  chan bool
// 	pool  *redis.Pool
// 	topic string
// }

// func (v *Ventilator) Receive(parsed *api.MessageParts) {
// 	v.raw <- parsed
// }

// func (v *Ventilator) Add(senseConn *SenseConn) {
// 	log.Printf("Adding: %s, %s", senseConn.SenseId, senseConn.Conn.RemoteAddr().String())
// 	v.in <- senseConn
// }

// func (v *Ventilator) Remove(senseId string) {
// 	v.out <- senseId
// }

// func (v *Ventilator) push(ack *AckMessageWrapper) {
// 	v.msg <- ack
// }

// func (v *Ventilator) Run() {
// 	fmt.Println("Ventilator running")
// outer:
// 	for {
// 		select {
// 		case m := <-v.in:
// 			v.store[m.SenseId] = m
// 		case m := <-v.out:
// 			delete(v.store, m)
// 		case m := <-v.msg:
// 			c, found := v.store[m.SenseId]
// 			if found {
// 				fmt.Println("Sending Ack", m.Id, "to", m.SenseId)
// 				fmt.Println("Signed with key", string(c.PrivKey))
// 				// TODO write to given sense!
// 				// c.Conn.WriteMessage(websocket.BinaryMessage, protoAckMessage)
// 			}
// 		case <-v.done:
// 			fmt.Println("Done")
// 			break outer
// 		default:

// 		}
// 	}
// 	close(v.in)
// 	close(v.out)
// 	close(v.msg)
// 	close(v.raw)
// }

// // Listen blocks and wait for messages to be published on the redis channel
// func (v *Ventilator) Listen() {

// 	for {
// 		// Get a connection from a pool
// 		c := v.pool.Get()
// 		psc := redis.PubSubConn{c}

// 		// Set up subscriptions
// 		psc.Subscribe(v.topic)

// 		// While not a permanent error on the connection.
// 		for c.Err() == nil {
// 			switch val := psc.Receive().(type) {
// 			case redis.Message:
// 				parsed, err := parse(val.Data)
// 				if err != nil {
// 					log.Println("Failed to parse message", err)
// 					continue
// 				}
// 				v.Receive(parsed)
// 			case redis.Subscription:
// 				fmt.Printf("%s: %s %d\n", val.Channel, val.Kind, val.Count)
// 			case error:
// 				fmt.Printf("%v", val)
// 			}
// 		}
// 		c.Close()
// 	}
// }

// func NewVentilator(done chan bool, topic string, pool *redis.Pool) *Ventilator {
// 	return &Ventilator{
// 		store: make(map[string]*SenseConn, 0),
// 		in:    make(chan *SenseConn, 2),
// 		out:   make(chan string, 2),
// 		msg:   make(chan *AckMessageWrapper, 16),
// 		raw:   make(chan *api.MessageParts, 16),
// 		done:  done,
// 		topic: topic,
// 		pool:  pool,
// 	}
// }

// // // Publish simulates messeji requests
// // func (v *Ventilator) Publish() {
// // 	t := time.NewTicker(1 * time.Second)
// // 	c := v.pool.Get()
// // 	for {
// // 		select {
// // 		case <-t.C:
// // 			names := make([]string, 0)
// // 			v.Lock()
// // 			for k, _ := range v.pairs {
// // 				names = append(names, k)
// // 			}
// // 			v.Unlock()
// // 			for _, name := range names {
// // 				jm := &api.JsonMessage{Name: name, Id: uint64(time.Now().Unix())}
// // 				buff, _ := json.Marshal(jm)
// // 				c.Do("PUBLISH", v.topic, buff)
// // 			}
// // 		}
// // 	}
// // }
