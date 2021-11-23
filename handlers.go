package main

import "fmt"

func HandleAdd(obj interface{}) {
	fmt.Println("[INFO]: Add was called")
}

func HandleDelete(obj interface{}) {
	fmt.Println("[INFO]: Delete was called")
}
