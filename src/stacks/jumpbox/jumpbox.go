package jumpbox

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
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

	jumpboxIP := NewPublicIP(coreStack, jumpboxNamingModule, resourceGroup)
	NewNetworkInterface(coreStack, jumpboxNaming, resourceGroup, jumpboxSubnet, jumpboxApplicationSecurityGroup, jumpboxNetworkSecurityGroup, jumpboxIP)
	virtualMachine := NewVirtualMachine(stack, naming, resourceGroup)

	providerFunc(stack)

	return stack
}

func NewVirtualMachine(stack t.TerraformStack, naming n.NamingModule, resourceGroup c.ResourceGroup) m.VirtualMachine {
	input := m.VirtualMachineConfig{
		Name:              naming.VirtualMachineOutput(),
		Location:          x.Config.Regions.Primary,
		ResourceGroupName: resourceGroup.Name(),
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
	}

	return m.NewVirtualMachine(stack, &naming.VirtualMachineOutput(), &input)
}
o

func NewPublicIP(
	stack cdktf.TerraformStack,
	naming Naming,
	group rg.ResourceGroup,
) publicip.PublicIp {

	input := publicip.PublicIpConfig{
		Name:                 PublicIpOutput(),
		Location:             Config.Regions.Primary,
		ResourceGroupName:    group.Name(),
		Sku:                  jsii.String("Basic"),
		AllocationMethod:     jsii.String("Dynamic"),
		IpVersion:            jsii.String("IPv4"),
		DomainNameLabel:      Config.ProjectName,
		IdleTimeoutInMinutes: jsii.Number(4),
	}

	return publicip.NewPublicIp(stack, Ids.PublicIPAddress, &input)
}

