package datefmt_test

import (
	"fmt"
	"time"

	"github.com/Nomango/datefmt"
)

func ExampleFastFormat() {
	t := time.Unix(1655689750, 0).In(time.UTC)
	s := datefmt.FastFormat(t, "yyyy-MM-dd HH:mm:ss")
	fmt.Println(s)
	// Output:
	// 2022-06-20 01:49:10
}
