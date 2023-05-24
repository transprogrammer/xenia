package core

import (
	"fmt"

	p "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	gn "github.com/transprogrammer/xenia/generated/naming"
	x "github.com/transprogrammer/xenia/internal/config"
	n "github.com/transprogrammer/xenia/internal/naming"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	asg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/applicationsecuritygroup"
	nic "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterface"
	nicasg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterfaceapplicationsecuritygroupassociation"
	nicnsg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterfacesecuritygroupassociation"
	nsg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networksecuritygroup"
	dns "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszone"
	dnsl "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszonevirtualnetworklink"
	ip "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/publicip"
	rg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/resourcegroup"
	vnet "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualnetwork"
	tf "github.com/hashicorp/terraform-cdk-go/cdktf"
)

type CoreStack struct {
	MongoDBNaming  gn.Naming
	JumpboxNaming  gn.Naming
	PublicIP       ip.PublicIp
	VirtualNetwork vnet.VirtualNetwork
}

func MakeCoreStack(scope constructs.Construct) CoreStack {
	stackName := fmt.Sprintf("%s-core", *x.Config.ProjectName)

	stack := tf.NewTerraformStack(scope, &stackName)

	NewAzureRMProvider(stack)

	naming := n.NewNamingModule(stack, []string{"xenia", "core"})
	mongoDBNaming := n.NewNamingModule(stack, []string{"xenia", "mongodb"})
	jumpboxNaming := n.NewNamingModule(stack, []string{"xenia", "jumpbox"})

	resourceGroup := NewResourceGroup(stack, naming)

	publicIP := NewPublicIP(stack, naming, resourceGroup)
	virtualNetwork := VirtualNetwork{NewVirtualNetwork(stack, naming, mongoDBNaming, jumpboxNaming, resourceGroup)}

	privateDNSZone := NewPrivateDNSZone(stack, resourceGroup)
	NewDNSZoneVNetLink(stack, naming, resourceGroup, privateDNSZone, virtualNetwork)

	return CoreStack{
		MongoDBNaming:  mongoDBNaming,
		JumpboxNaming:  jumpboxNaming,
		PublicIP:       publicIP,
		VirtualNetwork: virtualNetwork,
	}
}

func NewAzureRMProvider(stack tf.TerraformStack) p.AzurermProvider {
	input := p.AzurermProviderConfig{
		Features:       &p.AzurermProviderFeatures{},
		SubscriptionId: x.Config.SubscriptionId,
	}

	return p.NewAzurermProvider(stack, x.Ids.AzureRMProvider, &input)
}

func NewResourceGroup(stack tf.TerraformStack, naming n.NamingModule) rg.ResourceGroup {
	input := rg.ResourceGroupConfig{
		Name:     naming.ResourceGroupOutput(),
		Location: x.Config.Regions.Primary,
	}

	return rg.NewResourceGroup(stack, x.Ids.ResourceGroup, &input)
}

func NewPublicIP(stack tf.TerraformStack, naming n.NamingModule, group rg.ResourceGroup) ip.PublicIp {
	input := ip.PublicIpConfig{
		Name:                 naming.PublicIpOutput(),
		Location:             x.Config.Regions.Primary,
		ResourceGroupName:    group.Name(),
		Sku:                  jsii.String("Basic"),
		AllocationMethod:     jsii.String("Dynamic"),
		IpVersion:            jsii.String("IPv4"),
		DomainNameLabel:      x.Config.ProjectName,
		IdleTimeoutInMinutes: jsii.Number(4),
	}

	return ip.NewPublicIp(stack, x.Ids.PublicIPAddress, &input)
}

// HACK: Inline subnets too enable updating in-place. <>
// SEE: https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/deployment-modes#incremental-mode <>

const MongoDBSubnetIndex = 0
const VirtualMachineSubnetIndex = 1

func NewVirtualNetwork(stack tf.TerraformStack, naming n.NamingModule, mongoDBNaming n.NamingModule, jumpboxNaming n.NamingModule, resourceGroup rg.ResourceGroup) vnet.VirtualNetwork {

	subnets := []vnet.VirtualNetworkSubnet{}
	subnets[MongoDBSubnetIndex] = vnet.VirtualNetworkSubnet{
		Name:          mongoDBNaming.SubnetOutput(),
		AddressPrefix: x.Config.Subnets.MongoDB.AddressPrefix,
	}

	subnets[VirtualMachineSubnetIndex] = vnet.VirtualNetworkSubnet{
		Name:          jumpboxNaming.SubnetOutput(),
		AddressPrefix: x.Config.Subnets.VirtualMachine.AddressPrefix,
	}

	input := vnet.VirtualNetworkConfig{
		Name:              naming.VirtualNetworkOutput(),
		AddressSpace:      x.Config.AddressSpace,
		Location:          resourceGroup.Location(),
		ResourceGroupName: resourceGroup.Name(),
		Subnet:            subnets,
	}

	return vnet.NewVirtualNetwork(stack, x.Ids.VirtualNetwork, &input)
}

type VirtualNetwork struct {
	vnet.VirtualNetwork
}

func (virtualNetwork VirtualNetwork) mongoDBSubnet() vnet.VirtualNetworkSubnetOutputReference {
	index := float64(MongoDBSubnetIndex)
	return virtualNetwork.Subnet().Get(&index)
}

