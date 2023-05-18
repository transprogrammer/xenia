package main

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/transprogrammer/xenia/internal/core"
	"github.com/transprogrammer/xenia/internal/mongodb"
)

func main() {
	app := cdktf.NewApp(nil)

	coreStack := core.MakeCoreStack(app)
	jumpbox.NewJumpboxStack(app, coreStack)

	mongodb.NewMongoDBStack(coreStack)

	app.Synth()
}
