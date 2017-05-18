// Copyright 2017
// Author: catlittlechen@gmail.com

package nexus

// DB is the interface of db, user must implement itself
type DB interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Del(key string) error
	Close() error

	String() string
}
