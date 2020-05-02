package communication

import "moody-go/models"

const ()

type Client interface {
	Connect() error
	Discover()
	Forward(situation string)
	Update()
	Listen()
	Close()
}

type CommInterface struct {
	clients        []Client
	ConnectedNodes *models.ConnectedList
	DataTable      *models.DataTable
}

func (ifc *CommInterface) Init(confObj map[string]interface{}) error {
	return nil
}

func (ifc *CommInterface) Listen() {
	// TODO
}

func (ifc *CommInterface) Close() {
	// TODO
}
