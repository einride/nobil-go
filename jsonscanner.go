package nobil

import (
	"encoding/json"
	"io"
)

type JSONScanner struct {
	d                   *json.Decoder
	hasScannedToFirst   bool
	err                 error
	jsonChargingStation jsonChargingStation
}

func NewJSONScanner(r io.Reader) *JSONScanner {
	return &JSONScanner{
		d: json.NewDecoder(r),
	}
}

func (s *JSONScanner) Scan() bool {
	if s.err != nil {
		return false
	}
	if !s.hasScannedToFirst {
		if !s.scanToFirst() {
			return false
		}
	}
	if !s.d.More() {
		return false
	}
	s.err = s.d.Decode(&s.jsonChargingStation)
	return s.err == nil
}

func (s *JSONScanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

func (s *JSONScanner) ChargingStation() *ChargingStation {
	return s.jsonChargingStation.ChargingStation()
}

func (s *JSONScanner) scanToFirst() bool {
	s.hasScannedToFirst = true
	for {
		token, err := s.d.Token()
		if err != nil {
			s.err = err
			break
		}
		if delim, ok := token.(json.Delim); ok && delim == '[' {
			return true
		}
	}
	return false
}
