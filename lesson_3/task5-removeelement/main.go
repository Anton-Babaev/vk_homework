package main

import "fmt"

func RemoveElement(nums []int, index int) ([]int, error) {
	if index < 0 || index >= len(nums) {
		return nil, fmt.Errorf("индекс %d вне диапазона [0, %d)", index, len(nums))
	}

	result := make([]int, 0, len(nums)-1)
	result = append(result, nums[:index]...)
	result = append(result, nums[index+1:]...)
	return result, nil
}

func main() {
	original := []int{1, 2, 2, 3, 1}
	result, err := RemoveElement(original, 3)
	fmt.Println(result, err)
	fmt.Println("исходный слайс:", original)

	result, err = RemoveElement([]int{1, 2, 3}, 5)
	fmt.Println(result, err)

	result, err = RemoveElement([]int{1, 2, 3}, -1)
	fmt.Println(result, err)
}

// Почему append(nums[:index], nums[index+1:]...) испортит исходный слайс:
//
// nums[:index] — это СРЕЗ (slice) исходного массива nums, а не копия.
// Он смотрит в тот же массив, что nums. append к нему с достаточной
// ёмкостью пишет НОВЫЕ элементы прямо в исходный массив, начиная
// с позиции index. Это перезаписывает элементы nums после index.
//
// Пример: nums = [1, 2, 2, 3, 1], index = 3.
// nums[:3] = [1, 2, 2] (cap = 5, весь nums).
// append(nums[:3], nums[4:]...) добавит 1 на позицию 3 — перезапишет 3.
// Результат: nums = [1, 2, 2, 1, 1] — исходный испорчен.
//
// Решение: make + append в НОВЫЙ слайс (как в RemoveElement) или copy.
