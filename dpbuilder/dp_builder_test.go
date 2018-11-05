package dpbuilder_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"

	artifacts "github.com/kun-lun/artifacts/pkg/apis"
	. "github.com/kun-lun/deployment-producer/dpbuilder"
)

var _ = Describe("DpBuilder", func() {
	var (
		m *artifacts.Manifest
	)
	BeforeEach(func() {

		platform := artifacts.Platform{
			Type: "php",
		}

		networks := []artifacts.VirtualNetwork{
			{
				Name: "vnet-1",
				Subnets: []artifacts.Subnet{
					{
						Range:   "10.10.0.0/24",
						Gateway: "10.10.0.1",
						Name:    "snet-1",
					},
				},
			}}

		loadBalancers := []artifacts.LoadBalancer{
			{
				Name: "kunlun-wenserver-lb",
				SKU:  "standard",
			},
		}

		networkSecurityGroups := []artifacts.NetworkSecurityGroup{
			{
				Name: "nsg_1",
				NetworkSecurityRules: []artifacts.NetworkSecurityRule{
					{
						Name:                     "allow-ssh",
						Priority:                 100,
						Direction:                "Inbound",
						Access:                   "Allow",
						Protocol:                 "Tcp",
						SourcePortRange:          "*",
						DestinationPortRange:     "22",
						SourceAddressPrefix:      "*",
						DestinationAddressPrefix: "*",
					},
				},
			},
		}

		vmGroups := []artifacts.VMGroup{
			{
				Name: "jumpbox",
				Meta: yaml.MapSlice{
					{
						Key:   "group_type",
						Value: "jumpbox",
					},
				},
				SKU:   artifacts.VMStandardDS1V2,
				Count: 1,
				Type:  "VM",
				Storage: &artifacts.VMStorage{
					Image: &artifacts.Image{
						Offer:     "offer1",
						Publisher: "ubuntu",
						SKU:       "sku1",
						Version:   "latest",
					},
					OSDisk: &artifacts.OSDisk{},
					DataDisks: []artifacts.DataDisk{
						{
							DiskSizeGB: 10,
						},
					},
					AzureFiles: []artifacts.AzureFile{
						{
							StorageAccount: "storage_account_1",
							Name:           "azure_file_1",
							MountPoint:     "/mnt/azurefile_1",
						},
					},
				},
				OSProfile: artifacts.VMOSProfile{
					AdminName: "kunlun",
				},
				NetworkInfos: []artifacts.VMNetworkInfo{
					{
						SubnetName:               networks[0].Subnets[0].Name,
						LoadBalancerName:         loadBalancers[0].Name,
						NetworkSecurityGroupName: networkSecurityGroups[0].Name,
						PublicIP:                 "dynamic",
						Outputs: []artifacts.VMNetworkOutput{
							{
								IP:       "172.16.8.4",
								PublicIP: "13.75.71.162",
								Host:     "andliuubuntu.eastasia.cloudapp.azure.com",
							},
						},
					},
				},
				Roles: []artifacts.Role{
					{
						Name: "builtin/jumpbox",
					},
				},
			},
			{
				Name:  "d2v3_group",
				SKU:   artifacts.VMStandardDS1V2,
				Count: 2,
				Type:  "VM",
				OSProfile: artifacts.VMOSProfile{
					AdminName: "kunlun",
				},
				Storage: &artifacts.VMStorage{
					OSDisk: &artifacts.OSDisk{},
					DataDisks: []artifacts.DataDisk{
						{
							DiskSizeGB: 10,
						},
					},
					AzureFiles: []artifacts.AzureFile{},
				},
				NetworkInfos: []artifacts.VMNetworkInfo{
					{
						SubnetName:       networks[0].Subnets[0].Name,
						LoadBalancerName: loadBalancers[0].Name,
						Outputs: []artifacts.VMNetworkOutput{
							{
								IP: "172.16.8.4",
							},
							{
								IP: "172.16.8.4",
							},
						},
					},
				},
				Roles: []artifacts.Role{
					{
						Name: "builtin/php_web_role",
					},
				},
			},
		}

		storageAccounts := []artifacts.StorageAccount{
			{
				Name:     "storage_account_1",
				SKU:      "standard",
				Location: "eastus",
			},
		}

		databases := []artifacts.Database{
			{
				MigrationInformation: &artifacts.MigrationInformation{
					OriginHost:     "asd",
					OriginDatabase: "asd",
					OriginUsername: "asd",
					OriginPassword: "asd",
				},
				Engine:              artifacts.MysqlDB,
				EngineVersion:       "5.7",
				Cores:               2,
				Storage:             5,
				BackupRetentionDays: 35,
				Username:            "dbuser",
				Password:            "abcd1234!",
			},
		}

		// The checker add needed resource to manifest
		m = &artifacts.Manifest{
			Schema:                "v0.1",
			IaaS:                  "azure",
			Location:              "eastus",
			Platform:              &platform,
			VMGroups:              vmGroups,
			VNets:                 networks,
			LoadBalancers:         loadBalancers,
			StorageAccounts:       storageAccounts,
			NetworkSecurityGroups: networkSecurityGroups,
			Databases:             databases,
		}

	})
	Describe("Produce", func() {
		Context("Everything OK", func() {
			It("should produce deployments and hosts correctly", func() {
				dpProducer := DeploymentBuilder{}
				x, y, err := dpProducer.Produce(*m)
				Expect(err).To(BeNil())
				x_str, err := json.Marshal(x)
				Expect(err).To(BeNil())
				println(string(x_str))
				y_str, err := json.Marshal(y)
				Expect(err).To(BeNil())
				println(string(y_str))
			})
		})
	})
})
