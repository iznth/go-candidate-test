package main

import (
	"bytes"
	"math/rand"
	"sync"
	"time"
)

var (
	letters = "abcdefghijklmonpqrstuvwxyzABCDEFGHIJKLMONPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Create a random request ID
func newRequestId(idLen int) string {
	buff := new(bytes.Buffer)
	for i := 0; i < idLen; i++ {
		buff.WriteByte(letters[rand.Intn(len(letters))])
	}
	return buff.String()
}

type State string

var (
	StateUnknown  State = "Unknown"
	StateNew      State = "New"
	StateBusy     State = "Busy"
	StateFinished State = "Finished"
)

type Request struct {
	RequestState State
	Val          int
}

func NewRequest(val int) *Request {
	return &Request{
		RequestState: StateNew,
		Val:          val,
	}
}

type RequestManager struct {
	requests map[string]*Request
	mu       sync.RWMutex
}

func NewRequestManager() *RequestManager {
	return &RequestManager{
		requests: make(map[string]*Request),
	}
}

/*
Complete the below functions to enable concurrent processing of requests, querying their state and marking them as finished
*/

func (rm *RequestManager) QueueRequest(req *Request) (requestId string) {
	requestId = newRequestId(10)
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.requests[requestId] = req
	return requestId
}

func (rm *RequestManager) CompleteRequest(requestId string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if req, exists := rm.requests[requestId]; exists {
		req.RequestState = StateFinished
	}
}

func (rm *RequestManager) QueryRequestState(requestId string) State {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	if req, exists := rm.requests[requestId]; exists {
		return req.RequestState
	}
	return StateUnknown
}

// "I recommend using a main method and really making your solution work hard" -->
func main() {
	rm := NewRequestManager()

	// example
	req := NewRequest(42)
	reqID := rm.QueueRequest(req)

	// simulate
	time.Sleep(1 * time.Second)
	rm.CompleteRequest(reqID)

	state := rm.QueryRequestState(reqID)
	println("Request state:", state)
}
