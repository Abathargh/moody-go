package models

type DataMode uint8

const (
	Stopped DataMode = iota
	Collecting
	Predicting
)

type DatasetEntryRequest struct {
	Situation uint64             `json:"situation"`
	Entry     map[string]float64 `json:"entry" validate:"nonzero"`
}

type NeuralPredictionRequest struct {
	DatasetEntryRequest string             `json:"dataset"`
	Query               map[string]float64 `json:"query"`
}

type NeuralPredictionResponse struct {
	Situation int `json:"situation" validate:"nonzero,min=0"`
}

type NeuralState struct {
	Mode    DataMode `json:"mode" validate:"min=0,max=2"`
	Dataset string   `json:"dataset"` // not nonzero because it may be empty when stopped, maybe check
}

type DatasetMeta struct {
	Name string   `json:"name" validate:"nonzero"`
	Keys []string `json:"keys" validate:"nonzero"`
}
