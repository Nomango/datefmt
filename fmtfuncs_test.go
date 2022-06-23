package datefmt

import "testing"

func Test_formatFuncs(t *testing.T) {
	check := func(t *testing.T, f func(p []byte, v, w int) []byte, v, w int, r string) {
		var p []byte
		p = f(p, v, w)
		if string(p) != r {
			t.Errorf("format(%d, %d) = %s; want %s", v, w, p, r)
		}
	}
	t.Run("formatNumProbably2Digits", func(t *testing.T) {
		check(t, formatNumProbably2Digits, 0, 0, "0")
		check(t, formatNumProbably2Digits, 9, 0, "9")
		check(t, formatNumProbably2Digits, 99, 0, "99")
		check(t, formatNumProbably2Digits, 999, 0, "999")
		check(t, formatNumProbably2Digits, 9, 1, "9")
		check(t, formatNumProbably2Digits, 99, 1, "99")
		check(t, formatNumProbably2Digits, 999, 1, "999")
		check(t, formatNumProbably2Digits, 9, 2, "09")
		check(t, formatNumProbably2Digits, 99, 2, "99")
		check(t, formatNumProbably2Digits, 999, 2, "999")
		check(t, formatNumProbably2Digits, 9, 3, "009")
		check(t, formatNumProbably2Digits, 99, 3, "099")
		check(t, formatNumProbably2Digits, 999, 3, "999")
	})

	t.Run("formatNumProbably3Digits", func(t *testing.T) {
		check(t, formatNumProbably3Digits, 0, 0, "0")
		check(t, formatNumProbably3Digits, 9, 0, "9")
		check(t, formatNumProbably3Digits, 999, 0, "999")
		check(t, formatNumProbably3Digits, 9999, 0, "9999")
		check(t, formatNumProbably3Digits, 9, 1, "9")
		check(t, formatNumProbably3Digits, 999, 1, "999")
		check(t, formatNumProbably3Digits, 9999, 1, "9999")
		check(t, formatNumProbably3Digits, 9, 3, "009")
		check(t, formatNumProbably3Digits, 999, 3, "999")
		check(t, formatNumProbably3Digits, 9999, 3, "9999")
		check(t, formatNumProbably3Digits, 9, 4, "0009")
		check(t, formatNumProbably3Digits, 999, 4, "0999")
		check(t, formatNumProbably3Digits, 9999, 4, "9999")
	})

	t.Run("formatNumProbably4Digits", func(t *testing.T) {
		check(t, formatNumProbably4Digits, 0, 0, "0")
		check(t, formatNumProbably4Digits, 9, 0, "9")
		check(t, formatNumProbably4Digits, 9999, 0, "9999")
		check(t, formatNumProbably4Digits, 99999, 0, "99999")
		check(t, formatNumProbably4Digits, 9, 1, "9")
		check(t, formatNumProbably4Digits, 9999, 1, "9999")
		check(t, formatNumProbably4Digits, 99999, 1, "99999")
		check(t, formatNumProbably4Digits, 9, 4, "0009")
		check(t, formatNumProbably4Digits, 9999, 4, "9999")
		check(t, formatNumProbably4Digits, 99999, 4, "99999")
		check(t, formatNumProbably4Digits, 9, 5, "00009")
		check(t, formatNumProbably4Digits, 9999, 5, "09999")
		check(t, formatNumProbably4Digits, 99999, 5, "99999")
	})

	t.Run("formatNum", func(t *testing.T) {
		check(t, formatNum, -1, 2, "-01")
		check(t, formatNum, 9, 30, "000000009")
	})

	check2 := func(t *testing.T, f func(p []byte, v uint) []byte, v uint, r string) {
		var p []byte
		p = f(p, v)
		if string(p) != r {
			t.Errorf("format(%d) = %s; want %s", v, p, r)
		}
	}
	t.Run("formatMax999", func(t *testing.T) {
		check2(t, formatMax999, 9, "9")
		check2(t, formatMax999, 99, "99")
		check2(t, formatMax999, 999, "999")
	})
	t.Run("formatMax9999", func(t *testing.T) {
		check2(t, formatMax9999, 9, "9")
		check2(t, formatMax9999, 99, "99")
		check2(t, formatMax9999, 999, "999")
		check2(t, formatMax9999, 9999, "9999")
	})
}
