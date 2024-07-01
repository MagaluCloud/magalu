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
