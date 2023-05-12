package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const ConfigFile = "config.json"

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

type ImageReference struct {
	Publisher *string `json:"publisher"`
	Offer     *string `json:"offer"`
	Sku       *string `json:"sku"`
	Version   *string `json:"version"`
}

type VirtualMachine struct {
	Size               *string         `json:"size"`
	StorageAccountType *string         `json:"storageAccountType"`
	ImageReference     *ImageReference `json:"imageReference"`
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
	AddressSpace    []*string        `json:"addressSpace"`
	Subnets         *Subnets         `json:"subnets"`
	VirtualMachine  *VirtualMachine  `json:"virtualMachine"`
	DatabaseAccount *DatabaseAccount `json:"databaseAccount"`
}

func MakeConfig() *Config {
	config := makeConfig()
	config := makeConfigFrom(config)
	config.ids = makeIds()

	return config
}

func makeConfig() *config {
	file, err := os.Open(ConfigFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var cfg config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func makeConfigFrom(jso *config) *Config {
	return &Config{
		subscriptionId: jso.SubscriptionId,
		projectName:    jso.ProjectName,
		regions: &Regions{
			primary:   jso.Regions.Primary,
			secondary: jso.Regions.Secondary,
		},
		addressSpace: jso.AddressSpace,
		subnets: &Subnets{
			virtualMachine: &Subnet{
				postfix:       jso.Subnets.VirtualMachine.Postfix,
				addressPrefix: jso.Subnets.VirtualMachine.AddressPrefix,
			},
			mongoDB: &Subnet{
				postfix:       jso.Subnets.MongoDB.Postfix,
				addressPrefix: jso.Subnets.MongoDB.AddressPrefix,
			},
		},
		virtualMachine: &VirtualMachine{
			size:               jso.VirtualMachine.Size,
			storageAccountType: jso.VirtualMachine.StorageAccountType,
			imageReference: &ImageReference{
				publisher: jso.VirtualMachine.ImageReference.Publisher,
				offer:     jso.VirtualMachine.ImageReference.Offer,
				sku:       jso.VirtualMachine.ImageReference.Sku,
				version:   jso.VirtualMachine.ImageReference.Version,
			},
		},
		databaseAccount: &DatabaseAccount{
			kind:                    jso.DatabaseAccount.Kind,
			serverVersion:           jso.DatabaseAccount.ServerVersion,
			offerType:               jso.DatabaseAccount.OfferType,
			backupPolicyType:        jso.DatabaseAccount.BackupPolicyType,
			defaultConsistencyLevel: jso.DatabaseAccount.DefaultConsistencyLevel,
			capabilities:            jso.DatabaseAccount.Capabilities,
		},
	}
}
