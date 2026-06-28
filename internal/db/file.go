package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/tiendc/go-deepcopy"
)

const (
	permFile = 0644
	permDir  = 0755
)

type File[T any] struct {
	file     string
	cache    *T
	lastLoad time.Time
	mu       sync.Mutex
}

func NewFile[T any](file string) *File[T] {
	return &File[T]{
		file:  file,
		cache: new(T),
	}
}

func (s *File[T]) Save(data T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache = &data
	bytes, err := yaml.Marshal(s.cache)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	if _, err := os.Stat(s.file); errors.Is(err, os.ErrNotExist) {
		if err := s.init(); err != nil {
			return fmt.Errorf("create file: %w", err)
		}
	}

	return os.WriteFile(s.file, bytes, permFile)
}

func (s *File[T]) Load() (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	info, err := os.Stat(s.file)
	if errors.Is(err, os.ErrNotExist) {
		if err := s.init(); err != nil {
			return *new(T), fmt.Errorf("create file: %w", err)
		}
		info, err = os.Stat(s.file)
		if err != nil {
			return *new(T), err
		}
	} else if err != nil {
		return *new(T), err
	}

	mod := info.ModTime().After(s.lastLoad)
	if !mod {
		var data T
		if err := deepcopy.Copy(&data, s.cache); err != nil {
			return *new(T), fmt.Errorf("copy from cache: %w", err)
		}
		return data, nil
	}

	bytes, err := os.ReadFile(s.file)
	if err != nil {
		return *new(T), err
	}

	var data T
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return *new(T), fmt.Errorf("unmarshal yaml: %w", err)
	}

	if err := deepcopy.Copy(s.cache, data); err != nil {
		return *new(T), fmt.Errorf("copy to cache: %w", err)
	}
	s.lastLoad = info.ModTime()

	return data, nil
}

func (s *File[_]) init() error {
	if err := os.MkdirAll(filepath.Dir(s.file), permDir); err != nil {
		return err
	}
	if err := os.WriteFile(s.file, nil, permFile); err != nil {
		return err
	}
	return nil
}
