// Copyright 2017
// Author: catlittlechen@gmail.com

package nexus

import (
	"errors"
)

var (
	RootIndex = 0
)

var (
	ErrConfFormat = errors.New("wrong conf format")
	ErrConf       = errors.New("wrong conf")

	ErrNodeNotFound    = errors.New("node not found")
	ErrNodeReplicate   = errors.New("node replicate")
	ErrNodeHasChildren = errors.New("node has children")

	ErrKeyNotFound = errors.New("key not found")
)
