package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const ConfigFile = "Config.json"

var Cfg Config = makeConfig()

type Regions struct {
	Primary   *string `json:"primary"`
	Secondary *string `json:"secondary"`
}

type Subnet struct {
	Postfix       *string `json:"postfix"`
	AddressPrefix *string `json:"addressPrefix"`
}

type Subnets struct {
	VirtualMachine *Subnet `json:"virtualMachine"`
	MongoDB        *Subnet `json:"mongoDB"`
}

type Image struct {
	Publisher *string `json:"publisher"`
	Offer     *string `json:"offer"`
	Sku       *string `json:"sku"`
	Version   *string `json:"version"`
}

type VirtualMachine struct {
	Size               *string `json:"size"`
	StorageAccountType *string `json:"storageAccountType"`
	Image              *Image  `json:"imageReference"`
	AdminUsername      *string `json:"adminUsername"`
	SSHPublicKey       *string `json:"sshPublicKey"`
}

type DatabaseAccount struct {
	Kind                    *string   `json:"kind"`
	ServerVersion           *string   `json:"serverVersion"`
	OfferType               *string   `json:"offerType"`
	BackupPolicyType        *string   `json:"backupPolicyType"`
	DefaultConsistencyLevel *string   `json:"defaultConsistencyLevel"`
	Capabilities            []*string `json:"capabilities"`
}

type Config struct {
	SubscriptionId  *string          `json:"subscriptionId"`
	ProjectName     *string          `json:"projectName"`
	Regions         *Regions         `json:"regions"`
	AddressSpace    *[]*string       `json:"addressSpace"`
	Subnets         *Subnets         `json:"subnets"`
	VirtualMachine  *VirtualMachine  `json:"virtualMachine"`
	DatabaseAccount *DatabaseAccount `json:"databaseAccount"`
}

func makeConfig() Config {
	panik := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	file, err := os.Open(ConfigFile)
	panik(err)

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	panik(err)

	var cfg Config
	err = json.Unmarshal(bytes, &cfg)
	panik(err)

	return cfg
}
