package config

import (
	"github.com/aws/jsii-runtime-go"
)

var Id Ids = makeIds()

type Ids struct {
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

func makeIds() Ids {
	return Ids{
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
