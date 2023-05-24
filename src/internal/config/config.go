package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Regions struct {
	Primary   *string `json:"primary"`
	Secondary *string `json:"secondary"`
}

type Subnets struct {
	Jumpbox *string `json:"jumpbox"`
	MongoDB *string `json:"mongoDB"`
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

type Configuration struct {
	SubscriptionId  *string          `json:"subscriptionId"`
	ProjectName     *string          `json:"projectName"`
	Regions         *Regions         `json:"regions"`
	AddressSpace    *[]*string       `json:"addressSpace"`
	Subnets         *Subnets         `json:"subnets"`
	VirtualMachine  *VirtualMachine  `json:"virtualMachine"`
	DatabaseAccount *DatabaseAccount `json:"databaseAccount"`
}

const ConfigFile = "Config.json"

var Config Configuration = MakeConfig()

func MakeConfig() Configuration {

	file, err := os.Open(ConfigFile)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var config Configuration
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	return config
}
