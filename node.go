// Copyright 2017
// Author: catlittlechen@gmail.com

package nexus

import (
	"net/url"
	"strconv"
	"strings"
)

// Node root's parent is itself
type Node struct {
	Index    int
	Parent   *Node
	DB       DB
	children []*Node
}

// getNode
func (node *Node) getNode(key string, hashKey [16]byte, height int, clean bool) (n *Node, err error) {
	if len(node.children) == 0 {
		return nil, nil
	}

	// node with its children to share the keys
	index := int(hashValue(hashKey, height)) % len(node.children)
	if node.Index != RootIndex {
		index = int(hashValue(hashKey, height)) % (len(node.children) + 1)
		// if node should hold this key; then return
		if index == len(node.children) {
			n = node
			return
		}
	}

	if _, err = node.children[index].Get(key); nil != err && ErrKeyNotFound != err {
		return
	}

	if err == nil {
		// find the node; return node
		if !clean {
			return node.children[index], nil
		}

		// find the node but need clean; then clean and find the bottom of tree
		if err = node.children[index].Del(key); nil != err {
			return nil, err
		}
	}

	// when this node not found
	if n, err = node.children[index].getNode(key, hashKey, height+1, clean); nil != err {
		return
	}
	// its child has not children; return this node
	if n == nil {
		n = node.children[index]
		return
	}

	return
}

// Set .
func (node *Node) Set(key string, value string) error {
	if node == nil {
		return ErrNodeNotFound
	}
	return node.DB.Set(key, value)
}

// Get must return err = ErrKeyNotFound when key not found
func (node *Node) Get(key string) (string, error) {
	if node == nil {
		return "", nil
	}
	return node.DB.Get(key)
}

// Del .
func (node *Node) Del(key string) error {
	if node == nil {
		return nil
	}
	return node.DB.Del(key)
}

// String output the conf of the node
func (node *Node) String() string {
	var strArrays []string
	strArrays = append(strArrays, "index="+strconv.Itoa(node.Index))
	strArrays = append(strArrays, "parent="+strconv.Itoa(node.Parent.Index))
	if node.DB != nil {
		strArrays = append(strArrays, "db="+url.QueryEscape(node.DB.String()))
	}
	if node.children != nil {
		for _, n := range node.children {
			strArrays = append(strArrays, "child="+url.QueryEscape(n.String()))
		}
	}

	result := strings.Join(strArrays, "&")
	return result
}

// Close close the node
func (node *Node) Close() (err error) {
	err = node.DB.Close()
	return
}
