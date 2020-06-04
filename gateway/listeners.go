package main

import (
	"gateway/model"
	"log"
)

// Observer that forwards data to the appropriate service
type DataForwarder struct {
	incomingDataChan chan model.DataEvent
	// TODO cache in here to save calls to external services?
}

// Creates a new DataForwarder that is an observer on the global
// data table.
// bufferSize specifies the incoming data buffered channel size.
func NewDataForwarder(bufferSize int) *DataForwarder {
	return &DataForwarder{
		incomingDataChan: make(chan model.DataEvent, bufferSize),
	}
}

// Spawns a goroutine waiting for changes in the observed resource.
// The struct on which the method is called is a pointer
// for future use with the cache
func (df *DataForwarder) ListenForUpdates() {
	for evt := range df.incomingDataChan {
		switch neuralState {
		case Stopped:
			log.Printf("New data service[%s]: %s", evt.ChangedKey, evt.ChangedValue)
		case Collecting:

		}
	}
}

// Observer that expose data to the web app

// TODO
