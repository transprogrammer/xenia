package main

import (
	"io/ioutil"
	"os"

	"github.com/aws/jsii-runtime-go"
)

const ConfigFile = "config.json"

type ConfigJSON struct {
	SubscriptionId  string   `json:"subscriptionId"`
	ProjectName     string   `json:"projectName"`
	PrimaryRegion   string   `json:"primaryRegion"`
	SecondaryRegion string   `json:"secondaryRegion"`
	AddressSpace    []string `json:"addressSpace"`
	SubnetPrefixes  []string `json:"subnetPrefixes"`
}

type Config struct {
	projectName     *string
	subscriptionId  *string
	primaryRegion   *string
	secondaryRegion *string
	addressSpace    []*string
	subnetPrefixes  []*string
}

func makeConfig() Config {
	json := makeConfigJSON()
	
  toJsii := func(str string) interface{} {
		return jsii.String(str)
	}

	return Config{
		projectName:     jsii.String(json.ProjectName),
		subscriptionId:  jsii.String(json.SubscriptionId),
		primaryRegion:   jsii.String(json.PrimaryRegion),
		secondaryRegion: jsii.String(json.SecondaryRegion),
		addressSpace:    Map(json.AddressSpace, toJsii),
		subnetPrefixes:  Map(json.SubnetPrefixes, toJsii),
	}
}

func makeConfigJSON() ConfigJSON {
	file, err := os.Open("ConfigFile"
	if err != nil {
		panic(err)
	}
	defer file.Close()

  bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var json ConfigJSON
	err = json.Unmarshal(bytes, &ConfigJSON)
	if err != nil {
		panic(err)
	}

	return json
}

func (cfg Config) ProjectName() *string {
	return cfg.projectName
}

func (cfg Config) ProviderName() *string {
	return cfg.providerName
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

func (cfg Config) SubnetPrefixes() []*string {
	return cfg.subnetPrefixes
}

