package datefmt

func formatString(p []byte, s string) []byte {
	return append(p, s...)
}

func formatNumProbably2Digits(p []byte, v, w int) []byte {
	if v >= 0 && v <= 99 {
		if w == 2 {
			return formatZero2Max99(p, uint(v))
		}
		if w == 1 {
			return formatMax99(p, uint(v))
		}
	}
	return formatNum(p, v, w)
}

func formatNumProbably3Digits(p []byte, v, w int) []byte {
	if v >= 0 && v <= 999 && w == 3 {
		return formatZero3Max999(p, uint(v))
	}
	return formatNum(p, v, w)
}

func formatNumProbably4Digits(p []byte, v, w int) []byte {
	if v >= 0 && v <= 9999 && w == 4 {
		return formatZero4Max9999(p, uint(v))
	}
	return formatNum(p, v, w)
}

func formatMax9(p []byte, v uint) []byte {
	return append(p, byte('0'+v))
}

func formatMax99(p []byte, v uint) []byte {
	if v > 9 {
		return formatMax9(formatMax9(p, v/10), v%10)
	}
	return formatMax9(p, v)
}

func formatMax999(p []byte, v uint) []byte {
	if v > 99 {
		return formatMax9(formatMax9(formatMax9(p, v/100), (v/10)%10), v%10)
	}
	return formatMax99(p, v)
}

func formatMax9999(p []byte, v uint) []byte {
	if v > 999 {
		return formatMax9(formatMax9(formatMax9(formatMax9(p, v/1000), (v/100)%10), (v/10)%10), v%10)
	}
	return formatMax999(p, v)
}

func formatZero2Max99(p []byte, v uint) []byte {
	if v > 9 {
		return formatMax99(p, v)
	}
	return formatMax9(formatMax9(p, 0), v)
}

func formatZero3Max999(p []byte, v uint) []byte {
	if v > 99 {
		return formatMax999(p, v)
	}
	return formatZero2Max99(formatMax9(p, 0), v)
}

func formatZero4Max9999(p []byte, v uint) []byte {
	if v > 999 {
		return formatMax9999(p, v)
	}
	return formatZero3Max999(formatMax9(p, 0), v)
}

func formatNum(p []byte, v, w int) []byte {
	u := uint(v)
	if v < 0 {
		p = append(p, '-')
		u = uint(-v)
	}

	var buf [10]byte
	if !(w <= 1 || (w == 2 && u > 9) || (w == 4 && u > 999)) {
		buf = [10]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	}

	n := len(buf)
	i := n
	for u >= 10 {
		i--
		buf[i] = byte('0' + u%10)
		u /= 10
	}
	i--
	buf[i] = byte('0' + u)

	if w > 0 && w > n-i {
		if w > 9 {
			w = 9
		}
		return append(p, buf[n-w:]...)
	}
	return append(p, buf[i:]...)
}
