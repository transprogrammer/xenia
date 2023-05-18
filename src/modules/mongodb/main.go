package main

import (
	"fmt"

	ii "github.com/aws/jsii-runtime-go"
	"github.com/transprogrammer/xenia/generated/naming"

	asg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/applicationsecuritygroup"
	dbacct "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/cosmosdbaccount"
	nic "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterface"
	nicasg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterfaceapplicationsecuritygroupassociation"
	nicnsg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterfacesecuritygroupassociation"
	nsg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networksecuritygroup"
	pdnsz "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszone"
	pdnszvnl "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszonevirtualnetworklink"
	pe "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privateendpoint"
	prov "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	ip "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/publicip"
	rg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/resourcegroup"
	vm "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualmachine"
	vnet "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualnetwork"
)

func makeNamingModule(suffixes []*string) naming.Naming {
	return naming.NewNaming(Stk, Ids.NamingModule, &naming.NamingConfig{
		Prefix:               &[]*string{Cfg.ProjectName},
		UniqueIncludeNumbers: ii.Bool(false),
		Suffix:               &suffixes,
	})
}

var Nme naming.Naming = makeNamingModule([]*string{})
var MongoNaming naming.Naming = makeNamingModule([]*string{ii.String("mongo")})

var Rg rg.ResourceGroup = rg.NewResourceGroup(Stk, Ids.ResourceGroup,
	&rg.ResourceGroupConfig{
		Name:     Nme.ResourceGroupOutput(),
		Location: Cfg.Regions.Primary,
	},
)

var Ip ip.PublicIp = ip.NewPublicIp(Stk, Ids.PublicIPAddress,
	&ip.PublicIpConfig{
		Name:                 Nme.PublicIpOutput(),
		Location:             Rg.Location(),
		ResourceGroupName:    Rg.Name(),
		Sku:                  ii.String("Basic"),
		AllocationMethod:     ii.String("Dynamic"),
		IpVersion:            ii.String("IPv4"),
		DomainNameLabel:      Cfg.ProjectName,
		IdleTimeoutInMinutes: ii.Number(4),
	},
)

// ???: Inline subnet too enable updating in-place. <>
// SEE: https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/deployment-modes#incremental-mode <>
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

var VNet vnet.VirtualNetwork = makeVirtualNetwork()
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

var VM vm.VirtualMachine = vm.NewVirtualMachine(Stk, Ids.VirtualMachine, &vm.VirtualMachineConfig{
	Name:              Nme.VirtualMachineOutput(),
	Location:          Cfg.Regions.Primary,
	ResourceGroupName: Rg.Name(),
	VmSize:            Cfg.VirtualMachine.Size,
	StorageImageReference: &vm.VirtualMachineStorageImageReference{
		Publisher: Cfg.VirtualMachine.Image.Publisher,
		Offer:     Cfg.VirtualMachine.Image.Offer,
		Sku:       Cfg.VirtualMachine.Image.Sku,
		Version:   Cfg.VirtualMachine.Image.Version,
	},
	StorageOsDisk: &vm.VirtualMachineStorageOsDisk{
		Name:            ii.String("osdisk"),
		CreateOption:    ii.String("FromImage"),
		ManagedDiskType: Cfg.VirtualMachine.StorageAccountType,
	},
	NetworkInterfaceIds: &[]*string{NIC.Id()},
	OsProfile: &vm.VirtualMachineOsProfile{
		ComputerName:  Nme.VirtualMachineOutput(),
		AdminUsername: Cfg.VirtualMachine.AdminUsername,
		AdminPassword: Cfg.VirtualMachine.SSHPublicKey,
	},
	OsProfileLinuxConfig: &vm.VirtualMachineOsProfileLinuxConfig{
		DisablePasswordAuthentication: ii.Bool(true),
		SshKeys: &vm.VirtualMachineOsProfileLinuxConfigSshKeys{
			Path:    ii.String("/home/" + *Cfg.VirtualMachine.AdminUsername + "/.ssh/authorized_keys"),
			KeyData: Cfg.VirtualMachine.SSHPublicKey,
		},
	},
})

var dbAccount dbacct.CosmosdbAccount = dbacct.NewCosmosdbAccount(Stk, Ids.CosmosDBAccount, &dbacct.CosmosdbAccountConfig{
	Name:                       Nme.CosmosdbAccountOutput(),
	Location:                   Cfg.Regions.Primary,
	ResourceGroupName:          Rg.Name(),
	Kind:                       ii.String("MongoDB"),
	OfferType:                  ii.String("Standard"),
	MongoServerVersion:         ii.String("4.2"),
	PublicNetworkAccessEnabled: ii.Bool(false),
	ConsistencyPolicy: &dbacct.CosmosdbAccountConsistencyPolicy{
		ConsistencyLevel: ii.String("Eventual"),
	},
	GeoLocation: &[]*dbacct.CosmosdbAccountGeoLocation{
		{
			Location:         Cfg.Regions.Secondary,
			FailoverPriority: ii.Number(0),
			ZoneRedundant:    ii.Bool(false),
		},
	},
	Capabilities: &[]*dbacct.CosmosdbAccountCapabilities{
		{
			Name: ii.String("DisabledRateLimitingResponses"),
		},
		{
			Name: ii.String("EnableServerless"),
		},
	},
})

// ???: Add database? <>
// TODO: Add sg assocs <>
// TODO: Hardening <>

var PrivateEndpoint pe.PrivateEndpoint = pe.NewPrivateEndpoint(Stk, Ids.PrivateEndpoint, &pe.PrivateEndpointConfig{
	Name:              Nme.PrivateEndpointOutput(),
	Location:          Cfg.Regions.Primary,
	ResourceGroupName: Rg.Name(),
	SubnetId:          MongoDBSubnet.Id(),
	PrivateServiceConnection: &pe.PrivateEndpointPrivateServiceConnection{
		Name:                        ii.String("cosmosdb"),
		PrivateConnectionResourceId: dbAccount.Id(),
		SubresourceNames:            &[]*string{ii.String("MongoDB")},
	},
})

var PrivateDNSZone pdnsz.PrivateDnsZone = pdnsz.NewPrivateDnsZone(Stk, Ids.PrivateDNSZone, &pdnsz.PrivateDnsZoneConfig{
	Name:              ii.String("privatelink.mongo.cosmos.azure.com"),
	ResourceGroupName: Rg.Name(),
})

var PrivateDNSZoneVirtualNetworkLink pdnszvnl.PrivateDnsZoneVirtualNetworkLink = pdnszvnl.NewPrivateDnsZoneVirtualNetworkLink(Stk, Ids.PrivateDNSZoneVirtualNetworkLink, &pdnszvnl.PrivateDnsZoneVirtualNetworkLinkConfig{
	Name:                ii.String(fmt.Sprintf("%s-vnetlink", *MongoNaming.PrivateDnsZoneOutput())),
	ResourceGroupName:   Rg.Name(),
	PrivateDnsZoneName:  PrivateDNSZone.Name(),
	VirtualNetworkId:    VNet.Id(),
	RegistrationEnabled: ii.Bool(true),
})

func main() {
	prov.NewAzurermProvider(Stk, Ids.AzureRMProvider, &prov.AzurermProviderConfig{
		Features:       &prov.AzurermProviderFeatures{},
		SubscriptionId: Cfg.SubscriptionId,
	})

	App.Synth()
}
