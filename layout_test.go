package datefmt_test

import (
	"fmt"
	"time"

	"github.com/Nomango/datefmt"
)

func ExampleLayout() {
	l := datefmt.NewLayout("yyyy-MM-dd HH:mm:ss")
	t := time.Date(2022, time.June, 20, 21, 49, 10, 0, time.UTC)
	s := l.Format(t)

	fmt.Println(l, "=", s)
	// Output:
	// yyyy-MM-dd HH:mm:ss = 2022-06-20 21:49:10
}
