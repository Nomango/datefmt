package datefmt_test

import (
	"fmt"
	"time"

	"github.com/Nomango/datefmt"
)

func ExampleFastFormat() {
	t := time.Unix(1655732950, 181999999).In(time.UTC)
	s := datefmt.FastFormat(t, "yyyy-MM-dd hh:HH:mm:ss.SSS a")
	fmt.Println(s)
	// Output:
	// 2022-06-20 01:13:49:10.181 PM
}
