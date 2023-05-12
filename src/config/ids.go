package config

import "github.com/aws/jsii-runtime-go"

type Ids struct {
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

var ids IDS = &Ids{
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
}
