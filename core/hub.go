package core

import (
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/hello/haneda/sense"
	"sync"
	"time"
)

var (
	ErrSenseNotConnected = errors.New("sense not connected")
	ErrInvalidMessage    = errors.New("invalid mp. missing sense id")
)

type ConnectionRemover interface {
	Remove(senseId sense.SenseId)
}

type ConnectionAdder interface {
	Add(senseId sense.SenseId) chan *sense.MessageParts
}

type Hub struct {
	sync.Mutex
	pairs          map[string]time.Time // map[string] because map[SenseId] always returns a different pointer value
	chanPairs      map[string]chan *sense.MessageParts
	chanBufferSize int
	logger         log.Logger
	stats          chan *HelloStat
}

func NewHub(logger log.Logger, stats chan *HelloStat) *Hub {
	return &Hub{
		pairs:          make(map[string]time.Time),
		chanPairs:      make(map[string]chan *sense.MessageParts),
		chanBufferSize: 2,
		logger:         logger,
		stats:          stats,
	}
}

func (h *Hub) Send(mp *sense.MessageParts) error {

	if mp == nil || mp.SenseId == "" {
		return ErrInvalidMessage
	}

	h.Lock()
	c, found := h.chanPairs[string(mp.SenseId)]
	h.Unlock()
	if !found {
		return ErrSenseNotConnected
	}
	c <- mp // potentially blocking?
	return nil
}

func (h *Hub) Remove(senseId sense.SenseId) {
	h.logger.Log("action", "removing", "sense_id", senseId)
	h.Lock()
	defer h.Unlock()
	when, found := h.pairs[string(senseId)]
	if !found {
		h.logger.Log("not_found", string(senseId))
		return
	}
	delete(h.pairs, string(senseId))
	duration := time.Now().Sub(when).Seconds()
	stats := &HelloStat{CurrConns: hFloat(len(h.pairs))}
	// Just making sure duration can't be negative
	if duration > 0 {
		stats.ConnDuration = hInt64(int(duration))
	}
	h.stats <- stats
	h.logger.Log("action", "removed", "sense_id", senseId)
}

func (h *Hub) Add(senseId sense.SenseId) chan *sense.MessageParts {
	h.logger.Log("action", "adding", "sense_id", senseId)
	h.Lock()
	defer h.Unlock()
	h.pairs[string(senseId)] = time.Now()
	h.stats <- &HelloStat{CurrConns: hFloat(len(h.pairs))}
	c := make(chan *sense.MessageParts, h.chanBufferSize)
	h.chanPairs[string(senseId)] = c
	return c
}
