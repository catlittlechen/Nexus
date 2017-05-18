// Copyright 2017
// Author: catlittlechen@gmail.com

package main

import (
	"fmt"
	"github.com/catlittlechen/Nexus"
	"strconv"
)

func main() {
	root, err := nexus.NewNodeManager(NewRedis)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Printf("root: %+v\n", root)

	if err = root.NewNode(""); nil != err {
		fmt.Println(err)
		return
	}
	fmt.Printf("root: %+v\n", root)

	conf := []string{
		`{"addr":"127.0.0.1:6379", "db":0, "poolsize":1024, "idle":120}`,
		`{"addr":"127.0.0.1:6379", "db":1, "poolsize":1024, "idle":120}`,
		`{"addr":"127.0.0.1:6379", "db":2, "poolsize":1024, "idle":120}`,
	}
	if err = root.AddNode(nexus.RootIndex, conf); nil != err {
		fmt.Println(err)
		return
	}

	for i := 0; i <= 1024; i += 1 {
		if err = root.Set("hello"+strconv.Itoa(i), "world"); nil != err {
			fmt.Println(err)
			return
		}
	}

	fmt.Println(root.String())

	conf = []string{
		`{"addr":"127.0.0.1:6379", "db":3, "poolsize":1024, "idle":120}`,
		`{"addr":"127.0.0.1:6379", "db":4, "poolsize":1024, "idle":120}`,
		`{"addr":"127.0.0.1:6379", "db":5, "poolsize":1024, "idle":120}`,
	}
	if err = root.AddNode(2, conf); nil != err {
		fmt.Println(err)
		return
	}
	for i := 0; i <= 1024; i += 1 {
		if err = root.Set("hello"+strconv.Itoa(i), "world"); nil != err {
			fmt.Println(err)
			return
		}
	}

	configure := root.String()
	fmt.Println(configure)

	root2, err := nexus.NewNodeManager(NewRedis)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Printf("root: %+v\n", root2)

	if err = root2.NewNode(configure); nil != err {
		fmt.Println(err)
		return
	}
	fmt.Printf("root: %+v\n", root2)

	for i := 0; i <= 10; i += 1 {
		value, err := root2.Get("hello" + strconv.Itoa(i))
		fmt.Printf("key: hello%s value:%s err: %s\n", strconv.Itoa(i), value, err)
	}

	for i := 0; i <= 10; i += 1 {
		if err := root2.Del("hello" + strconv.Itoa(i)); nil != err {
			fmt.Println(err)
			return
		}
	}

	for i := 0; i <= 10; i += 1 {
		value, err := root2.Get("hello" + strconv.Itoa(i))
		fmt.Printf("key: hello%s value:%s err: %s\n", strconv.Itoa(i), value, err)
	}
}
