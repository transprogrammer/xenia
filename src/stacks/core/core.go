package core

import (
	"fmt"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	. "github.com/transprogrammer/xenia/internal/config"
	. "github.com/transprogrammer/xenia/internal/naming"
	. "github.com/transprogrammer/xenia/internal/stacks"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	asg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/applicationsecuritygroup"
	nsg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networksecuritygroup"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszone"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszonevirtualnetworklink"
	rg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/resourcegroup"
	vnet "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualnetwork"
)

type CoreStack struct {
	TerraformStack cdktf.TerraformStack
	MongoDBNaming  Naming
	MongoDBSubnet  vnet.VirtualNetworkSubnetOutputReference
	JumpboxNaming  Naming

	JumpboxSubnet vnet.VirtualNetworkSubnetOutputReference
	JumpboxASG    asg.ApplicationSecurityGroup
	JumpboxNSG    nsg.NetworkSecurityGroup
}

func (stack CoreStack) Stack() cdktf.TerraformStack {
	return stack.TerraformStack
}

const (
	JumpboxIndex = iota
	MongoDBIndex
)

func NewStack(scope constructs.Construct) CoreStack {

	stack := cdktf.NewTerraformStack(scope, Stacks.Core)

	NewAzureRMProvider(stack)

	naming := NewNaming(stack, Names.Core)
	mongoDBNaming := NewNaming(stack, Names.MongoDB)
	jumpboxNaming := NewNaming(stack, Names.Jumpbox)

	rg := NewResourceGroup(stack, naming)

	jumpboxASG := NewASG(stack, jumpboxNaming, rg)
	jumpboxNSG := NewNSG(stack, jumpboxNaming, rg, jumpboxASG)

	subnetInputs := make([]vnet.VirtualNetworkSubnet, 2)

	jumpboxSubnetInput := NewSubnetInput(stack, jumpboxNaming, jumpboxNSG, Config.Subnets.Jumpbox)
	subnetInputs[JumpboxIndex] = jumpboxSubnetInput

	mongoDBSubnetInput := NewSubnetInput(stack, mongoDBNaming, nil, Config.Subnets.MongoDB)
	subnetInputs[MongoDBIndex] = mongoDBSubnetInput

	vnet := NewVNet(stack, naming, rg, subnetInputs)

	jumpboxSubnet := GetSubnet(vnet, JumpboxIndex)
	mongoDBSubnet := GetSubnet(vnet, MongoDBIndex)

	privateDNSZone := NewPrivateDNSZone(stack, rg)
	NewDNSZoneVNetLink(stack, naming, rg, privateDNSZone, vnet)

	return CoreStack{
		TerraformStack: stack,
		MongoDBNaming:  mongoDBNaming,
		MongoDBSubnet:  mongoDBSubnet,
		JumpboxNaming:  jumpboxNaming,
		JumpboxSubnet:  jumpboxSubnet,
		JumpboxASG:     jumpboxASG,
		JumpboxNSG:     jumpboxNSG,
	}
}

func NewASG(stack cdktf.TerraformStack, naming Naming, rg rg.ResourceGroup) asg.ApplicationSecurityGroup {

	input := asg.ApplicationSecurityGroupConfig{
		Name:              naming.ApplicationSecurityGroupOutput(),
		Location:          Config.Regions.Primary,
		ResourceGroupName: rg.Name(),
	}

	return asg.NewApplicationSecurityGroup(
		stack,
		Ids.ASG,
		&input,
	)
}

func NewNSG(stack cdktf.TerraformStack, naming Naming, rg rg.ResourceGroup, asg asg.ApplicationSecurityGroup) nsg.NetworkSecurityGroup {

	input := nsg.NetworkSecurityGroupConfig{
		Name:              naming.NetworkSecurityGroupOutput(),
		Location:          Config.Regions.Primary,
		ResourceGroupName: rg.Name(),
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

	return nsg.NewNetworkSecurityGroup(
		stack,
		Ids.NetworkSecurityGroup,
		&input,
	)
}

func NewSubnetInput(stack cdktf.TerraformStack, naming Naming, nsg nsg.NetworkSecurityGroup, addressPrefix *string) vnet.VirtualNetworkSubnet {

	return vnet.VirtualNetworkSubnet{
		Name:          naming.SubnetOutput(),
		AddressPrefix: addressPrefix,
		SecurityGroup: nsg.Id(),
	}
}

// HACK: Inline subnets too enable updating in-place. <>
// SEE: https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/deployment-modes#incremental-mode <>

func NewVNet(stack cdktf.TerraformStack, naming Naming, rg rg.ResourceGroup, subnets []vnet.VirtualNetworkSubnet) vnet.VirtualNetwork {
	input := vnet.VirtualNetworkConfig{
		Name:              naming.VirtualNetworkOutput(),
		AddressSpace:      Config.AddressSpace,
		Location:          rg.Location(),
		ResourceGroupName: rg.Name(),
		Subnet:            subnets,
	}

	return vnet.NewVirtualNetwork(stack, Ids.VNet, &input)
}

// NOTE: Wrap VNet to provide access to subnets by inde <>
func GetSubnet(vnet vnet.VirtualNetwork, index float64) vnet.VirtualNetworkSubnetOutputReference {
	return vnet.Subnet().Get(&index)
}

func NewPrivateDNSZone(stack cdktf.TerraformStack, rg rg.ResourceGroup) privatednszone.PrivateDnsZone {
	input := privatednszone.PrivateDnsZoneConfig{
		Name:              jsii.String("privatelink.mongo.cosmos.azure.com"),
		ResourceGroupName: rg.Name(),
	}

	return privatednszone.NewPrivateDnsZone(stack, Ids.PrivateDNSZone, &input)
}

func NewDNSZoneVNetLink(stack cdktf.TerraformStack, naming Naming, rg rg.ResourceGroup, privateDnsZone privatednszone.PrivateDnsZone, vnet vnet.VirtualNetwork) privatednszonevirtualnetworklink.PrivateDnsZoneVirtualNetworkLink {
	name := fmt.Sprintf("%-vnetlink", naming.PrivateDnsZoneOutput())

	input := privatednszonevirtualnetworklink.PrivateDnsZoneVirtualNetworkLinkConfig{
		Name:                &name,
		ResourceGroupName:   rg.Name(),
		PrivateDnsZoneName:  privateDnsZone.Name(),
		VirtualNetworkId:    vnet.Id(),
		RegistrationEnabled: jsii.Bool(true),
	}

	return privatednszonevirtualnetworklink.NewPrivateDnsZoneVirtualNetworkLink(stack, Ids.PrivateDNSZoneVirtualNetworkLink, &input)
}
