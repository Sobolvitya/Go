package main

import (
	"fmt"
	"time"
	"strconv"
	"math/rand"
)
const FREQUENCY int = 500
const WAIT_TIME = 100
const CHECK_TIME = 1000
const CONCLUSION_TIME = 5000


var totalWorkerCount int = 2
var totalClientCount int = 5
var totalCountOfRejections = 0
var activeWorkerCount = 0
var responseCount = 0

var totalTaskGeneratedInLast5Seconds = 0
var totalTaskSolvedInLast5Seconds = 0
var totalTaskRejectedInLast5Seconds = 0


var write = make(chan Chunk, 20)
var stopClient = make(chan bool)
var stopWorker = make(chan bool)
var fail = make(chan bool, 20)
var activateWorker = make(chan bool, 20)
var increaseRequestCount = make(chan bool)
var taskGenerated = make(chan bool, 20)
var taskSolved = make(chan bool, 20)

type Chunk struct {
	channel chan string
	fileName string
}


func worker(number string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		
		case data := <- write:
			msg := data.fileName
			msg = msg + " : w" + number
			time.Sleep((time.Millisecond * time.Duration(r.Int31n(110))))
			data.channel <- msg
			activateWorker <- true
		case <- stopWorker:
			return
		}
	}
}

func client(number string) {

	tick := time.NewTicker(time.Duration(FREQUENCY) * time.Millisecond)
	var clientSendBackChanek = make(chan string)

	for t := range tick.C {

		dataToSend := Chunk{clientSendBackChanek, string(t.String() + " : c" + number)}
		write <- dataToSend
		taskGenerated <-  true
		ticker := time.NewTimer(WAIT_TIME * time.Millisecond)
		select {
			case res := <- clientSendBackChanek:
				res = res + ": c" + number
				increaseRequestCount <- true
				taskSolved <- true
			case <-ticker.C:
				<-clientSendBackChanek // possible bug, because we can wait for response forever
				fail <- true
			case <-stopClient:
				<-clientSendBackChanek
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
	tickerToCheckState := time.NewTicker(CHECK_TIME * time.Millisecond).C
	tickerToPrint := time.NewTicker(CONCLUSION_TIME * time.Millisecond).C
	for {
		select {
		case <- tickerToPrint:
			printData()
			clearConclusionData()
		case <- tickerToCheckState:

			if float64(totalCountOfRejections) > float64(responseCount) * 0.2 {
				fmt.Println("Add Worker")
				go worker("2")
				totalWorkerCount += 1
			}

			if float64(activeWorkerCount) < 0.5 * float64(totalWorkerCount) && totalWorkerCount > 0 {
				fmt.Println("Delete Worker")
				stopWorker <- true
				totalWorkerCount -= 1
			}

			totalCountOfRejections = 0
			responseCount = 0
			activeWorkerCount = 0

		case <- fail:
			increaseRejectionCounter()
		case <- taskGenerated:
			totalTaskGeneratedInLast5Seconds += 1
		case <- taskSolved:
			totalTaskSolvedInLast5Seconds += 1
		case <- increaseRequestCount:
			responseCount += 1
		case <- activateWorker:
			activeWorkerCount += 1
		}
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
	fmt.Println("\n===================================================")
	fmt.Println("Total worker count: ", totalWorkerCount)
	fmt.Println("Total client count: ", totalClientCount)
	fmt.Println("Client activity in milis: ", FREQUENCY)
	fmt.Println("Amount of tasks generated in last 5s : ", totalTaskGeneratedInLast5Seconds)
	fmt.Println("Amount of tasks solved in last 5s : ", totalTaskSolvedInLast5Seconds)
	fmt.Println("Amount of tasks rejected in last 5s : ", totalTaskRejectedInLast5Seconds)
	fmt.Println("Amout of active worker in last second:", activeWorkerCount)
	fmt.Println("===================================================\n")
}

func main() {
	fmt.Println("............Start........................")
	setUp()
	time.Sleep(1550 * time.Second)
}
