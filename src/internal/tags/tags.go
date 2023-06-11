package tags

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/internal/config"
	"github.com/transprogrammer/xenia/internal/stack"
)

func AddTags(stacks []stack.Stack) {
	for _, stack := range stacks {
		tfStack := stack.Stack()
		tfstack.DependsOn = []cdktf.ITerraformStack{tfStack}
		cdktf.Aspects_Of(stack).Add(NewTagsAspect("Project", config.Config.ProjectName))
	}
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
