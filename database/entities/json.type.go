package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//this struct for model type jsonb
type TypeJson []map[string]interface{}

func (c TypeJson) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *TypeJson) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var containers []interface{}
	if err := json.Unmarshal(source, &containers); err != nil {
		return err
	}

	for index := range containers {
		t, ok := containers[index].(map[string]interface{})
		if ok {
			*c = append(*c, t)
		} else {
			return errors.New("Type assertion .(map[string]interface{}) failed.")
		}
	}
	return nil
}
