package main

import (
	"github.com/aws/jsii-runtime-go"
	tf "github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/generated/naming"
	cfg "github.com/transprogrammer/xenia/internal/config"

	asg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/applicationsecuritygroup"
	nic "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterface"
	nicasg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterfaceapplicationsecuritygroupassociation"
	nicnsg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterfacesecuritygroupassociation"
	nsg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networksecuritygroup"
	prov "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	ip "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/publicip"
	rg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/resourcegroup"
	vnet "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualnetwork"
)

func main() {
	makeApp().Synth()
}

func makeApp() tf.App {
	cfg := cfg.MakeConfig()
	app := tf.NewApp(nil)
	stk := tf.NewTerraformStack(app, cfg.ProjectName())

	makeAzureRMProvider(cfg, stk)

	nme := makeNamingModule(cfg, stk, []*string{})
	rg := makeResourceGroup(cfg, stk, nme)

	ip := makePublicIPAddress(cfg, stk, nme, rg)
	vnet := makeVirtualNetwork(cfg, stk, nme, rg)

	nsg := makeNetworkSecurityGroup(cfg, stk, nme, rg)
	asg := makeApplicationSecurityGroup(cfg, stk, nme, rg)

	nic := makeNetworkInterface(cfg, stk, nme, rg, vnet, ip)
	associateNICWithNSG(cfg, stk, nme, nic, nsg)
	associateNICWithASG(cfg, stk, nme, nic, asg)

	//makeVirtualMachine(cfg, stk, nme, rg, nic, vnet)

	return app
}

func makeAzureRMProvider(cfg *cfg.Config, stack tf.TerraformStack) prov.AzurermProvider {
	return prov.NewAzurermProvider(stack, cfg.Ids().AzureRMProvider(), &prov.AzurermProviderConfig{
		Features:       &prov.AzurermProviderFeatures{},
		SubscriptionId: cfg.SubscriptionId(),
	})
}

func makeNamingModule(cfg *cfg.Config, stack tf.TerraformStack, suffixes []*string) *naming.Naming {
	prefix := []*string{cfg.ProjectName()}

	namingModule := naming.NewNaming(stack, cfg.Ids().NamingModule(), &naming.NamingConfig{
		Prefix:               &prefix,
		UniqueIncludeNumbers: jsii.Bool(false),
		Suffix:               &suffixes,
	})

	return &namingModule
}

func makeResourceGroup(cfg *cfg.Config, stack tf.TerraformStack, naming *naming.Naming) *rg.ResourceGroup {
	resourceGroup := rg.NewResourceGroup(stack, cfg.Ids().ResourceGroup(), &rg.ResourceGroupConfig{
		Name:     (*naming).ResourceGroupOutput(),
		Location: cfg.Regions().Primary(),
	})
	return &resourceGroup
}

// ???: Inline subnet too enable updating in-place. <>
// SEE: https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/deployment-modes#incremental-mode <>
func makeVirtualNetwork(cfg *cfg.Config, stack tf.TerraformStack, naming *naming.Naming, rg *rg.ResourceGroup) *vnet.VirtualNetwork {

	subnets := []vnet.VirtualNetworkSubnet{
		vnet.VirtualNetworkSubnet{
			Name:          (*makeNamingModule(cfg, stack, []*string{cfg.Subnets().MongoDB().Postfix()})).SubnetOutput(),
			AddressPrefix: cfg.Subnets().MongoDB().AddressPrefix(),
		},
		vnet.VirtualNetworkSubnet{
			Name:          (*makeNamingModule(cfg, stack, []*string{cfg.Subnets().MongoDB().Postfix()})).SubnetOutput(),
			AddressPrefix: cfg.Subnets().VirtualMachine().AddressPrefix(),
		},
	}

	addressSpace := cfg.AddressSpace()

	virtualNetwork := vnet.NewVirtualNetwork(stack, cfg.Ids().VirtualNetwork(), &vnet.VirtualNetworkConfig{
		Name:              (*naming).VirtualNetworkOutput(),
		AddressSpace:      &addressSpace,
		Location:          (*rg).Location(),
		ResourceGroupName: (*rg).Name(),

		Subnet: subnets,
	})

	return &virtualNetwork
}

func makePublicIPAddress(cfg *cfg.Config, stack tf.TerraformStack, naming *naming.Naming, rg *rg.ResourceGroup) *ip.PublicIp {
	publicIp := ip.NewPublicIp(stack, cfg.Ids().PublicIPAddress(), &ip.PublicIpConfig{
		Name:                 (*naming).PublicIpOutput(),
		Location:             (*rg).Location(),
		ResourceGroupName:    (*rg).Name(),
		Sku:                  jsii.String("Basic"),
		AllocationMethod:     jsii.String("Dynamic"),
		IpVersion:            jsii.String("IPv4"),
		DomainNameLabel:      cfg.ProjectName(),
		IdleTimeoutInMinutes: jsii.Number(4),
	})
	return &publicIp
}

