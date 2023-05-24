package config

import (
	i "github.com/aws/jsii-runtime-go"
)

type Identifiers struct {
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

var Ids Identifiers = Identifiers{
	ApplicationSecurityGroup:         i.String("application_security_group"),
	AzureRMProvider:                  i.String("azurerm"),
	CosmosDBAccount:                  i.String("cosmosdb_account"),
	NamingModule:                     i.String("naming"),
	NetworkInterface:                 i.String("network_interface"),
	NetworkInterfaceASGAssociation:   i.String("network_interface_asg_association"),
	NetworkInterfaceNSGAssociation:   i.String("network_interface_nsg_association"),
	NetworkSecurityGroup:             i.String("network_security_group"),
	PrivateDNSZone:                   i.String("private_dns_zone"),
	PrivateDNSZoneGroup:              i.String("private_dns_zone_group"),
	PrivateDNSZoneVirtualNetworkLink: i.String("private_dns_zone_virtual_network_link"),
	PrivateEndpoint:                  i.String("private_endpoint"),
	PublicIPAddress:                  i.String("public_ip_address"),
	ResourceGroup:                    i.String("resource_group"),
	Subnet:                           i.String("subnet"),
	VirtualMachine:                   i.String("virtual_machine"),
	VirtualNetwork:                   i.String("virtual_network"),
}
