package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"kafka/pkg/models"
	"sync"
)

var configs *models.Configs
var once = &sync.Mutex{}

func Load() *models.Configs {
	if configs == nil {
		once.Lock()
		defer once.Unlock()
		if configs == nil {
			configs = loadFromFile()
		}
	}
	return configs
}

func loadFromFile() *models.Configs {
	loadedConfigs := &models.Configs{}

	file, err := ioutil.ReadFile("configs/config.yaml")
	if err != nil {
		fmt.Println("load config file error")
	}

	err = yaml.Unmarshal(file, loadedConfigs)

	if err != nil {
		fmt.Println("unmarshall config file error")
	}
	return loadedConfigs
}
