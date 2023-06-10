package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/internal/config"
)

func main() {
	app := cdktf.NewApp(nil)

	stacks := make([]cdktf.TerraformStack, 0, 4)

	coreStack := core.NewStack(app)
	stacks = append(stacks, coreStack)

	jumpboxStack := jumpbox.NewStack(app, coreStack)
	stacks = append(stacks, jumpboxStack)

	mongodbStack := mongodb.NewStack(app, coreStack)
	stacks = append(stacks, mongodbStack)

	aksStack := aks.NewStack(app, coreStack)
	stacks = append(stacks, aksStack)

	for _, stack := range stacks {
		stack.DependsOn = []cdktf.ITerraformStack{coreStack}
		cdktf.Aspects_Of(stack).Add(NewTagsAspect("Project", config.Config.ProjectName))
	}

	app.Synth()
}

type Taggable interface {
	TagsInput() *map[string]*string
	SetTags(val *map[string]*string)
}

type TagsAspect struct {
	Tags *map[string]*string
}

func (taa TagsAspect) Visit(node constructs.IConstruct) {
	if taggable, ok := node.(Taggable); ok {
		existing := *taggable.TagsInput()
		tags := *taa.Tags
		maps.Copy(existing, tags) // requires Go 1.18
		taggable.SetTags(&existing)
	}
}

func NewTagsAspect(tags *map[string]*string) *TagsAspect {
	return &TagsAspect{Tags: tags}
}
