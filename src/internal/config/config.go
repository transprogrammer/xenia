package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const ConfigFile = "config.json"

type regionsJSON struct {
	Primary   *string `json:"primary"`
	Secondary *string `json:"secondary"`
}
type regions struct {
	primary   *string
	secondary *string
}

type subnetsJSON struct {
	VirtualMachine *string `json:"virtualMachine"`
	MongoDB        *string `json:"mongoDB"`
}
type subnets struct {
	virtualMachine *string
	mongoDB        *string
}

type imageReferenceJSON struct {
	Publisher *string `json:"publisher"`
	Offer     *string `json:"offer"`
	Sku       *string `json:"sku"`
	Version   *string `json:"version"`
}
type imageReference struct {
	publisher *string
	offer     *string
	sku       *string
	version   *string
}

type virtualMachineJSON struct {
	Size               *string         `json:"size"`
	StorageAccountType *string         `json:"storageAccountType"`
	ImageReference     *imageReference `json:"imageReference"`
}
type virtualMachine struct {
	size               *string
	storageAccountType *string
	imageReference     *imageReference
}

type databaseAccountJSON struct {
	Kind                    *string   `json:"kind"`
	ServerVersion           *string   `json:"serverVersion"`
	OfferType               *string   `json:"offerType"`
	BackupPolicyType        *string   `json:"backupPolicyType"`
	DefaultConsistencyLevel *string   `json:"defaultConsistencyLevel"`
	Capabilities            []*string `json:"capabilities"`
}
type databaseAccount struct {
	kind                    *string
	serverVersion           *string
	offerType               *string
	backupPolicyType        *string
	defaultConsistencyLevel *string
	capabilities            []*string
}

type configJSON struct {
	SubscriptionId  *string              `json:"subscriptionId"`
	ProjectName     *string              `json:"projectName"`
	Regions         *string              `json:"regions"`
	AddressSpace    []*string            `json:"addressSpace"`
	Subnets         *subnets             `json:"subnets"`
	VirtualMachine  *virtualMachineJSON  `json:"virtualMachine"`
	DatabaseAccount *databaseAccountJSON `json:"databaseAccount"`
	ids             *ids
}

func MakeConfig() *Config {
	configJSON := makeConfigJSON()
	return makeConfigFromJSON
	config
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

func (cfg Config) ProjectName() *string {
	return cfg.projectName
}

func (cfg Config) SubscriptionId() *string {
	return cfg.subscriptionId
}

func (cfg Config) AddressSpace() []*string {
	return cfg.addressSpace
}

func (cfg Config) Subnets() *subnets {
	return cfg.subnets
}

func (cfg Config) VirtualMachine() *virtualMachine {
	return cfg.virtualMachine
}

func (cfg Config) DatabaseAccount() *databaseAccount {
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
