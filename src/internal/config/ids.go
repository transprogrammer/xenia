package config

import "github.com/aws/jsii-runtime-go"

type ids struct {
	azurerm_provider  *string
	naming_module     *string
	resource_group    *string
	virtual_network   *string
	subnet            *string
	public_ip_address *string
}

func MakeIds() *ids {
	return &ids{
		azurerm_provider:  jsii.String("azurerm"),
		naming_module:     jsii.String("naming"),
		resource_group:    jsii.String("resource_group"),
		virtual_network:   jsii.String("virtual_network"),
		public_ip_address: jsii.String("public_ip_address"),
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
