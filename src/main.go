package main

import (
	ii "github.com/aws/jsii-runtime-go"
	tf "github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/generated/naming"

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

var App tf.App = tf.NewApp(nil)
var Stk tf.TerraformStack = tf.NewTerraformStack(App, Cfg.ProjectName)

var Nme naming.Naming = makeNamingModule([]*string{})

var Rg rg.ResourceGroup = rg.NewResourceGroup(Stk, Ids.ResourceGroup, &rg.ResourceGroupConfig{
	Name:     Nme.ResourceGroupOutput(),
	Location: Cfg.Regions.Primary,
},
)

var Ip ip.PublicIp = ip.NewPublicIp(Stk, Ids.PublicIPAddress, &ip.PublicIpConfig{
	Name:                 Nme.PublicIpOutput(),
	Location:             Rg.Location(),
	ResourceGroupName:    Rg.Name(),
	Sku:                  ii.String("Basic"),
	AllocationMethod:     ii.String("Dynamic"),
	IpVersion:            ii.String("IPv4"),
	DomainNameLabel:      Cfg.ProjectName,
	IdleTimeoutInMinutes: ii.Number(4),
})

// ???: Inline subnet too enable updating in-place. <>
// SEE: https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/deployment-modes#incremental-mode <>
var VNet vnet.VirtualNetwork = makeVirtualNetwork()

var mongoDBIndex = float64(0)
var vmIndex = float64(1)

func makeVirtualNetwork() vnet.VirtualNetwork {
	subnets := []vnet.VirtualNetworkSubnet{}
	subnets[int(mongoDBIndex)] = vnet.VirtualNetworkSubnet{
		Name:          (makeNamingModule([]*string{Cfg.Subnets.MongoDB.Postfix})).SubnetOutput(),
		AddressPrefix: Cfg.Subnets.MongoDB.AddressPrefix,
	}

	subnets[int(vmIndex)] = vnet.VirtualNetworkSubnet{
		Name:          (makeNamingModule([]*string{Cfg.Subnets.MongoDB.Postfix})).SubnetOutput(),
		AddressPrefix: Cfg.Subnets.VirtualMachine.AddressPrefix,
	}

	return vnet.NewVirtualNetwork(Stk, Ids.VirtualNetwork, &vnet.VirtualNetworkConfig{
		Name:              Nme.VirtualNetworkOutput(),
		AddressSpace:      Cfg.AddressSpace,
		Location:          Rg.Location(),
		ResourceGroupName: Rg.Name(),
		Subnet:            subnets,
	})
}

var MongoDBSubnet vnet.VirtualNetworkSubnetOutputReference = VNet.Subnet().Get(&mongoDBIndex)
var VMSubnet vnet.VirtualNetworkSubnetOutputReference = VNet.Subnet().Get(&vmIndex)

var NSG nsg.NetworkSecurityGroup = nsg.NewNetworkSecurityGroup(Stk, Ids.NetworkSecurityGroup, &nsg.NetworkSecurityGroupConfig{
	Name:              Nme.NetworkSecurityGroupOutput(),
	Location:          Rg.Location(),
	ResourceGroupName: Rg.Name(),
	SecurityRule: nsg.NetworkSecurityGroupSecurityRule{
		Name:                     ii.String("SSH"),
		Description:              ii.String("Allow SSH"),
		Priority:                 ii.Number(100),
		Direction:                ii.String("Inbound"),
		Access:                   ii.String("Allow"),
		Protocol:                 ii.String("Tcp"),
		SourcePortRange:          ii.String("*"),
		DestinationPortRange:     ii.String("22"),
		SourceAddressPrefix:      ii.String("*"),
		DestinationAddressPrefix: ii.String("*"),
	},
})

var ASG asg.ApplicationSecurityGroup = asg.NewApplicationSecurityGroup(Stk, Ids.ApplicationSecurityGroup, &asg.ApplicationSecurityGroupConfig{
	Name:              Nme.ApplicationSecurityGroupOutput(),
	Location:          Rg.Location(),
	ResourceGroupName: Rg.Name(),
})

var NIC nic.NetworkInterface = nic.NewNetworkInterface(Stk, Ids.NetworkInterface, &nic.NetworkInterfaceConfig{
	Name:              Nme.NetworkInterfaceOutput(),
	Location:          Rg.Location(),
	ResourceGroupName: Rg.Name(),

	IpConfiguration: nic.NetworkInterfaceIpConfiguration{
		Name:              ii.String("ipconfig"),
		Primary:           ii.Bool(true),
		SubnetId:          VMSubnet.Id(),
		PublicIpAddressId: Ip.Id(),
	},
})

var nicNSGAssociation nicnsg.NetworkInterfaceSecurityGroupAssociation = nicnsg.NewNetworkInterfaceSecurityGroupAssociation(Stk, Ids.NetworkInterfaceNSGAssociation, &nicnsg.NetworkInterfaceSecurityGroupAssociationConfig{
	NetworkInterfaceId:     NIC.Id(),
	NetworkSecurityGroupId: NSG.Id(),
})

var nicASGAssociation nicasg.NetworkInterfaceApplicationSecurityGroupAssociation = nicasg.NewNetworkInterfaceApplicationSecurityGroupAssociation(Stk, Ids.NetworkInterfaceASGAssociation, &nicasg.NetworkInterfaceApplicationSecurityGroupAssociationConfig{
	NetworkInterfaceId:         NIC.Id(),
	ApplicationSecurityGroupId: ASG.Id(),
})

func makeNamingModule(suffixes []*string) naming.Naming {
	return naming.NewNaming(Stk, Ids.NamingModule, &naming.NamingConfig{
		Prefix:               &[]*string{Cfg.ProjectName},
		UniqueIncludeNumbers: ii.Bool(false),
		Suffix:               &suffixes,
	})
}

func main() {
	prov.NewAzurermProvider(Stk, Ids.AzureRMProvider, &prov.AzurermProviderConfig{
		Features:       &prov.AzurermProviderFeatures{},
		SubscriptionId: Cfg.SubscriptionId,
	})

	App.Synth()

	//makeVirtualMachine(Cfg, stk, nme, Rg, nic, vnet)
}

// func makeVirtualMachine(Cfg *Cfg.Config, Stk tf.TerraformStack, naming *naming.Naming, Rg Rg.ResourceGroup, nic *nic.NetworkInterface) *vm.VirtualMachine {
// 	virtualMachine := vm.NewVirtualMachine(Stk, Cfg.Ids().VirtualMachine(), &vm.VirtualMachineConfig{
// 		Name:              Nme.VirtualMachineOutput(),
// 		Location:          Rg.Location(),
// 		ResourceGroupName: Rg.Name(),
//
// 		NetworkInterfaceIds: []string{NIC.Id()},
// 		VmSize:              ii.String("Standard_B1s"),
//	}
//
// 	return &virtualMachine
// }
