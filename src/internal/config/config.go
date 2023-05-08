package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const ConfigFile = "config.json"

type regions struct {
	primary   *string `json:"primary"`
	secondary *string `json:"secondary"`
}

type subnets struct {
	virtualMachine *string `json:"virtualMachine"`
	mongoDB        *string `json:"mongoDB"`
}

type imageReference struct {
	publisher *string `json:"publisher"`
	offer     *string `json:"offer"`
	sku       *string `json:"sku"`
	version   *string `json:"version"`
}

type virtualMachine struct {
	size               *string        `json:"size"`
	storageAccountType *string        `json:"storageAccountType"`
	imageReference     imageReference `json:"imageReference"`
}

type databaseAccount struct {
	kind                    *string   `json:"kind"`
	serverVersion           *string   `json:"serverVersion"`
	offerType               *string   `json:"offerType"`
	backupPolicyType        *string   `json:"backupPolicyType"`
	defaultConsistencyLevel *string   `json:"defaultConsistencyLevel"`
	capabilities            []*string `json:"capabilities"`
}

type Config struct {
	subscriptionId *string        `json:"subscriptionId"`
	projectName    *string        `json:"projectName"`
	regions        *string        `json:"regions"`
	addressSpace   []*string      `json:"addressSpace"`
	subnets        subnets        `json:"subnets"`
	virtualMachine virtualMachine `json:"virtualMachine"`
	ids            *ids
}

func MakeConfig() *Config {
	config := makeConfig()
	config.ids = makeIds()

	return config
}

func makeConfig() *Config {
	file, err := os.Open(ConfigFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func (cfg Config) projectName() *string {
	return cfg.projectName
}

func (cfg Config) subscriptionId() *string {
	return cfg.subscriptionId
}

func (cfg Config) addressSpace() []*string {
	return cfg.addressSpace
}

func (cfg Config) subnets() map[*string]*string {
	return cfg.subnets
}

func (cfg Config) virtualMachine() *virtualMachine {
	return cfg.virtualMachine
}

func (cfg Config) databaseAccount() *databaseAccount {
	return cfg.databaseAccount
}

func (cfg Config) Ids() *ids {
	return cfg.ids
}

func (rgn regions) Primary() *string {
	return rgn.primary
}

func (rgn regions) Secondary() *string {
	return rgn.secondary
}

func (sbn subnets) VirtualMachine() *string {
	return sbn.virtualMachine
}

func (sbn subnets) MongoDB() *string {
	return sbn.mongoDB
}

func (img imageReference) Publisher() *string {
	return img.publisher
}

func (img imageReference) Offer() *string {
	return img.offer
}

func (img imageReference) Sku() *string {
	return img.sku
}

func (img imageReference) Version() *string {
	return img.version
}

func (vm virtualMachine) Size() *string {
	return vm.size
}

func (vm virtualMachine) StorageAccountType() *string {
	return vm.storageAccountType
}

func (vm virtualMachine) ImageReference() *imageReference {
	return vm.imageReference
}

func (db databaseAccount) Kind() *string {
	return db.kind
}

func (db databaseAccount) ServerVersion() *string {
	return db.serverVersion
}

func (db databaseAccount) OfferType() *string {
	return db.offerType
}

func (db databaseAccount) BackupPolicyType() *string {
	return db.backupPolicyType
}

func (db databaseAccount) DefaultConsistencyLevel() *string {
	return db.defaultConsistencyLevel
}

func (db databaseAccount) Capabilities() []*string {
	return db.capabilities
}