func (virtualNetwork VirtualNetwork) VirtualMachineSubnet() vnet.VirtualNetworkSubnetOutputReference {
	index := float64(VirtualMachineSubnetIndex)
	return virtualNetwork.Subnet().Get(&index)
}

func NewNetworkSecurityGroup(stack tf.TerraformStack, naming n.NamingModule, resourceGroup rg.ResourceGroup) nsg.NetworkSecurityGroup {
	input := nsg.NetworkSecurityGroupConfig{
		Name:              naming.NetworkSecurityGroupOutput(),
		Location:          x.Config.Regions.Primary,
		ResourceGroupName: resourceGroup.Name(),
		SecurityRule: nsg.NetworkSecurityGroupSecurityRule{
			Name:                     jsii.String("SSH"),
			Description:              jsii.String("Allow SSH"),
			Priority:                 jsii.Number(100),
			Direction:                jsii.String("Inbound"),
			Access:                   jsii.String("Allow"),
			Protocol:                 jsii.String("Tcp"),
			SourcePortRange:          jsii.String("*"),
			DestinationPortRange:     jsii.String("22"),
			SourceAddressPrefix:      jsii.String("*"),
			DestinationAddressPrefix: jsii.String("*"),
		},
	}

	return nsg.NewNetworkSecurityGroup(stack, x.Ids.NetworkSecurityGroup, &input)
}

func NewApplicationSecurityGroup(stack tf.TerraformStack, naming n.NamingModule, resourceGroup rg.ResourceGroup) asg.ApplicationSecurityGroup {
	input := asg.ApplicationSecurityGroupConfig{
		Name:              naming.ApplicationSecurityGroupOutput(),
		Location:          x.Config.Regions.Primary,
		ResourceGroupName: resourceGroup.Name(),
	}

	return asg.NewApplicationSecurityGroup(stack, x.Ids.ApplicationSecurityGroup, &input)
}

func NewNetworkInterface(stack tf.TerraformStack, naming n.NamingModule, resourceGroup rg.ResourceGroup, virtualNetwork VirtualNetwork, publicIp ip.PublicIp) nic.NetworkInterface {
	input := nic.NetworkInterfaceConfig{
		Name:              naming.NetworkInterfaceOutput(),
		Location:          x.Config.Regions.Primary,
		ResourceGroupName: resourceGroup.Name(),

		IpConfiguration: nic.NetworkInterfaceIpConfiguration{
			Name:              jsii.String("ipconfig"),
			Primary:           jsii.Bool(true),
			SubnetId:          virtualNetwork.VirtualMachineSubnet().Id(),
			PublicIpAddressId: publicIp.Id(),
		},
	}

	return nic.NewNetworkInterface(stack, x.Ids.NetworkInterface, &input)
}

func NewNICASGAssocation(stack tf.TerraformStack, networkInterface nic.NetworkInterface, applicationSecurityGroup asg.ApplicationSecurityGroup) nicasg.NetworkInterfaceApplicationSecurityGroupAssociation {
	input := nicasg.NetworkInterfaceApplicationSecurityGroupAssociationConfig{
		NetworkInterfaceId:         networkInterface.Id(),
		ApplicationSecurityGroupId: applicationSecurityGroup.Id(),
	}

	return nicasg.NewNetworkInterfaceApplicationSecurityGroupAssociation(stack, x.Ids.NetworkInterfaceASGAssociation, &input)
}

func NewNICNSGAssocation(stack tf.TerraformStack, networkInterface nic.NetworkInterface, networkSecurityGroup nsg.NetworkSecurityGroup) nicnsg.NetworkInterfaceSecurityGroupAssociation {
	input := nicnsg.NetworkInterfaceSecurityGroupAssociationConfig{
		NetworkInterfaceId:     networkInterface.Id(),
		NetworkSecurityGroupId: networkSecurityGroup.Id(),
	}

	return nicnsg.NewNetworkInterfaceSecurityGroupAssociation(stack, x.Ids.NetworkInterfaceNSGAssociation, &input)
}

func NewPrivateDNSZone(stack tf.TerraformStack, resourceGroup rg.ResourceGroup) dns.PrivateDnsZone {
	input := dns.PrivateDnsZoneConfig{
		Name:              jsii.String("privatelink.mongo.cosmos.azure.com"),
		ResourceGroupName: resourceGroup.Name(),
	}

	return dns.NewPrivateDnsZone(stack, x.Ids.PrivateDNSZone, &input)
}

func NewDNSZoneVNetLink(stack tf.TerraformStack, naming n.NamingModule, resourceGroup rg.ResourceGroup, privateDnsZone dns.PrivateDnsZone, virtualNetwork vnet.VirtualNetwork) dnsl.PrivateDnsZoneVirtualNetworkLink {
	name := fmt.Sprintf("%-vnetlink", naming.PrivateDnsZoneOutput())

	input := dnsl.PrivateDnsZoneVirtualNetworkLinkConfig{
		Name:                &name,
		ResourceGroupName:   resourceGroup.Name(),
		PrivateDnsZoneName:  privateDnsZone.Name(),
		VirtualNetworkId:    virtualNetwork.Id(),
		RegistrationEnabled: jsii.Bool(true),
	}

	return dnsl.NewPrivateDnsZoneVirtualNetworkLink(stack, x.Ids.PrivateDNSZoneVirtualNetworkLink, &input)
}
