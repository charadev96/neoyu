package db

import "google.golang.org/protobuf/proto"

const (
	permFile = 0644
	permDir  = 0755
)

type Record interface {
	proto.Message
	GetId() string
}
