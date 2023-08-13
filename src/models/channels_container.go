package models

import "sync"

// since there are methods for a ChannelContainer
// and these methods would be the preferred way to
// add and remove channels, let's just make the
// fields of the struct private
type ChannelsContainer[T any] struct {
	mutex    sync.Mutex
	channels []chan T
}

func (self *ChannelsContainer[T]) Add() chan T {
	channel := make(chan T)
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.channels = append(self.channels, channel)
	return channel
}

func (self *ChannelsContainer[T]) RemoveAndClose(channel chan T) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	// if there are a lot of channels, this can be optimized
	// by using a map and removing by key
	for i, c := range self.channels {
		if c == channel {
			self.channels = append(self.channels[:i], self.channels[i+1:]...)
			close(channel)
		}
	}
}

func (self *ChannelsContainer[T]) Send(data T) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	for _, channel := range self.channels {
		channel <- data
	}
}
