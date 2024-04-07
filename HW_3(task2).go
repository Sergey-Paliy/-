package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func squaresGenerator(ch chan<- int) {
	for i := 1; ; i++ {
		ch <- i * i
	}
}

func main() {
	squares := make(chan int)

	go squaresGenerator(squares)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case square := <-squares:
			fmt.Println(square)
		case <-interrupt:
			fmt.Println("\nПолучен сигнал прерывания. Выход из программы.")
			close(squares)
			return
		}
	}
}
