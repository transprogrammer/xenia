package stack

import "github.com/hashicorp/terraform-cdk-go/cdktf"

type Stack interface {
	Stack() cdktf.TerraformStack
}
