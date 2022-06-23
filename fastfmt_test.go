package datefmt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Nomango/datefmt"
)

func ExampleFastFormat() {
	t := time.Unix(1655732950, 0).In(time.UTC)
	s := datefmt.FastFormat(t, "yyyy-MM-dd HH:mm:ss")
	fmt.Println(s)
	// Output:
	// 2022-06-20 13:49:10
}

func TestFastFormat(t *testing.T) {
	for _, tt := range tts {
		if tt.onlyStd {
			continue
		}
		r := datefmt.FastFormat(tt.t, tt.layout)
		if r != tt.formated {
			t.Errorf("FastFormat(%d, %s) = %s; want %s", tt.t.Unix(), tt.layout, r, tt.formated)
		}
	}
}
