package router

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/nehzata/test/events"
)

type handler[T events.Event] struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	v  reflect.Value
	ch chan T
}

func new_handler[T events.Event](r *Router, h func(T)) *handler[T] {
	ctx, cancel := context.WithCancel(context.Background())
	ret := &handler[T]{
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},

		v:  reflect.ValueOf(h),
		ch: make(chan T),
	}
	ret.wg.Add(1)
	go func() {
		defer close(ret.ch)
		defer ret.wg.Done()

		for {
			select {
			case <-ret.ctx.Done():
				return
			case evt := <-ret.ch:
				h(evt)
			}
		}
	}()
	return ret
}

type Router struct {
	handlers map[reflect.Type]*[]any
}

var r *Router = nil

func Init() {
	r = &Router{
		handlers: map[reflect.Type]*[]any{},
	}
}

func Close() {
	for k, h := range r.handlers {
		if len(*h) != 0 {
			fmt.Fprintf(os.Stderr, "dangling handler found %v", k)
		}
	}
}

func Subscribe[T events.Event](h func(T)) {
	val := reflect.ValueOf(h)

	if val.Kind() != reflect.Func {
		panic("trying to register an invalid listener")
	}

	eventType := val.Type().In(0).Elem()

	handlers, ok := r.handlers[eventType]
	if !ok {
		handlers = &[]any{}
		r.handlers[eventType] = handlers
	}

	*handlers = append(*handlers, new_handler(r, h))
}

func Unsubscribe[T events.Event](h func(T)) {
	val := reflect.ValueOf(h)

	if val.Kind() != reflect.Func {
		panic("trying to register an invalid listener")
	}

	eventType := val.Type().In(0).Elem()

	oldHandlers, ok := r.handlers[eventType]
	if !ok {
		return
	}

	newHandlers := &[]any{}
	for _, _hh := range *oldHandlers {
		if hh, ok := _hh.(*handler[T]); ok && hh.v.Pointer() == val.Pointer() {
			// Stop the event handler
			hh.cancel()

			// Wait for the handler to stop
			hh.wg.Wait()
		} else {
			*newHandlers = append(*newHandlers, _hh)
		}
	}

	r.handlers[eventType] = newHandlers
}

func Dispatch[T events.Event](evt T) {
	eventType := reflect.TypeOf(evt).Elem()
	if r, ok := r.handlers[eventType]; ok {
		for _, _h := range *r {
			h := _h.(*handler[T])
			h.ch <- evt
		}
	}

}
