package core

import (
	"github.com/hello/haneda/api"
)

type NoopBridge struct {
}

func (b *NoopBridge) PeriodicData(message *api.BatchedPeriodicData, key []byte) ([]byte, error) {
	return []byte{}, nil
}

func (b *NoopBridge) Logs(message *api.SenseLog, key []byte) error {
	return nil
}

func (b *NoopBridge) Pair(message *api.MorpheusCommand, privKey []byte) ([]byte, error) {
	return []byte{}, nil
}

func (b *NoopBridge) SenseState(message *api.SenseState, privKey []byte) ([]byte, error) {
	return []byte{}, nil
}
