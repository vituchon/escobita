// Routines for working config file located within a host OS filesystem
package services

import (
	"encoding/json"
	"io/ioutil"

	"local/escobita/model"
)

func LoadConfig(path string) (config model.Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}
	return
}
