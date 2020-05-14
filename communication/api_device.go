package communication

import (
	"github.com/Abathargh/moody-go/models"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

// Handles incoming data from a sensor, checking if the topic-specified
// service exists.
func DataHandler(datatype string, payload string) {
	Services[datatype] = &models.Service{Id: 1, Name: datatype}
	for _, service := range Services {
		if service.Name == datatype {
			DataTable.Add(datatype, payload)
			return
		}
	}
	log.Println("No such service")
}
