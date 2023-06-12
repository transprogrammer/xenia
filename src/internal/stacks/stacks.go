package stack

import (
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/resourcegroup"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/internal/config"
	"github.com/transprogrammer/xenia/internal/naming"
)

type Stack interface {
	Stack() cdktf.TerraformStack
}

func NewAzureRMProvider(stack cdktf.TerraformStack) provider.AzurermProvider {
	input := provider.AzurermProviderConfig{
		Features:       &provider.AzurermProviderFeatures{},
		SubscriptionId: config.Config.SubscriptionId,
	}

	return provider.NewAzurermProvider(stack, config.Ids.AzureRMProvider, &input)
}

func NewResourceGroup(stack cdktf.TerraformStack, naming naming.NamingModule) resourcegroup.ResourceGroup {
	input := resourcegroup.ResourceGroupConfig{
		Name:     naming.ResourceGroupOutput(),
		Location: config.Config.Regions.Primary,
	}

	return resourcegroup.NewResourceGroup(stack, config.Ids.ResourceGroup, &input)
}
