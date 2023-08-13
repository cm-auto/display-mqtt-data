package models

import "sync"

type DataContainer[T any] struct {
	RWMutex sync.RWMutex
	Data    []T
}
