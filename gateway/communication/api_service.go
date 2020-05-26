package communication

import (
	"github.com/Abathargh/moody-go/db"
	"github.com/Abathargh/moody-go/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type servicesResponse struct {
	Services     []*models.Service `json:"services"`
	ServiceCount int64             `json:"count"`
}

type serviceActivation struct {
	Action models.StateValue `validate:"min=0,max=1" json:"state"`
}

func getServices(w http.ResponseWriter, _ *http.Request) {
	services, count, err := db.GetAllServices()
	if err != nil {
		services = []*models.Service{}
	}
	resp := servicesResponse{
		Services:     services,
		ServiceCount: count,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	MustEncode(w, resp)
}

func postService(w http.ResponseWriter, r *http.Request) {
	newService := &models.Service{}
	w.Header().Set("Content-Type", "application/json")
	ok := MustValidate(r, newService)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}

	if err := db.AddService(newService); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		MustEncode(w, ErrorResponse{"Record with pk already exists"})
		return
	}
	w.WriteHeader(http.StatusOK)
	MustEncode(w, newService)
	if err := r.Body.Close(); err != nil {
		log.Println(err)
	}
}

// Get handler for /situation/{name}
func getService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}
	service, err := db.GetService(id)
	if err != nil {
		service = &models.Service{}
	}
	// here the response is the situation itself
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	MustEncode(w, service)
}

// Delete handler for /situation/{name}
func deleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}

	service, err := db.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		MustEncode(w, ErrorResponse{"The situation does not exist"})
	}

	if err := db.DeleteService(service); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		MustEncode(w, ErrorResponse{"An error occurred while trying to delete the situation"})
	}
	w.WriteHeader(http.StatusOK)
	MustEncode(w, service)
}

func patchService(w http.ResponseWriter, r *http.Request) {
	activation := &serviceActivation{}
	w.Header().Set("Content-Type", "application/json")
	ok := MustValidate(r, activation)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}
	service, err := db.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		MustEncode(w, ErrorResponse{"The service does not exist"})
		return
	}

	// Update in app Service Map
	tmpState := service.State
	service.State = activation.Action
	if err := db.PatchStateService(service); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		MustEncode(w, ErrorResponse{"Can't start the service"})
		return
	}
	if tmpState != activation.Action {
		switch activation.Action {
		case models.Stopped:
			delete(Services, service.Name)
		case models.Started:
			Services[service.Name] = service
		default:
			log.Println("a service can only be on or off")
		}
		log.Println(Services)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	MustEncode(w, service)

	if err := r.Body.Close(); err != nil {
		log.Println(err)
	}
}
