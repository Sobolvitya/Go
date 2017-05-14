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
	
	ticker := time.NewTicker(500 * time.Millisecond)
	
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	} ()
	time.Sleep(1550 * time.Millisecond)
	ticker.Stop()

	//c1 := make(chan string, 2)
	//
	//go func() {
	//	time.Sleep(1 * time.Second)
	//	c1 <- "c1"
	//}()
	//go func() {
	//	time.Sleep(1 * time.Second)
	//	c1 <- "c2"
	//}()
	//
	//msg1 := <-c1;
	//msg2 := <-c1;
	//
	//fmt.Println(msg1)
	//fmt.Println(msg2)
	
	
	//for i := 0; i < 2; i++ {
	//	select {
	//	case msg:= <-c1:
	//		fmt.Println(msg)
	//	case msg:= <-c2:
	//		fmt.Println(msg)
	//	}
	//}
	//
	//fmt.Println(<-messages)
	//fmt.Println(<-messages)
}
