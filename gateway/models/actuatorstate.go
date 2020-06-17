package models

type ActuatorState struct {
	Mode bool `json:"mode" validate:"nonzero"`
}
