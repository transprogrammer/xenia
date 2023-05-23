package jumpbox

import (
	"fmt"

	"github.com/aws/jsii-runtime-go"
	prv "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	tf "github.com/hashicorp/terraform-cdk-go/cdktf"
	x "github.com/transprogrammer/xenia/internal/config"
)

func NewJumpboxStack(app tf.App, providerFunc func(tf.TerraformStack) prv.AzurermProvider) tf.TerraformStack {
	stackName := fmt.Sprintf("%s-jumpbox", *x.Config.ProjectName)

	stack := tf.NewTerraformStack(app, &stackName)

	providerFunc(stack)

	return stack
}

var VM vm.VirtualMachine = vm.NewVirtualMachine(stack, Ids.VirtualMachine, &vm.VirtualMachineConfig{
	Name:              Nme.VirtualMachineOutput(),
	Location:          x.Config.Regions.Primary,
	ResourceGroupName: Rg.Name(),
	VmSize:            x.Config.VirtualMachine.Size,
	StorageImageReference: &vm.VirtualMachineStorageImageReference{
		Publisher: x.Config.VirtualMachine.Image.Publisher,
		Offer:     x.Config.VirtualMachine.Image.Offer,
		Sku:       x.Config.VirtualMachine.Image.Sku,
		Version:   x.Config.VirtualMachine.Image.Version,
	},
	StorageOsDisk: &vm.VirtualMachineStorageOsDisk{
		Name:            jsii.String("osdisk"),
		CreateOption:    jsii.String("FromImage"),
		ManagedDiskType: x.Config.VirtualMachine.StorageAccountType,
	},
	NetworkInterfaceIds: &[]*string{NIC.Id()},
	OsProfile: &vm.VirtualMachineOsProfile{
		ComputerName:  Nme.VirtualMachineOutput(),
		AdminUsername: x.Config.VirtualMachine.AdminUsername,
		AdminPassword: x.Config.VirtualMachine.SSHPublicKey,
	},
	OsProfileLinuxConfig: &vm.VirtualMachineOsProfileLinuxConfig{
		DisablePasswordAuthentication: jsii.Bool(true),
		SshKeys: &vm.VirtualMachineOsProfileLinuxConfigSshKeys{
			Path:    jsii.String("/home/" + *x.Config.VirtualMachine.AdminUsername + "/.ssh/authorized_keys"),
			KeyData: x.Config.VirtualMachine.SSHPublicKey,
		},
	},
})
