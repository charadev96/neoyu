package db

const (
	permFile = 0644
	permDir  = 0755
)

type Record interface {
	GetId() string
}
