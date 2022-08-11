package luna

func Luna(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(num int) int {
	var luhn int

	for i := 0; num > 0; i++ {
		cur := num % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		num = num / 10
	}

	return luhn % 10
}
