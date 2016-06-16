package api

type Msg struct {
	Name string `json:"name"`
}

// JsonMessage is the message published by redis
type JsonMessage struct {
	Name string `json:"name"`
	Id   uint64 `json:"id"`
}

// MessageParts
type MessageParts struct {
	Header []byte
	Body   []byte
	Sig    []byte
}

// Stat holds statistics
type Stat struct {
	Source string
	Count  int
}

// AuthenticateFunc validates username password against db
type AuthenticateFunc func(username, password string) bool
