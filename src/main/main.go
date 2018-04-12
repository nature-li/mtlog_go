package main

import (
	"fmt"
	"time"
	"mtlog"
	"flag"
)

var global  chan int

func logTest(logger *mtlog.Logger, round int64, message string) {
	var i int64
	for i = 0; i < round; i++ {
		logger.Trace(message)
		logger.Debug(message)
		logger.Info(message)
		logger.Warn(message)
		logger.Error(message)
		logger.Fatal(message)
		logger.Report(message)
	}

	global <- 1
}

func main() {
	var pRound = flag.Int64("round", 100000, "loop round")
	var pLittle = flag.Bool("little", true, "big or false")
	flag.Parse()

	global = make(chan int, 0)

	var message string
	if *pLittle {
		message = "01234567890中国"
	} else {
		for i := 0; i < 1024; i++ {
			message += "1234567890中国"
		}
	}
	logger := mtlog.NewLogger(true, mtlog.DEVELOP, mtlog.INFO, "logs", "server", 100 * 1024 * 1024, -1)
	if !logger.Start() {
		fmt.Println("logger.Start failed")
	}
	logger.SetLevel(mtlog.TRACE)

	round := *pRound
	start := time.Now()
	for t := 0; t < 4; t++ {
		go logTest(logger, round, message)
	}

	for t := 0; t < 4; t++ {
		<- global
	}
	logger.Stop()
	end := time.Now()

	nanoseconds := end.Sub(start).Nanoseconds()
	speed := 1000000000 * round * 7.0 / nanoseconds

	fmt.Printf("round: %v\n", round)
	fmt.Printf("little: %v\n", *pLittle)
	fmt.Printf("Nanoseconds: %v\n", nanoseconds)
	fmt.Printf("speed(QPS): %v\n", speed)
}

