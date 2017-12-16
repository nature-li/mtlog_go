package main

import (
	"fmt"
	"time"
	"mtlog"
	"flag"
)

func main() {
	var pRound = flag.Int64("round", 100000, "loop round")
	var pLittle = flag.Bool("little", true, "big or false")

	var message string

	if *pLittle {
		message = "01234567890中国"
	} else {
		for i := 0; i < 1024; i++ {
			message += "1234567890中国"
		}
	}
	flag.Parse()
	logger := mtlog.NewLogger(mtlog.DEVELOP, mtlog.INFO, "logs", "server", 10240, -1)
	if !logger.Start() {
		fmt.Println("logger.Start failed")
	}
	logger.SetLevel(mtlog.TRACE)

	var i int64 = 0
	round := *pRound
	start := time.Now()
	for i = 0; i < round; i++ {
		logger.Trace(message)
		logger.Debug(message)
		logger.Info(message)
		logger.Warn(message)
		logger.Error(message)
		logger.Fatal(message)
		logger.Report(message)
	}
	logger.Stop()
	end := time.Now()

	nanoseconds := end.Sub(start).Nanoseconds()
	speed := 1000000000 * round * 7.0 / nanoseconds

	fmt.Printf("round: %v\n", round)
	fmt.Printf("little: %v\n", *pLittle)
	fmt.Printf("Nanoseconds: %v\n", nanoseconds)
	fmt.Printf("speed: %v\n", speed)
}

