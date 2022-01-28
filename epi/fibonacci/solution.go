package fibonacci

type fibCalculator struct {
	cache map[int]int
}

func (fc *fibCalculator) calculate(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	if value, ok := fc.cache[n]; ok {
		return value
	}
	value := fc.calculate(n-1) + fc.calculate(n-2)
	fc.cache[n] = value

	return value
}

func Fibonacci(n int) int {
	fc := fibCalculator{
		cache: map[int]int{},
	}
	return fc.calculate(n)
}
