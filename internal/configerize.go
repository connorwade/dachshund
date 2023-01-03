package internal

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type CrawlerVars struct {
	StarterURL     string   `yaml:"starterURL"`
	AllowedDomains []string `yaml:"allowedDomains"`
	Selectors      struct {
		GetContent []string `yaml:"get-content"`
		CheckLinks []string `yaml:"check-links"`
	} `yaml:"selectors"`
	Colly struct {
		MaxDepth         int  `yaml:"maxDepth"`
		Async            bool `yaml:"async"`
		ParallelRequests int  `yaml:"parallelRequests"`
	} `yaml:"Colly"`
}

func CreateConfigFile(cfg CrawlerVars) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalln("Configuration file could not be created: ", err)
	}

	err = os.WriteFile("dachshund.yaml", data, 0644)
	if err != nil {
		log.Fatalln("Configuration file could not be created: ", err)
	}
}
