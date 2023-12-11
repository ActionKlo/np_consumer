package main

import (
	"fmt"
	"np_consumer/internal/kafka"
)

func main() {
	fmt.Println("Hello Consumer!")

	kafka.Reader()
}
