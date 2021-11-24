package main

import "fmt"

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("[INFO]: Add was called")
	c.queue.Add(obj)
}

func (c *controller) handleDelete(obj interface{}) {
	fmt.Println("[INFO]: Delete was called")
	c.queue.Add(obj)
}
