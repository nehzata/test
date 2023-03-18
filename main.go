package main

import (
	"fmt"
	"time"

	"github.com/nehzata/test/events"
	"github.com/nehzata/test/router"
)

func main() {
	router.Init()
	router.Subscribe(onTest1)
	router.Subscribe(onTest2)
	for i := 0; i < 5; i++ {
		router.Dispatch(&events.EventTest1{i})
		router.Dispatch(&events.EventTest2{i * 100})
		time.Sleep(time.Second)
	}
	router.Unsubscribe(onTest1)
	router.Unsubscribe(onTest2)
	router.Close()
}

func onTest1(evt *events.EventTest1) {
	fmt.Println("EventTest1", evt)
}

func onTest2(evt *events.EventTest2) {
	fmt.Println("EventTest2", evt)
}
