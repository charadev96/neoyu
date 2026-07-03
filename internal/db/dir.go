package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charadev96/neoyu/internal/common"

	"github.com/google/uuid"
)

type Dir[T Record] struct {
	dir     string
	records map[uuid.UUID]*File[T]
	mu      sync.Mutex
}

func NewDir[T Record](dir string) *Dir[T] {
	return &Dir[T]{
		dir:     dir,
		records: make(map[uuid.UUID]*File[T]),
	}
}

func (s *Dir[T]) List() ([]T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := make([]T, 0, len(s.records))

	if err := s.ensure(); err != nil {
		return list, fmt.Errorf("ensure store: %w", err)
	}

	files, err := os.ReadDir(s.dir)
	if err != nil {
		return list, err
	}

	for k, f := range files {
		name := f.Name()
		id, err := uuid.Parse(strings.TrimSuffix(name, ".yaml"))
		if err != nil {
			// todo: Log warning on unknown file
			continue
		}

		if _, ok := s.records[id]; !ok {
			s.records[id] = NewFile[T](filepath.Join(s.dir, name))
		}

		data, err := s.records[id].Load()
		if err != nil {
			return list, fmt.Errorf("%s: %w", k, err)
		}
		list = append(list, data)
	}

	return list, nil
}

func (s *Dir[T]) Save(data T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.MustParse(data.GetId())

	if err := s.ensure(); err != nil {
		return fmt.Errorf("ensure store: %w", err)
	}
	if _, ok := s.records[id]; !ok {
		s.records[id] = NewFile[T](s.recordPath(id))
	}

	if err := s.records[id].Save(data); err != nil {
		return fmt.Errorf("%s: %w", id, err)
	}

	return nil
}

func (s *Dir[T]) Load(id uuid.UUID) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ensure(); err != nil {
		return *new(T), fmt.Errorf("ensure store: %w", err)
	}
	if _, ok := s.records[id]; !ok {
		return *new(T), fmt.Errorf("%s: %w", id.String(), common.ErrNotExist)
	}

	data, err := s.records[id].Load()
	if err != nil {
		return *new(T), fmt.Errorf("%s: %w", id, common.ErrNotExist)
	}

	return data, nil
}

func (s *Dir[_]) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ensure(); err != nil {
		return fmt.Errorf("ensure store: %w", err)
	}

	if err := os.Remove(s.recordPath(id)); err != nil {
		return fmt.Errorf("remove: %w", err)
	}

	delete(s.records, id)

	return nil
}

func (s *Dir[_]) recordPath(id uuid.UUID) string {
	return filepath.Join(s.dir, id.String()+".yaml")
}

func (s *Dir[T]) ensure() error {
	if s.records == nil {
		s.records = make(map[uuid.UUID]*File[T])
	}

	if _, err := os.Stat(s.dir); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(s.dir, permDir); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
