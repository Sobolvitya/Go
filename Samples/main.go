package main

import "fmt"
import "time"
import "math/rand"

func printString(from string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
		time.Sleep((time.Second * time.Duration(r.Int31n(10))))
	}
}

func main() {
	go printString ("Viktor")
	printString("Vitya")
	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
}