func makeNetworkSecurityGroup(cfg *cfg.Config, stack tf.TerraformStack, naming *naming.Naming, rg *rg.ResourceGroup) *nsg.NetworkSecurityGroup {
	networkSecurityRule := nsg.NetworkSecurityGroupSecurityRule{
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
	}

	networkSecurityGroup := nsg.NewNetworkSecurityGroup(stack, cfg.Ids().NetworkSecurityGroup(), &nsg.NetworkSecurityGroupConfig{
		Name:              (*naming).NetworkSecurityGroupOutput(),
		Location:          (*rg).Location(),
		ResourceGroupName: (*rg).Name(),
		SecurityRule:      networkSecurityRule,
	})

	return &networkSecurityGroup
}

func makeApplicationSecurityGroup(cfg *cfg.Config, stack tf.TerraformStack, naming *naming.Naming, rg *rg.ResourceGroup) *asg.ApplicationSecurityGroup {
	applicationSecurityGroup := asg.NewApplicationSecurityGroup(stack, cfg.Ids().ApplicationSecurityGroup(), &asg.ApplicationSecurityGroupConfig{
		Name:              (*naming).ApplicationSecurityGroupOutput(),
		Location:          (*rg).Location(),
		ResourceGroupName: (*rg).Name(),
	})

	return &applicationSecurityGroup
}

func makeNetworkInterface(cfg *cfg.Config, stack tf.TerraformStack, naming *naming.Naming, rg *rg.ResourceGroup, vnet *vnet.VirtualNetwork, publicIp *ip.PublicIp) *nic.NetworkInterface {
	networkInterface := nic.NewNetworkInterface(stack, cfg.Ids().NetworkInterface(), &nic.NetworkInterfaceConfig{
		Name:              (*naming).NetworkInterfaceOutput(),
		Location:          (*rg).Location(),
		ResourceGroupName: (*rg).Name(),

		IpConfiguration: nic.NetworkInterfaceIpConfiguration{
			Name:              jsii.String("ipconfig"),
			Primary:           jsii.Bool(true),
			SubnetId:          (*vnet).Subnet().MongoDBSubnet().Id(),
			PublicIpAddressId: (*publicIp).Id(),
		},
	})

	return &networkInterface
}

func associateNICWithNSG(cfg *cfg.Config, stack tf.TerraformStack, naming *naming.Naming, nic *nic.NetworkInterface, nsg *nsg.NetworkSecurityGroup) {
	nicNSGAssociation := nicnsg.NewNetworkInterfaceSecurityGroupAssociation()(stack, cfg.Ids().NetworkInterfaceNSGAssociation(), &nicnsg.NetworkInterfaceSecurityGroupAssociationConfig{
		NetworkInterfaceId:     (*nic).Id(),
		NetworkSecurityGroupId: (*nsg).Id(),
	})

	return &nicNSGAssociation
}

func associateNICWithASG(cfg *cfg.Config, stack tf.TerraformStack, naming *naming.Naming, nic *nic.NetworkInterface, asg *asg.ApplicationSecurityGroup) {
	nicASGAssociation := nicasg.NewNetworkInterfaceApplicationSecurityGroupAssociation()(stack, cfg.Ids().NetworkInterfaceASGAssociation(), &nicasg.NetworkInterfaceApplicationSecurityGroupAssociationConfig{
		NetworkInterfaceId:                        (*nic).Id(),
		ApplicationSecurityGroupId:                (*asg).Id(),
		IpConfigurationName:                       jsii.String("ipconfig"),
		IpConfigurationPrivateIpAddressAllocation: jsii.String("Dynamic"),
	})

	return &nicASGAssociation
}

// func makeVirtualMachine(cfg *cfg.Config, stack tf.TerraformStack, naming *naming.Naming, rg *rg.ResourceGroup, nic *nic.NetworkInterface) *vm.VirtualMachine {
// 	virtualMachine := vm.NewVirtualMachine(stack, cfg.Ids().VirtualMachine(), &vm.VirtualMachineConfig{
// 		Name:              (*naming).VirtualMachineOutput(),
// 		Location:          (*rg).Location(),
// 		ResourceGroupName: (*rg).Name(),
//
// 		NetworkInterfaceIds: []string{(*nic).Id()},
// 		VmSize:              jsii.String("Standard_B1s"),
//	}
//
// 	return &virtualMachine
// }

// FIXME: Properly Select. <eris 2023-05-08>
func MongoDBSubnet(subnets []*vnet.VirtualNetworkSubnetOutputReference) *vnet.VirtualNetworkSubnetOutputReference {
	return subnets[0]
}

func VirtualMachineSubnet(subnets []*vnet.VirtualNetworkSubnetOutputReference) *vnet.VirtualNetworkSubnetOutputReference {
	return subnets[1]
}
