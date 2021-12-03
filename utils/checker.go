package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Action struct {
	Actions []string `yaml:"actions"`
}

func ActionExists(ActionName string) bool {

	// fetching actions.yml file for available actions
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = filepath.Join(path, "/actions.yml")
	blob, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("wrong filepath: %s", path)
	}
	// arr will contain all the actions available
	arr := new(Action)
	if err = yaml.Unmarshal(blob, arr); err != nil {
		panic("unmarshalling error: `actions.yml`")
	}

	for _, a := range arr.Actions {
		if a == ActionName {
			return true
		}
	}

	return false
}
