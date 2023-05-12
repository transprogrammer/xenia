package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/aws/jsii-runtime-go"
)

const ConfigFile = "Config.json"

var Cfg Config = makeConfig()
var Ids Id = makeIds()

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

type Id struct {
	AzureRMProvider                *string
	NamingModule                   *string
	ResourceGroup                  *string
	VirtualNetwork                 *string
	Subnet                         *string
	PublicIPAddress                *string
	NetworkInterface               *string
	NetworkInterfaceNSGAssociation *string
	NetworkInterfaceASGAssociation *string
	NetworkSecurityGroup           *string
	ApplicationSecurityGroup       *string
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

func makeIds() Id {
	return Id{
		AzureRMProvider:                jsii.String("azurerm"),
		NamingModule:                   jsii.String("naming"),
		ResourceGroup:                  jsii.String("resource_group"),
		VirtualNetwork:                 jsii.String("virtual_network"),
		Subnet:                         jsii.String("subnet"),
		PublicIPAddress:                jsii.String("public_ip_address"),
		NetworkInterface:               jsii.String("network_interface"),
		NetworkInterfaceNSGAssociation: jsii.String("network_interface_nsg_association"),
		NetworkInterfaceASGAssociation: jsii.String("network_interface_asg_association"),
		NetworkSecurityGroup:           jsii.String("network_security_group"),
		ApplicationSecurityGroup:       jsii.String("application_security_group"),
	}
}
