package utils

func ParseInt(something string) int {
	negative := false
	num := 0
	i := 0

	if len(something) > 0 && something[0] == '-' {
		negative = true
		i = 1
	}

	for ; i < len(something); i++ {
		v := something[i]
		if v >= '0' && v <= '9' {
			num = num*10 + int(v-'0')
		} else {
			break
		}
	}

	if negative {
		num = -num
	}

	return num
}
