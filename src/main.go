package main

import (
	p "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	t "github.com/hashicorp/terraform-cdk-go/cdktf"
	x "github.com/transprogrammer/xenia/internal/config"
	"github.com/transprogrammer/xenia/internal/core"
)

func main() {
	app := t.NewApp(nil)

	providerFunc := func(stack t.TerraformStack) p.AzurermProvider {
		input := p.AzurermProviderConfig{
			Features:       &p.AzurermProviderFeatures{},
			SubscriptionId: x.Config.SubscriptionId,
		}

		return p.NewAzurermProvider(stack, x.Ids.AzureRMProvider, &input)
	}

	core.MakeCoreStack(app, providerFunc)
	//jumpbox.NewJumpboxStack(app, coreStack)

	//mongodb.NewMongoDBStack(coreStack)

	app.Synth()
}
