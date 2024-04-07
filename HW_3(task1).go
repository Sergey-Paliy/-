package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	input := make(chan string)
	output := make(chan int)
	var wg sync.WaitGroup

	wg.Add(2)
	go square(&wg, input, output)
	go multiply(&wg, output)
	fmt.Println("Вводите числа, для завершения ввода введите стоп")

	for {
		var inputStr string
		fmt.Scanln(&inputStr)

		if inputStr == "стоп" {
			close(input)
			break
		}

		input <- inputStr
	}

	wg.Wait()
}

func square(wg *sync.WaitGroup, input <-chan string, output chan<- int) {
	defer wg.Done()
	for numStr := range input {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Println("Ошибка преобразования числа:", err)
			continue
		}
		square := num * num
		output <- square
		fmt.Println("Квадрат числа:", square)
	}
	close(output)
}

func multiply(wg *sync.WaitGroup, input <-chan int) {
	defer wg.Done()
	for square := range input {
		product := square * 2
		fmt.Println("Квадрат числа умноженного на 2:", product)
	}
}
