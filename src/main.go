package main

import (
	t "github.com/hashicorp/terraform-cdk-go/cdktf"
	c "github.com/transprogrammer/xenia/internal/core"
	j "github.com/transprogrammer/xenia/internal/jumpbox"
)

func main() {
	app := t.NewApp(nil)

	coreStack := c.MakeCoreStack(app)

	j.NewJumpboxStack(app, coreStack)
	//m.NewMongoDBStack(app, coreStack)
	//a.NewAKSStack(app, coreStack)

	app.Synth()
}
