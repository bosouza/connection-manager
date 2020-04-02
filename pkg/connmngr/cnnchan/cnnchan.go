package cnnchan

import (
	"fmt"
	"sync"

	"github.com/souza-bruno/connection-manager/pkg/connmngr"
)

func CreateChannelCnnFactory(callerToCalleeChan, calleeToCallerChan chan string) *ChannelCnnFactory {
	return &ChannelCnnFactory{
		callerToCalleeChan: callerToCalleeChan,
		calleeToCallerChan: calleeToCallerChan,
	}
}

type ChannelCnnFactory struct {
	occupied           bool
	mu                 sync.Mutex
	callerToCalleeChan chan string
	calleeToCallerChan chan string
}

func (f *ChannelCnnFactory) CreateConnection() (connmngr.Connection, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.occupied {
		return nil, fmt.Errorf("client is currently occupied")
	}

	f.occupied = true
	callerToCalleeChan := make(chan string)
	calleeToCallerChan := make(chan string)
	done := make(chan struct{})

	//create a goroutine to read stuff sent by the caller through callerToCalleeChan to
	//the callee through f.callerToCalleeChan. Since we don't control f.callerToCallee
	//we have to consider that it may be closed, and we shouldn't allow the panic from
	//writting to a closed channel propagate and crash the program
	go func() {
		defer func() {
			recover()
		}()
		for {
			select {
			//since we don't expose callerToCalleeChan we know it will never be closed,
			//so we don't have to check for that condition
			case v := <-callerToCalleeChan:
				f.callerToCalleeChan <- v
			case <-done:
				return
			}
		}
	}()

	//create a goroutine to read stuff sent by the callee through f.calleeToCallerChan to
	//the caller through calleToCallerChan.
	go func() {
		for {
			select {
			case v, ok := <-f.calleeToCallerChan:
				if !ok {
					//f.calleeToCaller was closed, nothing to do anymore
					return
				}
				//since we don't expose calleeToCaller we know it will never be closed,
				//so we don't have to worry about this channel write causing a panic
				calleeToCallerChan <- v
			case <-done:
				return
			}
		}
	}()

	go func() {
		select {
		case <-done:
			f.mu.Lock()
			f.occupied = false
			f.mu.Unlock()
			return
		}
	}()

	return &ChannelCnn{
		callerToCalleeChan: callerToCalleeChan,
		calleeToCallerChan: calleeToCallerChan,
		done:               done,
	}, nil
}

type ChannelCnn struct {
	callerToCalleeChan chan string
	calleeToCallerChan chan string
	done               chan struct{}
}

func (c *ChannelCnn) Send(msg string) error {
	//this cannot panic, since we never close the channel and it's not exposed
	c.callerToCalleeChan <- msg
	return nil
}

func (c *ChannelCnn) Receive() (msg string, err error) {
	msg, ok := <-c.calleeToCallerChan
	if !ok {
		return "", fmt.Errorf("channel is closed")
	}
	return msg, nil
}

func (c *ChannelCnn) Close() (err error) {
	//attempting to close a closed channel will panic, so we recover and return an error
	//if Close() is called twice
	defer func() {
		if panic := recover(); panic != nil {
			//overwriting the return value with an error in case someone tries to close twice
			err = fmt.Errorf("connection already closed")
		}
	}()
	close(c.done)
	return nil
}
