package api

type Msg struct {
	Name string `json:"name"`
}

type JsonMessage struct {
	Name string `json:"name"`
	Id   uint64 `json:"id"`
}

type MessageParts struct {
	Header []byte
	Body   []byte
	Sig    []byte
}

type Stat struct {
	Source string
	Count  int
}
