package core

import (
	"fmt"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/internal/config"
	"github.com/transprogrammer/xenia/internal/naming"
	"github.com/transprogrammer/xenia/internal/stack"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/applicationsecuritygroup"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networksecuritygroup"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszone"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszonevirtualnetworklink"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/publicip"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/resourcegroup"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualnetwork"
)

type CoreStack struct {
	TerraformStack                  cdktf.TerraformStack
	MongoDBNamingModule             naming.NamingModule
	MongoDBSubnet                   virtualnetwork.VirtualNetworkSubnet
	JumpboxNamingModule             naming.NamingModule
	JumpboxSubnet                   virtualnetwork.VirtualNetworkSubnet
	JumpboxApplicationSecurityGroup applicationsecuritygroup.ApplicationSecurityGroup
	JumpboxNetworkSecurityGroup     networksecuritygroup.NetworkSecurityGroup
}

func (s CoreStack) Stack() cdktf.TerraformStack {
	return s.TerraformStack
}

const (
	JumpboxIndex = iota
	MongoDBIndex
)

func NewStack(
	scope constructs.Construct,
) CoreStack {

	coreStack := cdktf.NewTerraformStack(scope, config.Stacks.Core)

	stack.NewAzureRMProvider(coreStack)

	namingModule := naming.NewNamingModule(coreStack, config.Names.Core)
	mongoDBNamingModule := naming.NewNamingModule(coreStack, config.Names.MongoDB)
	jumpboxNamingModule := naming.NewNamingModule(coreStack, config.Names.Jumpbox)

	resourceGroup := stack.NewResourceGroup(coreStack, namingModule)

	jumpboxApplicationSecurityGroup := NewApplicationSecurityGroup(coreStack, jumpboxNamingModule, resourceGroup)
	jumpboxNetworkSecurityGroup := NewNetworkSecurityGroup(coreStack, jumpboxNamingModule, resourceGroup, jumpboxApplicationSecurityGroup)

	subnetInputs := make([]virtualnetwork.VirtualNetworkSubnet, 2)

	jumpboxSubnetInput := NewSubnetInput(coreStack, jumpboxNamingModule, jumpboxNetworkSecurityGroup, config.Config.Subnets.Jumpbox)
	subnetInputs[JumpboxIndex] = jumpboxSubnetInput

	mongoDBSubnetInput := NewSubnetInput(coreStack, mongoDBNamingModule, nil, config.Config.Subnets.MongoDB)
	subnetInputs[MongoDBIndex] = mongoDBSubnetInput

	virtualNetwork := NewVirtualNetwork(coreStack, namingModule, resourceGroup, subnetInputs)

	jumpboxSubnet := GetSubnet(virtualNetwork, JumpboxIndex)
	mongoDBSubnet := GetSubnet(virtualNetwork, MongoDBIndex)

	privateDNSZone := NewPrivateDNSZone(coreStack, resourceGroup)
	NewDNSZoneVNetLink(coreStack, naming, resourceGroup, privateDNSZone, vnet)

	return CoreStack{
		TerraformStack:                  coreStack,
		MongoDBNamingModule:             mongoDBNamingModule,
		MongoDBSubnet:                   mongoDBSubnet,
		JumpboxNaming:                   jumpboxNamingModule,
		JumpboxSubnet:                   jumpboxSubnet,
		JumpboxApplicationSecurityGroup: jumpboxApplicationSecurityGroup,
		JumpboxNetworkSecurityGroup:     jumpboxNetworkSecurityGroup,
		VirtualNetwork:                  vnet,
	}
}

func NewApplicationSecurityGroup(
	stack cdktf.TerraformStack,
	naming naming.NamingModule,
	resourceGroup resourcegroup.ResourceGroup,
) applicationsecuritygroup.ApplicationSecurityGroup {

	input := applicationsecuritygroup.ApplicationSecurityGroupConfig{
		Name:              naming.ApplicationSecurityGroupOutput(),
		Location:          config.Config.Regions.Primary,
		ResourceGroupName: resourceGroup.Name(),
	}

	return applicationsecuritygroup.NewApplicationSecurityGroup(
		stack,
		config.Ids.ApplicationSecurityGroup,
		&input,
	)
}

