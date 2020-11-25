package sum

func Sum(numbers []int) (value int) {
	for _, number := range(numbers) {
		value += number
	}
	return value
}

func SumAll(numbers ...[]int) (value []int) {
	for _, arr := range(numbers) {
		value = append(value, Sum(arr))
	}
	return value
}

func SumAllTail(numbers ...[]int) (value []int) {
	value = SumAll(numbers...)

	for i, sumAll := range(value) {
		if len(numbers[i]) > 0 {
			value[i] = sumAll - numbers[i][0]
		}
	}
	return value
}
