package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const ConfigFile = "config.json"

type regionsJSO struct {
	Primary   *string `json:"primary"`
	Secondary *string `json:"secondary"`
}
type Regions struct {
	primary   *string
	secondary *string
}

type subnetJSO struct {
	Postfix       *string `json:"postfix"`
	AddressPrefix *string `json:"addressPrefix"`
}
type Subnet struct {
	postfix       *string
	addressPrefix *string
}

type subnetsJSO struct {
	VirtualMachine *subnetJSO `json:"virtualMachine"`
	MongoDB        *subnetJSO `json:"mongoDB"`
}
type Subnets struct {
	virtualMachine *Subnet
	mongoDB        *Subnet
}

type imageReferenceJSO struct {
	Publisher *string `json:"publisher"`
	Offer     *string `json:"offer"`
	Sku       *string `json:"sku"`
	Version   *string `json:"version"`
}
type ImageReference struct {
	publisher *string
	offer     *string
	sku       *string
	version   *string
}

type virtualMachineJSO struct {
	Size               *string            `json:"size"`
	StorageAccountType *string            `json:"storageAccountType"`
	ImageReferenceJSO  *imageReferenceJSO `json:"imageReference"`
}
type VirtualMachine struct {
	size               *string
	storageAccountType *string
	imageReference     *ImageReference
}

type databaseAccountJSO struct {
	Kind                    *string   `json:"kind"`
	ServerVersion           *string   `json:"serverVersion"`
	OfferType               *string   `json:"offerType"`
	BackupPolicyType        *string   `json:"backupPolicyType"`
	DefaultConsistencyLevel *string   `json:"defaultConsistencyLevel"`
	Capabilities            []*string `json:"capabilities"`
}
type DatabaseAccount struct {
	kind                    *string
	serverVersion           *string
	offerType               *string
	backupPolicyType        *string
	defaultConsistencyLevel *string
	capabilities            []*string
}

type configJSO struct {
	SubscriptionId     *string             `json:"subscriptionId"`
	ProjectName        *string             `json:"projectName"`
	RegionsJSO         *regionsJSO         `json:"regions"`
	AddressSpace       []*string           `json:"addressSpace"`
	SubnetsJSO         *subnetsJSO         `json:"subnets"`
	VirtualMachineJSO  *virtualMachineJSO  `json:"virtualMachine"`
	DatabaseAccountJSO *databaseAccountJSO `json:"databaseAccount"`
}
type Config struct {
	subscriptionId  *string
	projectName     *string
	regions         *Regions
	addressSpace    []*string
	subnets         *Subnets
	virtualMachine  *VirtualMachine
	databaseAccount *DatabaseAccount
	ids             *Ids
}

func MakeConfig() *Config {
	configJSO := makeConfigJSO()
	config := makeConfigFromJSO(configJSO)
	config.ids = makeIds()

	return config
}

func makeConfigJSO() *configJSO {
	file, err := os.Open(ConfigFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var cfg configJSO
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func makeConfigFromJSO(jso *configJSO) *Config {
	return &Config{
		subscriptionId: jso.SubscriptionId,
		projectName:    jso.ProjectName,
		regions: &Regions{
			primary:   jso.RegionsJSO.Primary,
			secondary: jso.RegionsJSO.Secondary,
		},
		addressSpace: jso.AddressSpace,
		subnets: &Subnets{
			virtualMachine: &Subnet{
				postfix:       jso.SubnetsJSO.VirtualMachine.Postfix,
				addressPrefix: jso.SubnetsJSO.VirtualMachine.AddressPrefix,
			},
			mongoDB: &Subnet{
				postfix:       jso.SubnetsJSO.MongoDB.Postfix,
				addressPrefix: jso.SubnetsJSO.MongoDB.AddressPrefix,
			},
		},
		virtualMachine: &VirtualMachine{
			size:               jso.VirtualMachineJSO.Size,
			storageAccountType: jso.VirtualMachineJSO.StorageAccountType,
			imageReference: &ImageReference{
				publisher: jso.VirtualMachineJSO.ImageReferenceJSO.Publisher,
				offer:     jso.VirtualMachineJSO.ImageReferenceJSO.Offer,
				sku:       jso.VirtualMachineJSO.ImageReferenceJSO.Sku,
				version:   jso.VirtualMachineJSO.ImageReferenceJSO.Version,
			},
		},
		databaseAccount: &DatabaseAccount{
			kind:                    jso.DatabaseAccountJSO.Kind,
			serverVersion:           jso.DatabaseAccountJSO.ServerVersion,
			offerType:               jso.DatabaseAccountJSO.OfferType,
			backupPolicyType:        jso.DatabaseAccountJSO.BackupPolicyType,
			defaultConsistencyLevel: jso.DatabaseAccountJSO.DefaultConsistencyLevel,
			capabilities:            jso.DatabaseAccountJSO.Capabilities,
		},
	}
}

func (cfg *Config) SubscriptionId() *string {
	return cfg.subscriptionId
}

func (cfg *Config) ProjectName() *string {
	return cfg.projectName
}

func (cfg *Config) Regions() *Regions {
	return cfg.regions
}

func (cfg *Config) AddressSpace() []*string {
	return cfg.addressSpace
}

func (cfg *Config) Subnets() *Subnets {
	return cfg.subnets
}

func (cfg *Config) VirtualMachine() *VirtualMachine {
	return cfg.virtualMachine
}

func (cfg *Config) DatabaseAccount() *DatabaseAccount {
	return cfg.databaseAccount
}

func (cfg *Config) Ids() *Ids {
	return cfg.ids
}

func (rgn *Regions) Primary() *string {
	return rgn.primary
}

func (rgn *Regions) Secondary() *string {
	return rgn.secondary
}

func (sbn *Subnets) VirtualMachine() *Subnet {
	return sbn.virtualMachine
}

func (sbn *Subnets) MongoDB() *Subnet {
	return sbn.mongoDB
}

func (sbn *Subnet) Postfix() *string {
	return sbn.postfix
}

func (sbn *Subnet) AddressPrefix() *string {
	return sbn.addressPrefix
}

func (img *ImageReference) Publisher() *string {
	return img.publisher
}

func (img *ImageReference) Offer() *string {
	return img.offer
}

func (img *ImageReference) Sku() *string {
	return img.sku
}

func (img *ImageReference) Version() *string {
	return img.version
}

func (vm *VirtualMachine) Size() *string {
	return vm.size
}

func (vm *VirtualMachine) StorageAccountType() *string {
	return vm.storageAccountType
}

func (vm *VirtualMachine) ImageReference() *ImageReference {
	return vm.imageReference
}

func (db *DatabaseAccount) Kind() *string {
	return db.kind
}

func (db *DatabaseAccount) ServerVersion() *string {
	return db.serverVersion
}

func (db *DatabaseAccount) OfferType() *string {
	return db.offerType
}

func (db *DatabaseAccount) BackupPolicyType() *string {
	return db.backupPolicyType
}

func (db *DatabaseAccount) DefaultConsistencyLevel() *string {
	return db.defaultConsistencyLevel
}

func (db *DatabaseAccount) Capabilities() []*string {
	return db.capabilities
}
