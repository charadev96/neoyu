package db

import (
	"fmt"
	"slices"

	"github.com/charadev96/neoyu/internal/common"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Lens[S proto.Message, T Record] struct {
	db       *File[S]
	selector func(S) *[]T
}

func NewLens[S proto.Message, T Record](s *File[S], get func(S) *[]T) *Lens[S, T] {
	return &Lens[S, T]{s, get}
}

func (s *Lens[_, T]) List() ([]T, error) {
	data, err := s.db.Load()
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	return *s.selector(data), nil
}

func (s *Lens[_, T]) Load(id uuid.UUID) (T, error) {
	data, err := s.db.Load()
	if err != nil {
		return *new(T), fmt.Errorf("load: %w", err)
	}
	list := *s.selector(data)

	i := s.index(list, id)
	if i == -1 {
		return *new(T), common.ErrNotExist
	}

	return list[i], nil
}

func (s *Lens[_, T]) Save(val T) error {
	data, err := s.db.Load()
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	list := s.selector(data)

	i := s.index(*list, uuid.MustParse(val.GetId()))
	if i == -1 {
		*list = append(*list, val)
	} else {
		(*list)[i] = val
	}

	if err := s.db.Save(data); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

func (s *Lens[_, T]) Delete(id uuid.UUID) error {
	data, err := s.db.Load()
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	list := s.selector(data)

	i := s.index(*list, id)
	if i == -1 {
		return common.ErrNotExist
	}
	*list = slices.Delete(*list, i, i+1)

	if err := s.db.Save(data); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

func (s Lens[_, T]) index(list []T, id uuid.UUID) int {
	return slices.IndexFunc(list, func(val T) bool {
		return val.GetId() == id.String()
	})
}
