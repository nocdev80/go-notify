# go-notify

Example 
```go
package main

import (
	"log"

	notify "github.com/nocv80/go-notify"
)

func main() {
	data := make(chan string)
	fx := func(observer *notify.ObserverNotify, event *notify.Event) {
		log.Println(observer.Filename, "  ", event.Name, "  ", event.Op.String())
	}

	notify.NewObserverNotify("./", "test.txt").
		FxCreate(fx).
		FxWrite(fx).
		FxChmod(fx).
		FxRemove(fx).
		FxRename(fx).
		Run()

	<-data
}
```

