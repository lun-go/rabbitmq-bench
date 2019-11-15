package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	cmd := flagParse()
	logmq.SetLevel(Option.logLevel)

	var buider Builder
	switch Option.builder {
	case "time":
		buider = &BodyTimer{}
	case "order":
	default:
		fmt.Println("builder is need")
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	switch cmd {
	case "producer":
		wg.Add(1)
		go startProducer(wg, &Option, buider)
		wg.Wait()

	case "consumer":
		wg.Add(1)
		go startConsumer(wg, &Option, buider)
		if Option.consumerTimeout > 0 {
			time.Sleep(Option.consumerTimeout)
		} else {
			select {}
		}

	default:
		if strings.HasPrefix(cmd, "prepare") {
			Option.prepareMode = strings.Replace(cmd, "prepare ", "", 1)
			err := Prepare(&Option)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
