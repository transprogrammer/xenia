package main

import (
	t "github.com/hashicorp/terraform-cdk-go/cdktf"
	c "github.com/transprogrammer/xenia/internal/core"
	j "github.com/transprogrammer/xenia/internal/jumpbox"
	m "github.com/transprogrammer/xenia/internal/mongodb"
)

func main() {
	app := t.NewApp(nil)

	coreStack := c.MakeCoreStack(app)

	j.NewJumpboxStack(app, coreStack)
	m.NewMongoDBStack(app, coreStack)

	app.Synth()
}
