package communication

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

// Handles incoming data from a sensor, checking if the topic-specified
// service exists.
func DataHandler(datatype string, payload string) {
	for _, service := range Services {
		if service.Name == datatype {
			DataTable.Add(datatype, payload)
			return
		}
	}
	log.Println("No such service")
}
