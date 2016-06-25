package sense

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Store struct {
	messages map[uint64]time.Time
	sync.Mutex
}

func (s *Store) Save(messageId uint64) error {
	s.Lock()
	defer s.Unlock()
	_, found := s.messages[messageId]
	if found {
		return errors.New(fmt.Sprintf("duplicate message id: %d", messageId))
	}
	s.messages[messageId] = time.Now()
	return nil
}

func (s *Store) Expire(messageId uint64) {
	s.Lock()
	defer s.Unlock()
	delete(s.messages, messageId)
	// fmt.Println("expiring", messageId)
}

func (s *Store) State() {
	t := time.NewTicker(time.Second * 1)
	for {

		select {
		case <-t.C:
			s.Lock()
			notAcked := len(s.messages)
			if notAcked > 0 {
				fmt.Println("State:")

				for messageId, when := range s.messages {
					fmt.Println("\t", messageId, when)
				}
				fmt.Println("Done:")
			}
			s.Unlock()

		}
	}
}

func NewStore() *Store {
	return &Store{
		messages: make(map[uint64]time.Time),
	}
}
