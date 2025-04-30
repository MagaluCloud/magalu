package cors

import "encoding/xml"

type CORSConfiguration struct {
	XMLName   xml.Name   `xml:"CORSConfiguration"`
	XMLns     string     `xml:"xmlns,attr"`
	CORSRules []CORSRule `xml:"CORSRule" json:"CORSRules"`
}

type CORSRule struct {
	AllowedOrigins []string `xml:"AllowedOrigin" json:"AllowedOrigins"`
	AllowedMethods []string `xml:"AllowedMethod" json:"AllowedMethods"`
	AllowedHeaders []string `xml:"AllowedHeader,omitempty" json:"AllowedHeaders,omitempty"`
	ExposeHeaders  []string `xml:"ExposeHeader,omitempty" json:"ExposeHeaders,omitempty"`
	MaxAgeSeconds  int      `xml:"MaxAgeSeconds,omitempty" json:"MaxAgeSeconds,omitempty"`
}
