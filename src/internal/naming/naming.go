package naming

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/generated/naming"
	"github.com/transprogrammer/xenia/internal/config"
)

type Naming naming.Naming

func NewNaming(stack cdktf.TerraformStack, suffixes []*string) naming.Naming {

	id := *config.Ids.Naming
	for _, suffix := range suffixes {
		id = id + "_" + *suffix
	}

	return naming.NewNaming(stack, &id, &naming.NamingConfig{
		Prefix:               &[]*string{config.Config.ProjectName},
		UniqueIncludeNumbers: jsii.Bool(false),
		Suffix:               &suffixes,
	})
}
