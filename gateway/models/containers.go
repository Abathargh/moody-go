package models

import (
	"sync"
)

// Observer Pattern

type NodesObservable interface {
	Attach(evtChan chan NodeEvent)
	Notify(evt NodeEvent)
}

type ActionPerformed int

const (
	Added ActionPerformed = iota
	Removed
)

type NodeEvent struct {
	Action ActionPerformed
	Node   Node
}

type DataObservable interface {
	Attach(evtChan chan<- map[string]string)
	Notify(evt map[string]string)
}

type Observer interface {
	ListenForUpdates()
}

// Connected List definition and API

type ConnectedList struct {
	mutex        sync.Mutex
	nodes        []Node
	evtListeners []chan<- NodeEvent
}

func (list *ConnectedList) Add(node Node) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.nodes = append(list.nodes, node)
	evt := NodeEvent{
		Action: Added,
		Node:   node,
	}
	list.Notify(evt)
}

func (list *ConnectedList) Remove(node Node) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	for index, elem := range list.nodes {
		if elem.MacAddress == node.MacAddress {
			toDelete := elem
			list.nodes[len(list.nodes)-1], list.nodes[index] = list.nodes[index], list.nodes[len(list.nodes)-1]
			list.nodes = list.nodes[:len(list.nodes)-1]
			evt := NodeEvent{
				Action: Removed,
				Node:   toDelete,
			}
			list.Notify(evt)
		}
	}
}

func (list *ConnectedList) Attach(evtChan chan<- NodeEvent) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.evtListeners = append(list.evtListeners, evtChan)
}

func (list *ConnectedList) Notify(evt NodeEvent) {
	for _, evtChan := range list.evtListeners {
		evtChan <- evt
	}
}

//State table implementation
type DataEvent struct {
	Keys   []string
	Values []string
}

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
	table.data[key] = value
	k, v := table.keyValues()
	evt := DataEvent{
		Keys:   k,
		Values: v,
	}
	table.Notify(evt)
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
