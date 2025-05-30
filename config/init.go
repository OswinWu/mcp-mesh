package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var once sync.Once

func Init(filePath string) error {
	once.Do(func() {
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			panic(err)
		}
	})
	return nil
}
