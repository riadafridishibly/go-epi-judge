package parity

func Parity(x int64) int16 {
	res := int16(0)
	for x != 0 {
		res ^= int16(x & 1)
		x >>= 1
	}

	return res
}
