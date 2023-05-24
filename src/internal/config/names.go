package config

type Namings struct {
	Core    []*string
	MongoDB []*string
	Jumpbox []*string
}

var Names Namings = Namings{
	Core:    []*string{Config.ProjectName, Stacks.Core},
	MongoDB: []*string{Config.ProjectName, Stacks.MongoDB},
	Jumpbox: []*string{Config.ProjectName, Stacks.Jumpbox},
}
