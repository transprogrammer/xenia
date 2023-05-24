package main

import (
	t "github.com/hashicorp/terraform-cdk-go/cdktf"
	c "github.com/transprogrammer/xenia/internal/core"
)

func main() {
	app := t.NewApp(nil)

	c.MakeCoreStack(app)

	//j.NewJumpboxStack(app, coreStack)
	//m.NewMongoDBStack(app, coreStack)

	app.Synth()
}
