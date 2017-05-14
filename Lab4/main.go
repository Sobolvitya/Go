package main

import (
	"fmt"
	"time"
	"math/rand"
)
const FREQUENCY int = 50
var write = make(chan string)
var read = make(chan string)
func worker(number string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		msg := <- write
		msg = msg + " : w" + number
		read <- msg
		time.Sleep((time.Second * time.Duration(r.Int31n(4))))
	}
}

func client(number string) {
	tick := time.NewTicker(time.Duration(FREQUENCY) * time.Millisecond)
	for t := range tick.C {
		write <- string(t.String() + " : c" + number)
		res := <- read
		res = res + ": c" + number
		fmt.Println(res)
	}
}

func main() {
	fmt.Println("Start.................")
	go worker("1")
	go worker("2")
	go client("1")
	go client("2")
	go client("3")
	go client("4")
	go client("5")
	go client("6")
	go client("7")
	time.Sleep(1550 * time.Second)
}
