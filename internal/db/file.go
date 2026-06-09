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

const perm = 0644

type FileStore[G any] struct {
	file     string
	cache    *G
	lastLoad time.Time
	mu       sync.Mutex
}

func NewFileStore[G any](file string) *FileStore[G] {
	return &FileStore[G]{
		file:  file,
		cache: new(G),
	}
}

func (s *FileStore[G]) Save(data G) error {
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

	return os.WriteFile(s.file, bytes, perm)
}

func (s *FileStore[G]) Load() (G, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	info, err := os.Stat(s.file)
	if errors.Is(err, os.ErrNotExist) {
		if err := s.init(); err != nil {
			return *new(G), fmt.Errorf("create file: %w", err)
		}
	} else if err != nil {
		return *new(G), err
	}

	mod := info.ModTime().After(s.lastLoad)
	if !mod {
		var data G
		if err := deepcopy.Copy(&data, s.cache); err != nil {
			return *new(G), fmt.Errorf("copy from cache: %w", err)
		}
		return data, nil
	}

	bytes, err := os.ReadFile(s.file)
	if err != nil {
		return *new(G), err
	}

	var data G
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return *new(G), fmt.Errorf("unmarshal yaml: %w", err)
	}

	if err := deepcopy.Copy(s.cache, data); err != nil {
		return *new(G), fmt.Errorf("copy to cache: %w", err)
	}
	s.lastLoad = info.ModTime()

	return data, nil
}

func (s *FileStore[G]) init() error {
	if err := os.MkdirAll(filepath.Dir(s.file), perm); err != nil {
		return err
	}
	if err := os.WriteFile(s.file, nil, perm); err != nil {
		return err
	}
	return nil
}
