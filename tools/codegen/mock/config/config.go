package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config map[string]Item

type Item struct {
	Exclude         bool         `yaml:"exclude"`
	InferParams     []InferParam `yaml:"infer_params,omitempty"`
	URLParams       []URLParam   `yaml:"url_params,omitempty"`
	JSONBody        []JSONField  `yaml:"json_body,omitempty"`
	WrappedResponse string       `yaml:"wrapped_response,omitempty"`
}

func (i *Item) InferParam(param string) *InferParam {
	for _, ip := range i.InferParams {
		if ip.Infer == param {
			return &ip
		}
	}

	return nil
}

type (
	InferParam struct {
		Infer string `yaml:"infer"`
		From  string `yaml:"from"`
	}

	URLParam struct {
		Param     string `yaml:"param,omitempty"`
		Name      string `yaml:"name,omitempty"`
		Type      string `yaml:"type,omitempty"`
		Omitemtpy bool   `yaml:"omitemtpy,omitempty"`
		Value     string `yaml:"val,omitempty"`
	}

	JSONField struct {
		Param     string `yaml:"param,omitempty"`
		Name      string `yaml:"name,omitempty"`
		Type      string `yaml:"type,omitempty"`
		Omitemtpy bool   `yaml:"omitemtpy,omitempty"`
		Value     string `yaml:"val,omitempty"`
	}
)

var C Config

func init() {
	cdata, err := os.Open(os.Args[1])
	if err != nil {
		panic("config: " + err.Error())
	}

	defer func() {
		if err := cdata.Close(); err != nil {
			log.Println("config:", err.Error())
		}
	}()

	if err := yaml.NewDecoder(cdata).Decode(&C); err != nil {
		panic("config: " + err.Error())
	}
}
