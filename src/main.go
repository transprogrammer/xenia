package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"

	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/linuxvirtualmachine"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/networkinterface"
	provider "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	"github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/subnet"
	vnet "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualnetwork"

	"github.com/transprogrammer/xenia/generated/naming"
)

func main() {
	app := makeApp()
	stack := makeStack(config, app)
	config := makeConfig(config, stack)

	makeProvider(config, stack)
	addResources(config, stack)

	app.Synth()
}

func makeApp() tf.App {
	return tf.NewApp(nil)
}

func makeStack(config Config, scope constructs.Construct) cdktf.TerraformStack {
	return cdktf.NewTerraformStack(scope, config.ProjectName)
}

func makeProvider(config Config, stack tf.TerraformStack) provider.AzurermProvider {
	return provider.NewAzurermProvider(stack, config.ProviderName, provider.AzurermProviderConfig{
		Features:       &prv.AzurermProviderFeatures{},
		SubscriptionId: SubscriptionId,
	})

}

func addResources(config Config, stack tf.TerraformStack) {
	naming := addNaming(config, stack)
	rg := addRG(config, stack, naming)
	vnet := addVNet(config, stack, naming, rg)

	// addNullResource(stack, naming)
}

func addNaming(config Config, stack tf.TerraformStack) naming.Naming {
	return naming.NewNaming(stack, jsii.String("naming"), &naming.NamingConfig{
		Prefix:               &[]*string{ProjectName},
		UniqueIncludeNumbers: jsii.Bool(false),
	})
}

func addRG(config Config, stack tf.TerraformStack, naming naming.Naming) rg.ResourceGroup {
	return rg.NewResourceGroup(stack, jsii.String("rg"), rg.ResourceGroupConfig{
		Name:     naming.ResourceGroupOutput(),
		Location: config.PrimaryRegion,
	})
}

func addVNet(config Config, stack tf.TerraformStack, naming naming.Naming, rg rg.ResourceGroup) vnet.VirtualNetwork {
	return vnet.NewVirtualNetwork(stack, jsii.String("vnet"), vnet.VirtualNetworkConfig{

		Name:              naming.VirtualNetworkOutput(),
		AddressSpace:      config.AddressSpace,
		Location:          rg.Location(),
		ResourceGroupName: rg.Name(),
	})
}

func addSubnet(config Config, stack tf.TerraformStack, naming naming.Naming, rg rg.ResourceGroup, vnet vnet.VirtualNetwork) subnet.Subnet {
	return subnet.NewSubnet(stack, jsii.String("subnet"), subnet.SubnetConfig{
		Name:               naming.SubnetOutput(),
		ResourceGroupName:  rg.Name(),
		VirtualNetworkName: vnet.Name(),
		AddressPrefixes:    config.AddressPrefixes,
	})

	// func addNullResource(stack tf.TerraformStack) *nullresource.Resource {
	// 	return nullresource.NewResource(stack, "null", &nullresource.ResourceConfig{})
	// }

	//Create the test Virtual Machine with its Network Interface
	vm_nic := networkinterface.NewNetworkInterface(stack, jsii.String("test_vm_nic"), &networkinterface.NetworkInterfaceConfig{
		Name:              jsii.String("test-vm-nic"),
		Location:          rg.Location(),
		ResourceGroupName: rg.Name(),

		IpConfiguration: &[]*networkinterface.NetworkInterfaceIpConfiguration{{
			Name:                       jsii.String("internal"),
			SubnetId:                   vm_nw_sn.Id(),
			PrivateIpAddressAllocation: jsii.String("Dynamic"),
		}},
	})

	vm := linuxvirtualmachine.NewLinuxVirtualMachine(stack, jsii.String("test_vm"), &linuxvirtualmachine.LinuxVirtualMachineConfig{
		Name:                jsii.String("test-vm"),
		Location:            rg.Location(),
		ResourceGroupName:   rg.Name(),
		Size:                jsii.String("Standard_F2"),
		AdminUsername:       jsii.String("adminuser"),
		NetworkInterfaceIds: &[]*string{vm_nic.Id()},

		AdminSshKey: &[]*linuxvirtualmachine.LinuxVirtualMachineAdminSshKey{{
			Username:  jsii.String("glados"),
			PublicKey: tf.Fn_File(jsii.String("~/.ssh/id_rsa.pub")),
		}},

		OsDisk: &linuxvirtualmachine.LinuxVirtualMachineOsDisk{
			Caching:            jsii.String("ReadWrite"),
			StorageAccountType: jsii.String("Standard_LRS"),
		},

		SourceImageReference: &linuxvirtualmachine.LinuxVirtualMachineSourceImageReference{
			Publisher: jsii.String("Canonical"),
			Offer:     jsii.String("UbuntuServer"),
			Sku:       jsii.String("16.04-LTS"),
			Version:   jsii.String("latest"),
		},
	})

	//Output stuff
	tf.NewTerraformOutput(stack, jsii.String("names"), &tf.TerraformOutputConfig{
		Value: &[]*string{vm.Name(), rg.Name()},
	})

	return stack
}
