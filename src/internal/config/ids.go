package config

import "github.com/aws/jsii-runtime-go"

type ids struct {
	azurerm_provider                  *string
	naming_module                     *string
	resource_group                    *string
	virtual_network                   *string
	subnet                            *string
	public_ip_address                 *string
	network_interface                 *string
	network_interface_nsg_association *string
	network_interface_asg_association *string
}

func MakeIds() *ids {
	return &ids{
		azurerm_provider:                  jsii.String("azurerm"),
		naming_module:                     jsii.String("naming"),
		resource_group:                    jsii.String("resource_group"),
		virtual_network:                   jsii.String("virtual_network"),
		public_ip_address:                 jsii.String("public_ip_address"),
		network_interface:                 jsii.String("network_interface"),
		network_interface_nsg_association: jsii.String("network_interface_nsg_association"),
		network_interface_asg_association: jsii.String("network_interface_asg_association"),
	}
}

func (i *ids) AzureRMProvider() *string {
	return i.azurerm_provider
}

func (i *ids) NamingModule() *string {
	return i.naming_module
}

func (i *ids) ResourceGroup() *string {
	return i.resource_group
}

func (i *ids) VirtualNetwork() *string {
	return i.virtual_network
}

func (i *ids) Subnet() *string {
	return i.subnet
}

func (i *ids) PublicIPAddress() *string {
	return i.public_ip_address
}

func (i *ids) NetworkInterface() *string {
	return i.network_interface
}

func (i *ids) NetworkInterfaceNSGAssociation() *string {
	return i.network_interface_nsg_association
}

func (i *ids) NetworkInterfaceASGAssociation() *string {
	return i.network_interface_asg_association
}
