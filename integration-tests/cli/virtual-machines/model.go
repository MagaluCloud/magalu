package model

type InstanceType struct {
	Disk   int    `json:"disk"`
	Gpu    int    `json:"gpu"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Ram    int    `json:"ram"`
	Sku    string `json:"sku"`
	Status string `json:"status"`
	Vcpus  int    `json:"vcpus"`
}

type InstanceTypesResponse struct {
	InstanceTypes []InstanceType `json:"instance_types"`
}

type CreateNewInstanceResponse struct {
	ID string `json:"id"`
}

type GetVMInstanceResponse struct {
	CreatedAt   string  `json:"created_at"`
	ID          string  `json:"id"`
	Image       Image   `json:"image"`
	MachineType Machine `json:"machine_type"`
	Name        string  `json:"name"`
	Network     Network `json:"network"`
	SshKeyName  string  `json:"ssh_key_name"`
	State       string  `json:"state"`
	Status      string  `json:"status"`
	UpdatedAt   string  `json:"updated_at"`
}

type Image struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
}

type Machine struct {
	Disk  int    `json:"disk"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	RAM   int    `json:"ram"`
	Vcpus int    `json:"vcpus"`
}

type Network struct {
	Ports []Port `json:"ports"`
	VPC   VPC    `json:"vpc"`
}

type Port struct {
	ID          string      `json:"id"`
	IpAddresses IpAddresses `json:"ipAddresses"`
	Name        string      `json:"name"`
}

type IpAddresses struct {
	IpV6Address      string `json:"ipV6Address"`
	PrivateIpAddress string `json:"privateIpAddress"`
	PublicIpAddress  string `json:"publicIpAddress"`
}

type VPC struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ImageResponse struct {
	Images []ImageDetails `json:"images"`
}

type ImageDetails struct {
	EndLifeAt            string              `json:"end_life_at"`
	EndStandardSupportAt string              `json:"end_standard_support_at"`
	ID                   string              `json:"id"`
	MinimumRequirements  MinimumRequirements `json:"minimum_requirements"`
	Name                 string              `json:"name"`
	Platform             string              `json:"platform"`
	ReleaseAt            string              `json:"release_at"`
	Status               string              `json:"status"`
	Version              string              `json:"version"`
}

type MinimumRequirements struct {
	Disk int `json:"disk"`
	Ram  int `json:"ram"`
	Vcpu int `json:"vcpu"`
}
