package communication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/models"
	"log"
	"net/http"
	"time"
)

// Observer that forwards data to the appropriate service
type DataForwarder struct {
	incomingDataChan chan models.DataEvent
	// TODO cache in here to save calls to external services?
}

// Creates a new DataForwarder that is an observer on the global
// data table.
// bufferSize specifies the incoming data buffered channel size.
func NewDataForwarder(bufferSize int) *DataForwarder {
	forwarder := &DataForwarder{
		incomingDataChan: make(chan models.DataEvent, bufferSize),
	}
	DataTable.Attach(forwarder.incomingDataChan)
	return forwarder
}

// Compares the keys representing the active services with the one tied to the dataset
// to see if they match and if the data should be forwarded to the next tier in the
// Moody architecture.
func sameStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aOccurrences := make(map[string]int)
	bOccurrences := make(map[string]int)

	for _, aElem := range a {
		aOccurrences[aElem]++
	}
	for _, bElem := range b {
		bOccurrences[bElem]++
	}

	for key, val := range aOccurrences {
		if val != bOccurrences[key] {
			return false
		}
	}
	return true
}

// Spawns a goroutine waiting for changes in the observed resource.
// The struct on which the method is called is a pointer
// for future use with the cache
func (df *DataForwarder) ListenForUpdates() {
	for evt := range df.incomingDataChan {
		switch NeuralState.Mode {
		case models.Stopped:
			log.Printf("New data service[%s]: %s", evt.ChangedKey, evt.ChangedValue)
		case models.Collecting:
			forwardToDataset(evt)
		case models.Predicting:
			forwardToNeural(evt)
		default:
			log.Println("Unexpected value for Neuralstate")
		}
	}
}

func forwardToDataset(evt models.DataEvent) {
	if sameStringSlice(ActiveServices.AsSlice(), DataTable.Keys()) {
		entry := evt.TableSnapshot
		// Decorate table copy with situation/time information
		entry["situation"] = float64(Situation.SituationId)
		hour, minute, _ := time.Now().Clock()
		entry["hour"] = float64(hour)
		entry["minute"] = float64(minute)
		req := models.DatasetEntryRequest{
			Entry: entry,
		}
		jsonTable, err := json.Marshal(req)
		if err != nil {
			log.Println(err)
			return
		}
		datasetEndpoint := fmt.Sprintf("%s/%s/%s", ApiGatewayAddress, "dataset", NeuralState.Dataset)
		resp, postErr := http.Post(datasetEndpoint, "application/json", bytes.NewBuffer(jsonTable))
		if postErr != nil {
			log.Println(postErr)
			return
		}
		if resp.StatusCode == http.StatusNotFound {
			log.Println("The specified dataset does not exist")
			return
		}
	}
}

func forwardToNeural(evt models.DataEvent) {
	if sameStringSlice(ActiveServices.AsSlice(), DataTable.Keys()) {
		query := evt.TableSnapshot
		req := models.NeuralPredictionRequest{
			DatasetEntryRequest: NeuralState.Dataset,
			Query:               query,
		}
		jsonTable, err := json.Marshal(req)
		if err != nil {
			log.Println(err)
			return
		}
		datasetEndpoint := fmt.Sprintf("%s/%s", ApiGatewayAddress, "neural")
		resp, postErr := http.Post(datasetEndpoint, "application/json", bytes.NewBuffer(jsonTable))
		if postErr != nil {
			log.Println(postErr)
			return
		}
		if resp.StatusCode == http.StatusNotFound {
			log.Println("The specified dataset does not exist")
			return
		}
		if resp.StatusCode == http.StatusUnprocessableEntity {
			log.Println("Wrong keys for the specified dataset")
			return
		}
		prediction := &models.NeuralPredictionResponse{}
		if err := models.ReadAndDecode(resp.Body, prediction); err != nil {
			log.Println(err)
			return
		}

		// Situation may not exist anymore, maybe check?
		CommForward(string(prediction.Situation))
	}
}

// Observer that expose data to the web app
type SocketioExposer struct {
	incomingDataChan     chan models.DataEvent
	incomingActuatorChan chan string
}

func NewSocketioExposer(bufferSize int) *SocketioExposer {
	exposer := SocketioExposer{
		incomingDataChan:     make(chan models.DataEvent, bufferSize),
		incomingActuatorChan: make(chan string, bufferSize),
	}
	DataTable.Attach(exposer.incomingDataChan)
	ActuatorIPs.Attach(exposer.incomingActuatorChan)
	return &exposer
}

func (se *SocketioExposer) ListenForUpdates() {
	for {
		select {
		case ip := <-se.incomingActuatorChan:
			evt := ActuatorConnectedEvent{IP: ip}
			SioServer.ForwardActuatorData(evt)
		case dataEvt := <-se.incomingDataChan:
			evt := ServiceDataEvent{
				Service: dataEvt.ChangedKey,
				Data:    dataEvt.ChangedValue,
			}
			SioServer.ForwardServiceData(evt)
		}
	}
}
