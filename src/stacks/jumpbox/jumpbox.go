package jumpbox

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	i "github.com/aws/jsii-runtime-go"
	p "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/provider"
	m "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v5/virtualmachine"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	x "github.com/transprogrammer/xenia/internal/config"
	c "github.com/transprogrammer/xenia/internal/core"
	n "github.com/transprogrammer/xenia/internal/naming"
)

type JumpboxStack struct {
	TerraformStack cdktf.TerraformStack
}

func (s JumpboxStack) Stack() cdktf.TerraformStack {
	return s.TerraformStack
}

func NewStack(scope constructs.Construct, coreStack c.CoreStack) JumpboxStack {
	stackName := fmt.Sprintf("%s-jumpbox", *x.Config.ProjectName)

	stack := t.NewTerraformStack(app, &stackName)

	naming := coreStack.JumpboxNaming
	
	resourceGroup := c.NewResourceGroup(stack, naming)
	publicIP := NewPublicIP(stack, naming, resourceGroup)

	applicationSecurityGroup := NewApplicationSecurityGroup(stack, naming, resourceGroup)
	networkSecurityGroup := NewNetworkSecurityGroup(stack, naming, resourceGroup)

	networkInterface := NewNetworkInterface(stack, naming, resourceGroup, virtualNetwork, publicIP)
	NewNICASGAssocation(stack, networkInterface, applicationSecurityGroup)
	NewNICNSGAssocation(stack, networkInterface, networkSecurityGroup)

	virtualMachine := NewVirtualMachine(stack, naming, resourceGroup)

	providerFunc(stack)

	return stack
}

func NewVirtualMachine(stack t.TerraformStack, naming n.NamingModule, resourceGroup c.ResourceGroup) m.VirtualMachine {
	input := m.VirtualMachineConfig{
		Name: 						naming.VirtualMachineOutput(),
		Location: 					x.Config.Regions.Primary,
		ResourceGroupName: 			resourceGroup.Name(),
	VmSize:            x.Config.VirtualMachine.Size,
	StorageImageReference: &m.VirtualMachineStorageImageReference{
		Publisher: x.Config.VirtualMachine.Image.Publisher,
		Offer:     x.Config.VirtualMachine.Image.Offer,
		Sku:       x.Config.VirtualMachine.Image.Sku,
		Version:   x.Config.VirtualMachine.Image.Version,
	},
	StorageOsDisk: &m.VirtualMachineStorageOsDisk{
		Name:            i.String("osdisk"),
		CreateOption:    i.String("FromImage"),
		ManagedDiskType: x.Config.VirtualMachine.StorageAccountType,
	},
	NetworkInterfaceIds: &[]*string{NIC.Id()},
	OsProfile: &m.VirtualMachineOsProfile{
		ComputerName:  Nme.VirtualMachineOutput(),
		AdminUsername: x.Config.VirtualMachine.AdminUsername,
		AdminPassword: x.Config.VirtualMachine.SSHPublicKey,
	},
	OsProfileLinuxConfig: &m.VirtualMachineOsProfileLinuxConfig{
		DisablePasswordAuthentication: i.Bool(true),
		SshKeys: &m.VirtualMachineOsProfileLinuxConfigSshKeys{
			Path:    i.String("/home/" + *x.Config.VirtualMachine.AdminUsername + "/.ssh/authorized_keys"),
			KeyData: x.Config.VirtualMachine.SSHPublicKey,
		},
	},
})
