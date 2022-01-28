package even_odd_array

func EvenOdd(a []int) {
	if len(a) < 2 {
		return
	}

	evenIndex := 0
	oddIndex := len(a) - 1

	for evenIndex < oddIndex {
		if a[evenIndex]%2 == 1 && a[oddIndex]%2 == 0 {
			a[evenIndex], a[oddIndex] = a[oddIndex], a[evenIndex]
		}

		if a[evenIndex]%2 == 0 {
			evenIndex++
		}

		if a[oddIndex]%2 == 1 {
			oddIndex--
		}
	}
}
