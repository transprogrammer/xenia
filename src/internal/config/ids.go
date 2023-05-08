package config

import "github.com/aws/jsii-runtime-go"

type Ids struct {
	azurerm_provider                  *string
	naming_module                     *string
	resource_group                    *string
	virtual_network                   *string
	subnet                            *string
	public_ip_address                 *string
	network_interface                 *string
	network_interface_nsg_association *string
	network_interface_asg_association *string
	network_security_group            *string
	application_security_group        *string
}

func makeIds() *Ids {
	return &Ids{
		azurerm_provider:                  jsii.String("azurerm"),
		naming_module:                     jsii.String("naming"),
		resource_group:                    jsii.String("resource_group"),
		virtual_network:                   jsii.String("virtual_network"),
		public_ip_address:                 jsii.String("public_ip_address"),
		network_interface:                 jsii.String("network_interface"),
		network_interface_nsg_association: jsii.String("network_interface_nsg_association"),
		network_interface_asg_association: jsii.String("network_interface_asg_association"),
		network_security_group:            jsii.String("network_security_group"),
		application_security_group:        jsii.String("application_security_group"),
	}
}

func (i *Ids) AzureRMProvider() *string {
	return i.azurerm_provider
}

func (i *Ids) NamingModule() *string {
	return i.naming_module
}

func (i *Ids) ResourceGroup() *string {
	return i.resource_group
}

func (i *Ids) VirtualNetwork() *string {
	return i.virtual_network
}

func (i *Ids) Subnet() *string {
	return i.subnet
}

func (i *Ids) PublicIPAddress() *string {
	return i.public_ip_address
}

func (i *Ids) NetworkInterface() *string {
	return i.network_interface
}

func (i *Ids) NetworkInterfaceNSGAssociation() *string {
	return i.network_interface_nsg_association
}

func (i *Ids) NetworkInterfaceASGAssociation() *string {
	return i.network_interface_asg_association
}

func (i *Ids) NetworkSecurityGroup() *string {
	return i.network_security_group
}

func (i *Ids) ApplicationSecurityGroup() *string {
	return i.application_security_group
}
