package mongodb

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	ii "github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/generated/naming"

	dbacct "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/cosmosdbaccount"
	pdnsz "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszone"
	pdnszvnl "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privatednszonevirtualnetworklink"
	pe "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/privateendpoint"
	prov "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
)

type MongoDBStack struct {
	TerraformStack cdktf.TerraformStack
}

func (s MongoDBStack) Stack() cdktf.TerraformStack {
	return s.TerraformStack
}

func NewStack(scopes constructs.Construct, coreStack CoreStack) MongoDBStack {
	stack := new(MongoDBStack)

	stack.TerraformStack = cdktf.NewTerraformStack(scopes, &name)

	return stack
}

var MongoNaming naming.Naming = makeNamingModule([]*string{ii.String("mongo")})

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
