package modules

import (
	"sync"
)

var (
	nodeInfo  DataTable
	connected ConnectedList
	onceNodes sync.Once
	onceData  sync.Once
)

// Virtualization of the generic Node
type Node struct {
	Mac       string
	Type      string
	Group     string
	Datatypes []string
}

type Action int

const (
	Added Action = iota
	Removed
)

// Observer Pattern

type NodesObservable interface {
	Attach(evtChan chan NodeEvent)
	Notify(evt NodeEvent)
}

type NodeEvent struct {
	Action Action
	Node   Node
}

type DataObservable interface {
	Attach(evtChan chan DataEvent)
	Notify(evt DataEvent)
}

type DataEvent struct {
	Datatype string
	Payload  string
}

type Observer interface {
	ListenForUpdates()
}

// Connected List definition and API

type ConnectedList struct {
	mutex        sync.Mutex
	nodes        []Node
	EvtListeners []chan NodeEvent
}

func ConnectedNodes() *ConnectedList {
	onceNodes.Do(func() {
		connected = ConnectedList{}
	})
	return &connected
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
		if elem.Mac == node.Mac {
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

func (list *ConnectedList) Attach(evtChan chan NodeEvent) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.EvtListeners = append(list.EvtListeners, evtChan)
}

func (list *ConnectedList) Notify(evt NodeEvent) {
	for _, evtChan := range list.EvtListeners {
		evtChan <- evt
	}
}

//State table implementation
type DataTable struct {
	mutex        sync.Mutex
	data         map[string]string
	EvtListeners []chan DataEvent
}

func NodesData() *DataTable {
	onceData.Do(func() {
		nodeInfo = DataTable{
			data: make(map[string]string),
		}
	})
	return &nodeInfo
}

func (table *DataTable) Add(key string, value string) {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	table.data[key] = value
	evt := DataEvent{
		Datatype: key,
		Payload:  value,
	}
	table.Notify(evt)
}

func (table *DataTable) Attach(evtChan chan DataEvent) {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	table.EvtListeners = append(table.EvtListeners, evtChan)
}

func (table *DataTable) Notify(evt DataEvent) {
	for _, evtChan := range table.EvtListeners {
		evtChan <- evt
	}
}