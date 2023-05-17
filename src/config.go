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

type Id struct {
	ApplicationSecurityGroup         *string
	AzureRMProvider                  *string
	CosmosDBAccount                  *string
	NamingModule                     *string
	NetworkInterface                 *string
	NetworkInterfaceASGAssociation   *string
	NetworkInterfaceNSGAssociation   *string
	NetworkSecurityGroup             *string
	PrivateDNSZone                   *string
	PrivateDNSZoneGroup              *string
	PrivateDNSZoneVirtualNetworkLink *string
	PrivateEndpoint                  *string
	PublicIPAddress                  *string
	ResourceGroup                    *string
	Subnet                           *string
	VirtualMachine                   *string
	VirtualNetwork                   *string
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
		ApplicationSecurityGroup:         jsii.String("application_security_group"),
		AzureRMProvider:                  jsii.String("azurerm"),
		CosmosDBAccount:                  jsii.String("cosmosdb_account"),
		NamingModule:                     jsii.String("naming"),
		NetworkInterface:                 jsii.String("network_interface"),
		NetworkInterfaceASGAssociation:   jsii.String("network_interface_asg_association"),
		NetworkInterfaceNSGAssociation:   jsii.String("network_interface_nsg_association"),
		NetworkSecurityGroup:             jsii.String("network_security_group"),
		PrivateDNSZone:                   jsii.String("private_dns_zone"),
		PrivateDNSZoneGroup:              jsii.String("private_dns_zone_group"),
		PrivateDNSZoneVirtualNetworkLink: jsii.String("private_dns_zone_virtual_network_link"),
		PrivateEndpoint:                  jsii.String("private_endpoint"),
		PublicIPAddress:                  jsii.String("public_ip_address"),
		ResourceGroup:                    jsii.String("resource_group"),
		Subnet:                           jsii.String("subnet"),
		VirtualMachine:                   jsii.String("virtual_machine"),
		VirtualNetwork:                   jsii.String("virtual_network"),
	}

}
