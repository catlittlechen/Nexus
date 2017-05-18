// Copyright 2017
// Author: catlittlechen@gmail.com

package nexus

import (
	"net/url"
	"strconv"
	"strings"
	"sync"
)

// NodeManager manage nodes
type NodeManager struct {
	lock        *sync.RWMutex
	f           func(string) (DB, error)
	globalNode  map[int]*Node
	globalIndex int
	node        *Node
}

// NewNodeManager .
func NewNodeManager(f func(string) (DB, error)) (nm *NodeManager, err error) {
	nm = new(NodeManager)
	nm.lock = new(sync.RWMutex)
	nm.f = f
	nm.globalIndex = -1
	nm.globalNode = make(map[int]*Node)
	return
}

// addNode save node into map
func (nm *NodeManager) addNode(node *Node) (ok bool) {
	if _, ok = nm.globalNode[node.Index]; ok {
		return
	}
	nm.globalNode[node.Index] = node
	if nm.globalIndex < node.Index {
		nm.globalIndex = node.Index + 1
	}
	return
}

// delNode delete node from map
func (nm *NodeManager) delNode(index int) (err error) {
	err = nm.globalNode[index].Close()
	return
}

// getIndex  get the max index
func (nm *NodeManager) getIndex() (index int) {
	index = nm.globalIndex
	nm.globalIndex += 1
	return
}

// NewNode build node tree
func (nm *NodeManager) NewNode(conf string) (err error) {
	nm.lock.Lock()
	defer nm.lock.Unlock()
	if len(conf) != 0 {
		nm.node, err = nm.newNode(conf)
	} else {
		nm.node = &Node{
			Index:    RootIndex,
			Parent:   nil,
			DB:       nil,
			children: nil,
		}
		nm.node.Parent = nm.node
		nm.addNode(nm.node)
	}
	return
}

// newNode build node tree and return the root of tree
func (nm *NodeManager) newNode(conf string) (node *Node, err error) {

	node = new(Node)
	strArrays := strings.Split(conf, "&")
	for _, str := range strArrays {
		array := strings.Split(str, "=")
		if len(array) != 2 {
			err = ErrConfFormat
			return
		}
		switch array[0] {
		case "index":
			if node.Index, err = strconv.Atoi(array[1]); nil != err {
				return
			}

			if ok := nm.addNode(node); ok {
				err = ErrNodeReplicate
				return
			}
		case "parent":
			var pIndex int
			if pIndex, err = strconv.Atoi(array[1]); nil != err {
				return
			}
			var ok bool
			if node.Parent, ok = nm.globalNode[pIndex]; !ok {
				err = ErrNodeNotFound
				return
			}
		case "db":
			if array[1], err = url.QueryUnescape(array[1]); nil != err {
				return
			}
			if node.DB, err = nm.f(array[1]); nil != err {
				return
			}
		case "child":
			if node.children == nil {
				node.children = make([]*Node, 0)
			}
			if array[1], err = url.QueryUnescape(array[1]); nil != err {
				return
			}
			var n *Node
			if n, err = nm.newNode(array[1]); nil != err {
				return
			}
			node.children = append(node.children, n)
		}
	}

	if node.Parent == nil || (node.Index != RootIndex && (node.Index == node.Parent.Index || node.DB == nil)) {
		err = ErrConf
		return
	}

	return
}

// AddNode add the children of index node
func (nm *NodeManager) AddNode(index int, conf []string) (err error) {
	nm.lock.Lock()
	defer nm.lock.Unlock()

	if len(conf) == 0 {
		return
	}

	node, ok := nm.globalNode[index]
	if !ok {
		err = ErrNodeNotFound
		return
	}

	if len(node.children) != 0 {
		err = ErrNodeHasChildren
		return
	}

	node.children = make([]*Node, len(conf))
	for index, cfg := range conf {

		node.children[index] = &Node{}
		node.children[index].Index = nm.getIndex()
		node.children[index].Parent = node
		if node.children[index].DB, err = nm.f(cfg); nil != err {
			node.children = nil
			index--
			for index >= 0 {
				node.children[index].DB.Close()
				index--
			}
			return
		}

	}

	for _, n := range node.children {
		nm.addNode(n)
	}
	return
}

func (nm *NodeManager) getNode(key string, clean bool) (*Node, error) {
	hashKey := MD5(key)
	return nm.node.getNode(key, hashKey, 0, clean)
}

// Set find the bottom of tree
func (nm *NodeManager) Set(key string, value string) error {
	nm.lock.RLock()
	defer nm.lock.RUnlock()

	node, err := nm.getNode(key, true)
	if nil != err {
		return err
	}
	return node.Set(key, value)
}

// Get find and return
func (nm *NodeManager) Get(key string) (string, error) {
	nm.lock.RLock()
	defer nm.lock.RUnlock()

	node, err := nm.getNode(key, false)
	if nil != err {
		return "", err
	}
	return node.Get(key)
}

// Del. find and delete
func (nm *NodeManager) Del(key string) error {
	nm.lock.RLock()
	defer nm.lock.RUnlock()

	_, err := nm.getNode(key, true)
	return err
}

// String .
func (nm *NodeManager) String() string {
	if nm.node == nil {
		return ""
	}
	return nm.node.String()
}
