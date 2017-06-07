// Copyright 2017
// Author: catlittlechen@gmail.com

package nexus

import (
	"errors"
)

var (
	// RootIndex the index of root
	RootIndex = 0
)

var (
	// ErrConfFormat the format of configuration is wrong
	ErrConfFormat = errors.New("wrong conf format")
	// ErrConf bad configuration
	ErrConf = errors.New("wrong conf")

	// ErrNodeNotFound node not found
	ErrNodeNotFound = errors.New("node not found")
	// ErrNodeReplicate node replicate
	ErrNodeReplicate = errors.New("node replicate")
	// ErrNodeHasChildren node has children
	ErrNodeHasChildren = errors.New("node has children")

	// ErrKeyNotFound key not found
	ErrKeyNotFound = errors.New("key not found")
)
