package main

import "log"

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	initConfig()
}

func main() {
	log.Printf("hello")
}
