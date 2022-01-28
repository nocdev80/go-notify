package notify

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type ObserverNotify struct {
	Filename     string
	Directory    string
	Watcher      *fsnotify.Watcher
	CurrentEvent *fsnotify.Event
	fxWrite      func(observer *ObserverNotify, event *Event)
	fxCreate     func(observer *ObserverNotify, event *Event)
	fxRemove     func(observer *ObserverNotify, event *Event)
	fxRename     func(observer *ObserverNotify, event *Event)
	fxChmod      func(observer *ObserverNotify, event *Event)
}
type Event fsnotify.Event

func (o *ObserverNotify) FxCreate(fxCreate func(observer *ObserverNotify, event *Event)) *ObserverNotify {
	o.fxCreate = fxCreate
	return o
}
func (o *ObserverNotify) FxWrite(fxWrite func(observer *ObserverNotify, event *Event)) *ObserverNotify {
	o.fxWrite = fxWrite
	return o
}
func (o *ObserverNotify) FxRemove(fxRemove func(observer *ObserverNotify, event *Event)) *ObserverNotify {
	o.fxRemove = fxRemove
	return o
}
func (o *ObserverNotify) FxRename(fxRename func(observer *ObserverNotify, event *Event)) *ObserverNotify {
	o.fxRename = fxRename
	return o
}
func (o *ObserverNotify) FxChmod(fxChmod func(observer *ObserverNotify, event *Event)) *ObserverNotify {
	o.fxChmod = fxChmod
	return o
}

func NewObserverNotify(directory string, filename string) *ObserverNotify {
	observer := &ObserverNotify{
		Filename:  filename,
		Directory: directory,
		fxWrite:   func(observer *ObserverNotify, event *Event) {},
		fxCreate:  func(observer *ObserverNotify, event *Event) {},
		fxRemove:  func(observer *ObserverNotify, event *Event) {},
		fxRename:  func(observer *ObserverNotify, event *Event) {},
		fxChmod:   func(observer *ObserverNotify, event *Event) {},
	}
	return observer
}

func (o *ObserverNotify) Run() {
	go func() {
		var err error
		o.Watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer o.Watcher.Close()

		err = o.Watcher.Add(o.Directory)
		if err != nil {
			log.Fatal(err)
		}

		done := make(chan bool)
		go func() {
			for {
				select {
				case event, ok := <-o.Watcher.Events:
					if !ok {
						return
					}
					if !strings.HasSuffix(event.Name, o.Filename) && o.Filename != "*" {
						continue
					}
					event1 := (Event)(event)
					switch {
					case event.Op&fsnotify.Write == fsnotify.Write:
						o.fxWrite(o, &event1)
					case event.Op&fsnotify.Create == fsnotify.Create:
						o.fxCreate(o, &event1)
					case event.Op&fsnotify.Remove == fsnotify.Remove:
						o.fxRemove(o, &event1)
					case event.Op&fsnotify.Rename == fsnotify.Rename:
						o.fxRename(o, &event1)
					case event.Op&fsnotify.Chmod == fsnotify.Chmod:
						o.fxChmod(o, &event1)
					}
				case err, ok := <-o.Watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()
		<-done
	}()
}
