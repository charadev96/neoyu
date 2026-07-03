package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/goccy/go-yaml"
	"google.golang.org/protobuf/proto"
)

type File[T proto.Message] struct {
	file     string
	cache    T
	lastLoad time.Time
	mu       sync.Mutex
}

func NewFile[T proto.Message](file string) *File[T] {
	return &File[T]{
		file:  file,
		cache: newMessage[T](),
	}
}

func (s *File[T]) Save(data T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache = data
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
		return proto.Clone(s.cache).(T), nil
	}

	bytes, err := os.ReadFile(s.file)
	if err != nil {
		return *new(T), err
	}

	data := newMessage[T]()
	if err := yaml.Unmarshal(bytes, data); err != nil {
		return *new(T), fmt.Errorf("unmarshal yaml: %w", err)
	}

	s.cache = proto.Clone(data).(T)
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

func newMessage[T proto.Message]() T {
	var t T
	return t.ProtoReflect().New().Interface().(T)
}
