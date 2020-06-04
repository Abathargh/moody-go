package communication

import (
	"gateway"
	"log"
	"strconv"
)

const (
	payloadBitSize = 64
)

// Handles incoming data from a sensor, checking if the topic-specified
// service exists. This is independent from the mqtt interface so that
// it may be used by other protocols in the future.
func DataHandler(datatype string, payload string) {
	_, err := strconv.ParseFloat(payload, payloadBitSize)
	if err != nil {
		log.Println("payload is not a number:" + payload)
		return
	}
	_, serviceIsActive := main.activeServices[datatype]
	if !serviceIsActive {
		log.Println("No such service")
		return
	}
	main.dataTable.Add(datatype, payload)
}
