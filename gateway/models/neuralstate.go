package models

import (
	"fmt"
	"net/http"
)

type DataMode uint8

const (
	Stopped DataMode = iota
	Collecting
	Predicting
)

type DatasetEntryRequest struct {
	Entry map[string]float64 `json:"entry" validate:"nonzero"`
}

type NeuralPredictionRequest struct {
	DatasetEntryRequest string             `json:"dataset"`
	Query               map[string]float64 `json:"query"`
}

type NeuralPredictionResponse struct {
	Situation int `json:"situation" validate:"nonzero,min=0"`
}

type NeuralState struct {
	Mode        DataMode `json:"mode" validate:"min=0,max=2"`
	Dataset     string   `json:"dataset" validate:"nonzero"`
	dataSetKeys []string `json:"-" validate:"-"` // private field for internal use only
}

type DatasetMeta struct {
	Name string   `json:"name" validate:"nonzero"`
	Keys []string `json:"keys" validate:"nonzero"`
}

func (s *NeuralState) DatasetKeysIfExists() ([]string, error) {
	if s.dataSetKeys != nil {
		return s.dataSetKeys, nil
	}

	resp, err := http.Get(fmt.Sprintf("http://dataset/%s", s.Dataset))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, NotFound
	}

	var meta DatasetMeta
	if err := ReadAndDecode(resp.Body, &meta); err != nil {
		return nil, err
	}

	s.dataSetKeys = meta.Keys
	return meta.Keys, nil
}
