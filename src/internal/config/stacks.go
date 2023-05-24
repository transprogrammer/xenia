package config

import (
	i "github.com/aws/jsii-runtime-go"
)

type StackNames struct {
	Core    *string
	MongoDB *string
	Jumpbox *string
}

var Stacks StackNames = StackNames{
	Core:    i.String("core"),
	MongoDB: i.String("mongodb"),
	Jumpbox: i.String("jumpbox"),
}
