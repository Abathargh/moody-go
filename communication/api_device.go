package communication

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

// Handles incoming data from a sensor, checking if the topic-specified
// service exists.
func DataHandler(datatype string, payload string) {
	for serviceName, _ := range Services {
		if serviceName == datatype {
			DataTable.Add(datatype, payload)
			log.Println(DataTable.Data)
			return
		}
	}
	log.Println("No such service")
}
