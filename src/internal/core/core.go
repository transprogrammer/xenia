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
	VirtualNetwork vnet.VirtualNetwork
}

const JumpboxIndex = 0
const MongoDBIndex = 1

func MakeCoreStack(scope constructs.Construct) CoreStack {
	stack := tf.NewTerraformStack(scope, x.Stacks.Core)

	NewAzureRMProvider(stack)

	naming := n.NewNamingModule(stack, x.Names.Core)
	mongoDBNaming := n.NewNamingModule(stack, x.Names.MongoDB)
	jumpboxNaming := n.NewNamingModule(stack, x.Names.Jumpbox)

	resourceGroup := NewResourceGroup(stack, naming)

	jumpboxASG := NewASG(stack, jumpboxNaming, resourceGroup)
	jumpboxNSG := NewNSG(stack, jumpboxNaming, resourceGroup, jumpboxASG)

	subnetInputs := make([]vnet.VirtualNetworkSubnet, 2)

	jumpboxSubnetInput := NewSubnetInput(stack, jumpboxNaming, jumpboxNSG, x.Config.Subnets.Jumpbox)
	subnetInputs[JumpboxIndex] = jumpboxSubnetInput

	mongoDBSubnetInput := NewSubnetInput(stack, mongoDBNaming, nil, x.Config.Subnets.MongoDB)
	subnetInputs[MongoDBIndex] = mongoDBSubnetInput

	vnet := NewVNet(stack, naming, resourceGroup, subnetInputs)

	jumpboxSubnet := GetSubnet(vnet, JumpboxIndex)
	// mongoDBSubnet := GetSubnet(vnet, MongoDBIndex)

	jumpboxIP := NewIP(stack, jumpboxNaming, resourceGroup)
	NewNIC(stack, jumpboxNaming, resourceGroup, jumpboxSubnet, jumpboxASG, jumpboxIP)

	privateDNSZone := NewPrivateDNSZone(stack, resourceGroup)
	NewDNSZoneVNetLink(stack, naming, resourceGroup, privateDNSZone, vnet)

	return CoreStack{
		MongoDBNaming:  mongoDBNaming,
		JumpboxNaming:  jumpboxNaming,
		VirtualNetwork: vnet,
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

func NewASG(stack tf.TerraformStack, naming n.NamingModule, resourceGroup rg.ResourceGroup) asg.ApplicationSecurityGroup {
	input := asg.ApplicationSecurityGroupConfig{
		Name:              naming.ApplicationSecurityGroupOutput(),
		Location:          x.Config.Regions.Primary,
		ResourceGroupName: resourceGroup.Name(),
	}

	return asg.NewApplicationSecurityGroup(stack, x.Ids.ApplicationSecurityGroup, &input)
}

func NewNSG(stack tf.TerraformStack, naming n.NamingModule, resourceGroup rg.ResourceGroup, asg asg.ApplicationSecurityGroup) nsg.NetworkSecurityGroup {
	input := nsg.NetworkSecurityGroupConfig{
		Name:              naming.NetworkSecurityGroupOutput(),
		Location:          x.Config.Regions.Primary,
		ResourceGroupName: resourceGroup.Name(),
		SecurityRule: nsg.NetworkSecurityGroupSecurityRule{
			Name:                                   jsii.String("SSH"),
			Description:                            jsii.String("Allow SSH"),
			Priority:                               jsii.Number(100),
			Direction:                              jsii.String("Inbound"),
			Access:                                 jsii.String("Allow"),
			Protocol:                               jsii.String("Tcp"),
			SourcePortRange:                        jsii.String("*"),
			DestinationPortRange:                   jsii.String("22"),
			SourceAddressPrefix:                    jsii.String("*"),
			DestinationAddressPrefix:               jsii.String("*"),
			DestinationApplicationSecurityGroupIds: &[]*string{asg.Id()},
		},
	}

	return nsg.NewNetworkSecurityGroup(stack, x.Ids.NetworkSecurityGroup, &input)
}

func NewSubnetInput(stack tf.TerraformStack, naming n.NamingModule, networkSecurityGroup nsg.NetworkSecurityGroup, addressPrefix *string) vnet.VirtualNetworkSubnet {
	return vnet.VirtualNetworkSubnet{
		Name:          naming.SubnetOutput(),
		AddressPrefix: addressPrefix,
		SecurityGroup: networkSecurityGroup.Id(),
	}
}

func NewIP(stack tf.TerraformStack, naming n.NamingModule, group rg.ResourceGroup) ip.PublicIp {
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

func NewVNet(stack tf.TerraformStack, naming n.NamingModule, resourceGroup rg.ResourceGroup, subnets []vnet.VirtualNetworkSubnet) vnet.VirtualNetwork {
	input := vnet.VirtualNetworkConfig{
		Name:              naming.VirtualNetworkOutput(),
		AddressSpace:      x.Config.AddressSpace,
		Location:          resourceGroup.Location(),
		ResourceGroupName: resourceGroup.Name(),
		Subnet:            subnets,
	}

	return vnet.NewVirtualNetwork(stack, x.Ids.VirtualNetwork, &input)
}

type VNet struct {
	vnet.VirtualNetwork
}

// NOTE: Wrap VNet to provide access to subnets by index. <>
func GetSubnet(vnet vnet.VirtualNetwork, index float64) vnet.VirtualNetworkSubnetOutputReference {
	return vnet.Subnet().Get(&index)
}

func (vnet VNet) VirtualMachineSubnet() vnet.VirtualNetworkSubnetOutputReference {
	index := float64(JumpboxIndex)
	return vnet.Subnet().Get(&index)
}

func NewNIC(stack tf.TerraformStack, naming n.NamingModule, resourceGroup rg.ResourceGroup, subnet vnet.VirtualNetworkSubnetOutputReference, asg asg.ApplicationSecurityGroup, ip ip.PublicIp) nic.NetworkInterface {
	input := nic.NetworkInterfaceConfig{
		Name:              naming.NetworkInterfaceOutput(),
		Location:          x.Config.Regions.Primary,
		ResourceGroupName: resourceGroup.Name(),

		IpConfiguration: nic.NetworkInterfaceIpConfiguration{
			Name:              jsii.String("ipconfig"),
			Primary:           jsii.Bool(true),
			SubnetId:          subnet.Id(),
			PublicIpAddressId: ip.Id(),
		},
	}
	nic := nic.NewNetworkInterface(stack, x.Ids.NetworkInterface, &input)

	asgInput := nicasg.NetworkInterfaceApplicationSecurityGroupAssociationConfig{
		NetworkInterfaceId:         nic.Id(),
		ApplicationSecurityGroupId: asg.Id(),
	}
	nicasg.NewNetworkInterfaceApplicationSecurityGroupAssociation(stack, x.Ids.NetworkInterfaceASGAssociation, &asgInput)

	nsgInput := nicnsg.NetworkInterfaceSecurityGroupAssociationConfig{
		NetworkInterfaceId:     nic.Id(),
		NetworkSecurityGroupId: nsg.Id(),
	}
	nicnsg.NewNetworkInterfaceSecurityGroupAssociation(stack, x.Ids.NetworkInterfaceNSGAssociation, &nsgInput)

	return nic
}

func NewPrivateDNSZone(stack tf.TerraformStack, resourceGroup rg.ResourceGroup) dns.PrivateDnsZone {
	input := dns.PrivateDnsZoneConfig{
		Name:              jsii.String("privatelink.mongo.cosmos.azure.com"),
		ResourceGroupName: resourceGroup.Name(),
	}

	return dns.NewPrivateDnsZone(stack, x.Ids.PrivateDNSZone, &input)
}

func NewDNSZoneVNetLink(stack tf.TerraformStack, naming n.NamingModule, resourceGroup rg.ResourceGroup, privateDnsZone dns.PrivateDnsZone, vnet vnet.VirtualNetwork) dnsl.PrivateDnsZoneVirtualNetworkLink {
	name := fmt.Sprintf("%-vnetlink", naming.PrivateDnsZoneOutput())

	input := dnsl.PrivateDnsZoneVirtualNetworkLinkConfig{
		Name:                &name,
		ResourceGroupName:   resourceGroup.Name(),
		PrivateDnsZoneName:  privateDnsZone.Name(),
		VirtualNetworkId:    vnet.Id(),
		RegistrationEnabled: jsii.Bool(true),
	}

	return dnsl.NewPrivateDnsZoneVirtualNetworkLink(stack, x.Ids.PrivateDNSZoneVirtualNetworkLink, &input)
}
