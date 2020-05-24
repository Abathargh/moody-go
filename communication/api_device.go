package communication

import (
	"log"
	"strconv"
)

const (
	payloadBitSize = 64
)

// Handles incoming data from a sensor, checking if the topic-specified
// service exists.
func DataHandler(datatype string, payload string) {
	for serviceName, _ := range Services {
		if serviceName == datatype {
			_, err := strconv.ParseFloat(payload, payloadBitSize)
			if err != nil {
				log.Println("payload is not a number:" + payload)
				return
			}
			DataTable.Add(datatype, payload)
			log.Println(DataTable)
			return
		}
	}
	log.Println("No such service")
}
