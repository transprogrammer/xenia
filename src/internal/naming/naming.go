package naming

import (
	ii "github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/generated/naming"
	cfg "github.com/transprogrammer/xenia/internal/config"
)

type NamingModule naming.Naming

func NewNamingModule(stack cdktf.TerraformStack, suffixes []*string) naming.Naming {

	id := *cfg.Ids.NamingModule
	for _, suffix := range suffixes {
		id = id + "_" + *suffix
	}

	return naming.NewNaming(stack, &id, &naming.NamingConfig{
		Prefix:               &[]*string{cfg.Config.ProjectName},
		UniqueIncludeNumbers: ii.Bool(false),
		Suffix:               &suffixes,
	})
}
