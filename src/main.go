package main

import (
	"cdk.tf/go/stack/generated/naming"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	tf "github.com/hashicorp/terraform-cdk-go/cdktf"

	prov "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	rg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/resourcegroup"
	vnet "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualnetwork"
)

var config = makeConfig()
var stackId = config.ProjectName()

var providerId = jsii.String("azurerm")
var namingId = jsii.String("azurerm_naming")

var rgId = jsii.String("rg")
var vnetId = jsii.String("vnet")
var subnetId = jsii.String("subnet")

var app = makeApp()

func main() {
	app.Synth()
}

func makeApp() tf.App {
	app := tf.NewApp(nil)

	makeStack(app)

	return app
}

func makeStack(scope constructs.Construct) {
	stack := tf.NewTerraformStack(scope, stackId)

	makeProvider(config, stack)

	naming := makeNaming(config, stack, []*string{})

	makeResources(config, stack, naming)
}

func makeProvider(config Config, stack tf.TerraformStack) prov.AzurermProvider {
	return prov.NewAzurermProvider(stack, providerId, &prov.AzurermProviderConfig{
		Features:       &prov.AzurermProviderFeatures{},
		SubscriptionId: config.SubscriptionId(),
	})
}

func makeNaming(config Config, stack tf.TerraformStack, suffixes []*string) naming.Naming {
	prefix := []*string{config.ProjectName()}

	return naming.NewNaming(stack, namingId, &naming.NamingConfig{
		Prefix:               &prefix,
		UniqueIncludeNumbers: jsii.Bool(false),
		Suffix:               &suffixes,
	})
}

func makeResources(config Config, stack tf.TerraformStack, naming naming.Naming) {
	rg := makeRG(config, stack, naming)

	makeVNet(config, stack, naming, rg)
}

func makeRG(config Config, stack tf.TerraformStack, naming naming.Naming) rg.ResourceGroup {
	return rg.NewResourceGroup(stack, rgId, &rg.ResourceGroupConfig{
		Name:     naming.ResourceGroupOutput(),
		Location: config.PrimaryRegion(),
	})
}

// ???: Inline subnet too enable updating in-place. <>
// SEE: https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/deployment-modes#incremental-mode <>
func makeVNet(config Config, stack tf.TerraformStack, naming naming.Naming, rg rg.ResourceGroup) vnet.VirtualNetwork {
	subnets := make([]vnet.VirtualNetworkSubnet, len(config.Subnets()))
	i := 0
	for k, v := range config.Subnets() {
		subnetNaming := makeNaming(config, stack, []*string{k})
		subnets[i] = vnet.VirtualNetworkSubnet{
			Name:          subnetNaming.SubnetOutput(),
			AddressPrefix: v,
			//SecurityGroup: nil,
		}
		i++
	}

	addressSpace := config.AddressSpace()

	return vnet.NewVirtualNetwork(stack, vnetId, &vnet.VirtualNetworkConfig{
		Name:              naming.VirtualNetworkOutput(),
		AddressSpace:      &addressSpace,
		Location:          rg.Location(),
		ResourceGroupName: rg.Name(),
		Subnet:            subnets,
	})
}
