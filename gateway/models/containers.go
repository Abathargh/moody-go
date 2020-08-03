package models

import (
	"fmt"
	"log"
	"sync"
)

// Observer Pattern

type DataEvent struct {
	ChangedKey    string
	ChangedValue  float64
	TableSnapshot map[string]float64
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
	data         map[string]float64
	evtListeners []chan<- DataEvent
}

func NewDataTable() *DataTable {
	return &DataTable{data: make(map[string]float64)}
}

func (table *DataTable) Add(key string, value float64) {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	current, ok := table.data[key]
	if ok && current == value {
		return
	}

	table.data[key] = value
	evt := DataEvent{
		ChangedKey:    key,
		ChangedValue:  value,
		TableSnapshot: table.data,
	}
	table.Notify(evt)
	log.Println(table)
}

func (table *DataTable) Remove(key string) {
	table.mutex.Lock()
	defer table.mutex.Unlock()
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

func (table *DataTable) Keys() []string {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	keys := make([]string, len(table.data))
	i := 0
	for key := range table.data {
		keys[i] = key
		i++
	}
	return keys
}

func (table *DataTable) String() string {
	var s string
	s = "{"
	for key, value := range table.data {
		s += fmt.Sprintf("%s: %f,", key, value)
	}
	if s[len(s)-1] != '{' {
		s = s[:len(s)-1]
	}
	s += "}"
	return s
}

// SynchronizedStringSet
type SynchronizedStringSet struct {
	evtListeners []chan<- string
	set          map[string]bool
	mutex        sync.Mutex
}

type StringObservable interface {
	Attach(evtChan chan<- string)
	Notify(evt string)
}

func NewSynchronizedStringSet() *SynchronizedStringSet {
	return &SynchronizedStringSet{
		set: make(map[string]bool),
	}
}

func (il *SynchronizedStringSet) Add(elem string) {
	il.mutex.Lock()
	defer il.mutex.Unlock()
	il.set[elem] = true
	il.Notify(elem)
}

func (il *SynchronizedStringSet) Remove(elem string) {
	il.mutex.Lock()
	defer il.mutex.Unlock()
	delete(il.set, elem)
}

func (il *SynchronizedStringSet) Clear() {
	il.mutex.Lock()
	defer il.mutex.Unlock()
	il.set = make(map[string]bool)
}

func (il *SynchronizedStringSet) Contains(elem string) bool {
	il.mutex.Lock()
	defer il.mutex.Unlock()
	return il.set[elem]
}

func (il *SynchronizedStringSet) Empty() bool {
	il.mutex.Lock()
	defer il.mutex.Unlock()
	return len(il.set) == 0
}

func (il *SynchronizedStringSet) AsSlice() []string {
	il.mutex.Lock()
	defer il.mutex.Unlock()

	i := 0
	keys := make([]string, len(il.set))
	for key := range il.set {
		keys[i] = key
		i++
	}
	return keys
}

func (il *SynchronizedStringSet) Attach(evtChan chan<- string) {
	il.mutex.Lock()
	defer il.mutex.Unlock()
	il.evtListeners = append(il.evtListeners, evtChan)
}

func (il *SynchronizedStringSet) Notify(evt string) {
	for _, evtChan := range il.evtListeners {
		evtChan <- evt
	}
}
