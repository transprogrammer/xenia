package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/aws/jsii-runtime-go"
)

const ConfigFile = "config.json"

type ConfigJSON struct {
	SubscriptionId  string            `json:"subscriptionId"`
	ProjectName     string            `json:"projectName"`
	PrimaryRegion   string            `json:"primaryRegion"`
	SecondaryRegion string            `json:"secondaryRegion"`
	AddressSpace    []string          `json:"addressSpace"`
	Subnets         map[string]string `json:"subnets"`
}

type Config struct {
	projectName     *string
	subscriptionId  *string
	primaryRegion   *string
	secondaryRegion *string
	addressSpace    []*string
	subnets         map[*string]*string
}

func sliceToJsii(strs []string) []*string {
	jsiiStrs := make([]*string, len(strs))
	for i, str := range strs {
		jsiiStrs[i] = jsii.String(str)
	}
	return jsiiStrs
}

func mapToJsii(strMap map[string]string) map[*string]*string {
	jsiiMap := make(map[*string]*string, len(strMap))
	for k, v := range strMap {
		jsiiMap[jsii.String(k)] = jsii.String(v)
	}
	return jsiiMap
}

func makeConfig() Config {
	json := makeConfigJSON()

	return Config{
		projectName:     jsii.String(json.ProjectName),
		subscriptionId:  jsii.String(json.SubscriptionId),
		primaryRegion:   jsii.String(json.PrimaryRegion),
		secondaryRegion: jsii.String(json.SecondaryRegion),
		addressSpace:    sliceToJsii(json.AddressSpace),
		subnets:         mapToJsii(json.Subnets),
	}
}

func makeConfigJSON() ConfigJSON {
	file, err := os.Open(ConfigFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var jso ConfigJSON
	err = json.Unmarshal(bytes, jso)
	if err != nil {
		panic(err)
	}

	return jso
}

func (cfg Config) ProjectName() *string {
	return cfg.projectName
}

func (cfg Config) SubscriptionId() *string {
	return cfg.subscriptionId
}

func (cfg Config) PrimaryRegion() *string {
	return cfg.primaryRegion
}

func (cfg Config) SecondaryRegion() *string {
	return cfg.secondaryRegion
}

func (cfg Config) AddressSpace() []*string {
	return cfg.addressSpace
}

func (cfg Config) Subnets() map[*string]*string {
	return cfg.subnets
}
