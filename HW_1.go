package main

import (
	"fmt"
)

func insertarray(array1 [4]int, array2 [5]int) [9]int {
	var combarray [9]int
	for i := 0; i < len(array1); i++ {
		combarray[i] = array1[i]
	}
	for i := 0; i < len(array2); i++ {
		combarray[i+len(array1)] = array2[i]
	}
	return combarray

}
func sortarray(combaray [9]int) [9]int {
	for i := 0; i < len(combaray)-1; i++ {
		for j := 0; j < len(combaray)-i-1; j++ {

			if combaray[j] > combaray[j+1] {
				combaray[j], combaray[j+1] = combaray[j+1], combaray[j]
			}
		}
	}

	return combaray
}

func main() {

	var array1 [4]int
	var array2 [5]int
	fmt.Println("Введите первый массив")
	for i := 0; i < len(array1); i++ {
		fmt.Scan(&array1[i])
	}
	fmt.Println("Введите второй массив")
	for i := 0; i < len(array2); i++ {
		fmt.Scan(&array2[i])
	}
	fmt.Println(insertarray(array1, array2))
	fmt.Println(sortarray(insertarray(array1, array2)))
}
