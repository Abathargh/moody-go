package models

type StateValue uint64

type Service struct {
	Id   uint64 `gorm:"column:id;primary_key" json:"id"`
	Name string `gorm:"column:name;unique" validate:"nonzero" json:"name"`
}

type ServiceRequest struct {
	Service uint64 `json:"serviceId" validate:"nonzero"`
}

type ServiceByName struct {
	ServiceName string `json:"name" validate:"nonzero"`
}

type ServicesResponse struct {
	Services     []string `json:"services"`
	ServiceCount int      `json:"count"`
}

type ServiceTableResponse struct {
	Table        map[string]string `json:"services"`
	ServiceCount int               `json:"count"`
}
