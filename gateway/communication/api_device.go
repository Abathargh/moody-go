package communication

import (
	"log"
	"strconv"
	"strings"
)

const (
	payloadBitSize = 64
)

// Handles incoming data from a sensor, checking if the topic-specified
// service exists. This is independent from the mqtt interface so that
// it may be used by other protocols in the future.
func DataHandler(datatype string, payload string) {
	floatPayload, err := strconv.ParseFloat(payload, payloadBitSize)
	if err != nil {
		log.Println("payload is not a number:" + payload)
		return
	}
	if !ActiveServices.Contains(datatype) {
		log.Printf("No such service '%s'\n", datatype)
		return
	}
	DataTable.Add(datatype, floatPayload)
}

func ActIPHandler(payload string) {
	if ok := isValidIP(payload); ok {
		log.Printf("a badly formatted ip was received from an actuator node: %s\n", payload)
		return
	}
	ActuatorIPs.Add(payload)
}

func isValidIP(ip string) bool {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return false
	}

	for _, octet := range octets {
		_, err := strconv.ParseUint(octet, 10, 8)
		if err != nil {
			return false
		}
	}
	return true
}
