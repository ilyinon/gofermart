package utils

func ValidLuhn(number string) bool {

	var sum int
	double := false

	for i := len(number) - 1; i >= 0; i-- {

		d := int(number[i] - '0')

		if double {

			d *= 2

			if d > 9 {
				d -= 9
			}
		}

		sum += d

		double = !double
	}

	return sum%10 == 0
}
