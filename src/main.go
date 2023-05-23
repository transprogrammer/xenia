package main

import (
	prv "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	tf "github.com/hashicorp/terraform-cdk-go/cdktf"
	x "github.com/transprogrammer/xenia/internal/config"
	"github.com/transprogrammer/xenia/internal/core"
	"github.com/transprogrammer/xenia/internal/jumpbox"
	"github.com/transprogrammer/xenia/internal/mongodb"
)

func main() {
	app := tf.NewApp(nil)

	providerFunc := func(stack tf.TerraformStack) prv.AzurermProvider {
		input := prv.AzurermProviderConfig{
			Features:       &prv.AzurermProviderFeatures{},
			SubscriptionId: x.Config.SubscriptionId,
		}

		return prv.NewAzurermProvider(stack, x.Ids.AzureRMProvider, &input)
	}

	coreStack := core.MakeCoreStack(app, providerFunc)
	jumpbox.NewJumpboxStack(app, coreStack)

	mongodb.NewMongoDBStack(coreStack)

	app.Synth()
}
