package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().UTC().Format(time.RFC3339))
	t, _ := time.Parse("20060102t150405", "20190217t095000")
	fmt.Println(t)
}
