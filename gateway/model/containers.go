package model

import (
	"fmt"
	"log"
	"sync"
)

// Observer Pattern

type DataEvent struct {
	ChangedKey   string
	ChangedValue string
	Keys         []string
	Values       []string
}

type DataObservable interface {
	Attach(evtChan chan<- DataEvent)
	Notify(evt DataEvent)
}

type Observer interface {
	ListenForUpdates()
}

//State table implementation

type DataTable struct {
	mutex        sync.Mutex
	data         map[string]string
	evtListeners []chan<- DataEvent
}

func NewDataTable() *DataTable {
	return &DataTable{data: make(map[string]string)}
}

func (table *DataTable) Add(key string, value string) {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	data, ok := table.data[key]
	if ok && data == value {
		return
	}
	table.data[key] = value
	k, v := table.keyValues()
	evt := DataEvent{
		Keys:   k,
		Values: v,
	}
	table.Notify(evt)
	log.Println(table)
}

func (table *DataTable) Remove(key string) {
	delete(table.data, key)
}

func (table *DataTable) Attach(evtChan chan<- DataEvent) {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	table.evtListeners = append(table.evtListeners, evtChan)
}

func (table *DataTable) Notify(evt DataEvent) {
	for _, evtChan := range table.evtListeners {
		evtChan <- evt
	}
}

func (table *DataTable) keyValues() ([]string, []string) {
	keys := make([]string, len(table.data))
	values := make([]string, len(table.data))
	i, j := 0, 0
	for key, value := range table.data {
		keys[i] = key
		values[j] = value
		i++
		j++
	}

	return keys, values
}

func (table *DataTable) String() string {
	var s string
	s = "{"
	for key, value := range table.data {
		s += fmt.Sprintf("%s: %s,", key, value)
	}
	if s[len(s)-1] != '{' {
		s = s[:len(s)-1]
	}
	s += "}"
	return s
}
