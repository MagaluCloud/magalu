package common

import (
	"encoding/xml"
	"time"
)

const namespace = "http://s3.amazonaws.com/doc/2006-03-01/"

type ObjectLockMode string

const (
	ObjectLockModeGovernance = ObjectLockMode("GOVERNANCE")
	ObjectLockModeCompliance = ObjectLockMode("COMPLIANCE")
)

// Object retention [Object]
type ObjectRetention struct {
	XMLName         xml.Name `xml:"Retention"`
	Namespace       string   `xml:"xmlns,attr"`
	Mode            ObjectLockMode
	RetainUntilDate string `xml:",omitempty"`
}

func DefaultObjectRetentionBody(retainUntilDate time.Time) ObjectRetention {
	return ObjectRetention{
		Namespace:       namespace,
		Mode:            ObjectLockModeCompliance,
		RetainUntilDate: retainUntilDate.UTC().Format(time.RFC3339),
	}

}

// Object lock [Bucket]
type ObjectLockRuleDefaultRetention struct {
	Days  int `xml:",omitempty"`
	Mode  ObjectLockMode
	Years int `xml:",omitempty"`
}

type ObjectLockRule struct {
	DefaultRetention ObjectLockRuleDefaultRetention
}

type ObjectLockingBody struct {
	XMLName           xml.Name `xml:"ObjectLockConfiguration,omitempty"`
	Namespace         string   `xml:"xmlns,attr"`
	ObjectLockEnabled string
	Rule              ObjectLockRule
}

var DefaultObjectLockingBody = ObjectLockingBody{
	ObjectLockEnabled: "Enabled",
	Namespace:         namespace,
	Rule: ObjectLockRule{
		DefaultRetention: ObjectLockRuleDefaultRetention{
			Mode: ObjectLockModeCompliance,
		},
	},
}
