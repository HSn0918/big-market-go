package util

// Convert
/**
 * 计算公式；
 * 1. 找到范围内最小的概率值，比如 0.1、0.02、0.003，需要找到的值是 0.003
 * 2. 基于1找到的最小值，0.003 就可以计算出百分比、千分比的整数值。这里就是1000
 */
func Convert(award float64) float64 {
	if award == 0 {
		return 1.0
	}

	max := 1.0
	for award < 1 {
		award *= 10
		max *= 10
	}
	return max
}
