package main

import (
	"fmt"
	"time"
//	"math/rand"
	"strconv"
	"math/rand"
)
const FREQUENCY int = 500
const WAIT_TIME = 100
const CHECK_TIME = 1000
const CONCLUSION_TIME = 5000


var totalWorkerCount int = 2
var totalClientCount int = 10
var totalCountOfRejections = 0
var activeWorkerCount = 0

var totalTaskGeneratedInLast5Seconds = 0
var totalTaskSolvedInLast5Seconds = 0
var totalTaskRejectedInLast5Seconds = 0


var write = make(chan string)
var read = make(chan string)
var stopClient = make(chan bool)
var stopWorker = make(chan bool)
var fail = make(chan bool, 20)
var activateWorker = make(chan bool, 20)
var taskGenerated = make(chan bool, 20)
var taskSolved = make(chan bool, 20)


func worker(number string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		
		case msg := <- write:
			msg = msg + " : w" + number
			time.Sleep((time.Millisecond * time.Duration(r.Int31n(110))))
			read <- msg
			activateWorker <- true
			
		case <- stopWorker:
			return
		}
	}
}

func client(number string) {
	tick := time.NewTicker(time.Duration(FREQUENCY) * time.Millisecond)
	for t := range tick.C {
		write <- string(t.String() + " : c" + number)
		taskGenerated <-  true
		ticker := time.NewTimer(WAIT_TIME * time.Millisecond).C
		select {
			case res := <- read:
				res = res + ": c" + number
				fmt.Println(res)
				//taskSolved <- true
			case <-ticker:
				<- read
				fail <- true
				fmt.Println("We fucked it up, no fucking responce =( ")
			case <-stopClient:
				<-read
				return
		}
	}
}

func setUp() {
	fmt.Println("............Initializing.................")
	
	fmt.Println("............Set up Workers................")
	
	for i := 0; i < totalWorkerCount; i++ {
		go worker(strconv.FormatInt(int64(i), 10))
	}
	
	fmt.Println("............Set Up Client...................")
	
	
	for i := 0; i < totalClientCount; i++ {
		go client(strconv.FormatInt(int64(i), 10))
	}
	
	fmt.Println("............Run controller...................")
	controller()
}

func controller() {
	tickerToCheckState := time.NewTimer(CHECK_TIME * time.Millisecond).C
	tickerToPrint := time.NewTimer(CONCLUSION_TIME * time.Millisecond).C
	
	select {
		case <- tickerToPrint:
			printData()
			clearConclusionData()
		case <- tickerToCheckState:
			
			if float64(totalCountOfRejections) > float64(totalClientCount) * 0.2 {
				go worker("2")
				totalWorkerCount += 1
				totalCountOfRejections = 0
			}
			
			if float64(activeWorkerCount) < 0.5 * float64(totalWorkerCount) && totalWorkerCount > 0 {
				stopWorker <- true
				totalWorkerCount -= 1
			}
		case <- fail:
			increaseRejectionCounter()
		case <- taskGenerated:
			totalTaskGeneratedInLast5Seconds += 1
		case <- taskSolved:
			totalTaskSolvedInLast5Seconds += 1
	}
	
}

func increaseRejectionCounter() {
	totalCountOfRejections += 1
	totalTaskRejectedInLast5Seconds += 1
}

func clearConclusionData() {
	totalTaskGeneratedInLast5Seconds = 0
	totalTaskRejectedInLast5Seconds = 0
	totalTaskSolvedInLast5Seconds = 0
}

func printData() {
	fmt.Println("Total worker count: ", totalWorkerCount)
	fmt.Println("Total client count: ", totalClientCount)
	fmt.Println("Client activity in milis: ", FREQUENCY)
	fmt.Println("Amount of tasks generated in last 5s : ", totalTaskGeneratedInLast5Seconds)
	fmt.Println("Amount of tasks solved in last 5s : ", totalTaskSolvedInLast5Seconds)
	fmt.Println("Amount of tasks rejected in last 5s : ", totalTaskRejectedInLast5Seconds)
	
}

func main() {
	fmt.Println("............Start........................")
	setUp()
	
	time.Sleep(1550 * time.Second)
}
