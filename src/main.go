package main

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/internal/stack"
	"github.com/transprogrammer/xenia/internal/tags"
	"github.com/transprogrammer/xenia/stacks/aks"
	"github.com/transprogrammer/xenia/stacks/core"
	"github.com/transprogrammer/xenia/stacks/jumpbox"
	"github.com/transprogrammer/xenia/stacks/mongodb"
)

func main() {
	app := cdktf.NewApp(nil)

	stacks := make([]stack.Stack, 0, 4)

	coreStack := core.NewStack(app)
	stacks = append(stacks, coreStack)

	jumpboxStack := jumpbox.NewStack(app, coreStack)
	stacks = append(stacks, jumpboxStack)

	mongodbStack := mongodb.NewStack(app, coreStack)
	stacks = append(stacks, mongodbStack)

	aksStack := aks.NewStack(app, coreStack)
	stacks = append(stacks, aksStack)

	tags.AddTags(stacks)

	app.Synth()
}
