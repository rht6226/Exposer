package main

import "fmt"

// Called when a deployment is added.
func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("[INFO]: Add was called")
	c.queue.Add(obj)
}

// called when a deployment is deleted.
func (c *controller) handleDelete(obj interface{}) {
	fmt.Println("[INFO]: Delete was called")
	c.queue.Add(obj)
}
