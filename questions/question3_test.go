package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNewRequestManager(t *testing.T) { // test proper instantiation
	rm := NewRequestManager()
	if rm == nil {
		t.Fatal("NewRequestManager returned nil")
	}
	if len(rm.requests) != 0 {
		t.Fatal("New RequestManager should have no requests initially")
	}
}

func TestQueueRequest(t *testing.T) { // test proper enqueue
	rm := NewRequestManager()
	req := NewRequest(42)
	requestId := rm.QueueRequest(req)
	if requestId == "" {
		t.Fatal("QueueRequest should return a non-empty request ID")
	}
	if len(rm.requests) != 1 {
		t.Fatal("RequestManager should have one request after queuing a request")
	}
	if rm.requests[requestId].Val != 42 {
		t.Fatalf("Expected request value to be 42, got %d", rm.requests[requestId].Val)
	}
}

func TestCompleteRequest(t *testing.T) { // test proper finish state
	rm := NewRequestManager()
	req := NewRequest(42)
	requestId := rm.QueueRequest(req)
	rm.CompleteRequest(requestId)
	if rm.requests[requestId].RequestState != StateFinished {
		t.Fatalf("Expected request state to be %s, got %s", StateFinished, rm.requests[requestId].RequestState)
	}
}

func TestQueryRequestState(t *testing.T) { // ensure requestable
	rm := NewRequestManager()
	req := NewRequest(42)
	requestId := rm.QueueRequest(req)
	state := rm.QueryRequestState(requestId)
	if state != StateNew {
		t.Fatalf("Expected request state to be %s, got %s", StateNew, state)
	}
	rm.CompleteRequest(requestId)
	state = rm.QueryRequestState(requestId)
	if state != StateFinished {
		t.Fatalf("Expected request state to be %s, got %s", StateFinished, state)
	}
}

func TestQueryUnknownRequestState(t *testing.T) { // check StateUnknown
	rm := NewRequestManager()
	state := rm.QueryRequestState("nonexistent")
	if state != StateUnknown {
		t.Fatalf("Expected request state to be %s, got %s", StateUnknown, state)
	}
}

func TestConcurrentAccess(t *testing.T) { // safe concurrent access
	rm := NewRequestManager()
	var wg sync.WaitGroup
	const numRoutines = 100

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			req := NewRequest(val)
			reqID := rm.QueueRequest(req)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			rm.CompleteRequest(reqID)
			state := rm.QueryRequestState(reqID)
			if state != StateFinished {
				t.Errorf("Expected request state to be %s, got %s", StateFinished, state)
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentQuery(t *testing.T) { // concurrent query
	rm := NewRequestManager()
	req := NewRequest(42)
	reqID := rm.QueueRequest(req)

	var wg sync.WaitGroup
	const numRoutines = 100

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			state := rm.QueryRequestState(reqID)
			if state != StateNew && state != StateFinished {
				t.Errorf("Expected request state to be either %s or %s, got %s", StateNew, StateFinished, state)
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentMethodAccess(t *testing.T) { // concurrent method access
	rm := NewRequestManager()
	var wg sync.WaitGroup
	const numRoutines = 100
	const numRequests = 10

	for i := 0; i < numRequests; i++ {
		req := NewRequest(i)
		reqID := rm.QueueRequest(req)
		for j := 0; j < numRoutines; j++ {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				rm.QueryRequestState(id)
			}(reqID)

			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				rm.CompleteRequest(id)
			}(reqID)

			wg.Add(1)
			go func() {
				defer wg.Done()
				req := NewRequest(rand.Intn(100))
				rm.QueueRequest(req)
			}()
		}
	}

	wg.Wait()
}