func NewNetworkSecurityGroup(
	stack cdktf.TerraformStack,
	naming naming.NamingModule,
	resourceGroup resourcegroup.ResourceGroup,
	asg applicationsecuritygroup.ApplicationSecurityGroup,
) networksecuritygroup.NetworkSecurityGroup {

	input := networksecuritygroup.NetworkSecurityGroupConfig{
		Name:              naming.NetworkSecurityGroupOutput(),
		Location:          config.Config.Regions.Primary,
		ResourceGroupName: resourceGroup.Name(),
		SecurityRule: networksecuritygroup.NetworkSecurityGroupSecurityRule{
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

	return networksecuritygroup.NewNetworkSecurityGroup(
		stack,
		config.Ids.NetworkSecurityGroup,
		&input,
	)
}

func NewSubnetInput(
	stack cdktf.TerraformStack,
	naming naming.NamingModule,
	networkSecurityGroup networksecuritygroup.NetworkSecurityGroup,
	addressPrefix *string,
) virtualnetwork.VirtualNetworkSubnet {

	return virtualnetwork.VirtualNetworkSubnet{
		Name:          naming.SubnetOutput(),
		AddressPrefix: addressPrefix,
		SecurityGroup: networkSecurityGroup.Id(),
	}
}

func NewPublicIP(
	stack cdktf.TerraformStack,
	naming naming.NamingModule,
	group resourcegroup.ResourceGroup,
) publicip.PublicIp {

	input := publicip.PublicIpConfig{
		Name:                 naming.PublicIpOutput(),
		Location:             config.Config.Regions.Primary,
		ResourceGroupName:    group.Name(),
		Sku:                  jsii.String("Basic"),
		AllocationMethod:     jsii.String("Dynamic"),
		IpVersion:            jsii.String("IPv4"),
		DomainNameLabel:      config.Config.ProjectName,
		IdleTimeoutInMinutes: jsii.Number(4),
	}

	return publicip.NewPublicIp(stack, config.Ids.PublicIPAddress, &input)
}

// HACK: Inline subnets too enable updating in-place. <>
// SEE: https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/deployment-modes#incremental-mode <>

func NewVirtualNetwork(stack cdktf.TerraformStack, naming naming.NamingModule, resourceGroup resourcegroup.ResourceGroup, subnets []virtualnetwork.VirtualNetworkSubnet) virtualnetwork.VirtualNetwork {
	input := virtualnetwork.VirtualNetworkConfig{
		Name:              naming.VirtualNetworkOutput(),
		AddressSpace:      config.Config.AddressSpace,
		Location:          resourceGroup.Location(),
		ResourceGroupName: resourceGroup.Name(),
		Subnet:            subnets,
	}

	return virtualnetwork.NewVirtualNetwork(stack, config.Ids.VirtualNetwork, &input)
}

// NOTE: Wrap VNet to provide access to subnets by indeconfig. <>
func GetSubnet(vnet virtualnetwork.VirtualNetwork, index float64) virtualnetwork.VirtualNetworkSubnetOutputReference {
	return vnet.Subnet().Get(&index)
}

func NewPrivateDNSZone(stack cdktf.TerraformStack, resourceGroup resourcegroup.ResourceGroup) privatednszone.PrivateDnsZone {
	input := privatednszone.PrivateDnsZoneConfig{
		Name:              jsii.String("privatelink.mongo.cosmos.azure.com"),
		ResourceGroupName: resourceGroup.Name(),
	}

	return privatednszone.NewPrivateDnsZone(stack, config.Ids.PrivateDNSZone, &input)
}

func NewDNSZoneVNetLink(stack cdktf.TerraformStack, naming naming.NamingModule, resourceGroup resourcegroup.ResourceGroup, privateDnsZone privatednszone.PrivateDnsZone, vnet virtualnetwork.VirtualNetwork) privatednszonevirtualnetworklink.PrivateDnsZoneVirtualNetworkLink {
	name := fmt.Sprintf("%-vnetlink", naming.PrivateDnsZoneOutput())

	input := privatednszonevirtualnetworklink.PrivateDnsZoneVirtualNetworkLinkConfig{
		Name:                &name,
		ResourceGroupName:   resourceGroup.Name(),
		PrivateDnsZoneName:  privateDnsZone.Name(),
		VirtualNetworkId:    virtualnetwork.Id(),
		RegistrationEnabled: jsii.Bool(true),
	}

	return privatednszonevirtualnetworklink.NewPrivateDnsZoneVirtualNetworkLink(stack, config.Ids.PrivateDNSZoneVirtualNetworkLink, &input)
}
