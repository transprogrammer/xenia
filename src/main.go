package main

import (
	"github.com/aws/jsii-runtime-go"
	tf "github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/generated/naming"
	"github.com/transprogrammer/xenia/internal/cfg"

	prov "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	rg "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/resourcegroup"
	vnet "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualnetwork"
)

func main() {
	makeApp().Synth()
}

func makeApp() tf.App {
	cfg := cfg.MakeConfig()
	app := tf.NewApp(nil)

	stack := tf.NewTerraformStack(app, cfg.ProjectName())
	makeAzureRMProvider(cfg, stack)

	naming := makeNamingModule(cfg, stack, []*string{})
	rg := makeResourceGroup(cfg, stack, naming)

	makeVirtualNetwork(cfg, stack, naming, rg)

	return app
}

func makeAzureRMProvider(cfg *cfg.Config, stack tf.TerraformStack) prov.AzurermProvider {
	return prov.NewAzurermProvider(stack, cfg.Ids().AzureRMProvider(), &prov.AzurermProviderConfig{
		Features:       &prov.AzurermProviderFeatures{},
		SubscriptionId: cfg.SubscriptionId(),
	})
}

func makeNamingModule(cfg *cfg.Config, stack tf.TerraformStack, suffixes []*string) naming.Naming {
	prefix := []*string{cfg.ProjectName()}

	return naming.NewNaming(stack, cfg.Ids().NamingModule(), &naming.NamingConfig{
		Prefix:               &prefix,
		UniqueIncludeNumbers: jsii.Bool(false),
		Suffix:               &suffixes,
	})
}

func makeResourceGroup(cfg *cfg.Config, stack tf.TerraformStack, naming naming.Naming) rg.ResourceGroup {
	return rg.NewResourceGroup(stack, cfg.Ids().ResourceGroup(), &rg.ResourceGroupConfig{
		Name:     naming.ResourceGroupOutput(),
		Location: cfg.PrimaryRegion(),
	})
}

// ???: Inline subnet too enable updating in-place. <>
// SEE: https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/deployment-modes#incremental-mode <>
func makeVirtualNetwork(cfg *cfg.Config, stack tf.TerraformStack, naming naming.Naming, rg rg.ResourceGroup) vnet.VirtualNetwork {
	subnets := make([]vnet.VirtualNetworkSubnet, len(cfg.Subnets()))
	i := 0
	for k, v := range cfg.Subnets() {
		subnetNaming := makeNamingModule(cfg, stack, []*string{k})

		subnets[i] = vnet.VirtualNetworkSubnet{
			Name:          subnetNaming.SubnetOutput(),
			AddressPrefix: v,
			//SecurityGroup: nil,
		}
		i++
	}

	addressSpace := cfg.AddressSpace()

	return vnet.NewVirtualNetwork(stack, cfg.Ids().VirtualNetwork(), &vnet.VirtualNetworkConfig{
		Name:              naming.VirtualNetworkOutput(),
		AddressSpace:      &addressSpace,
		Location:          rg.Location(),
		ResourceGroupName: rg.Name(),
		Subnet:            subnets,
	})
}
