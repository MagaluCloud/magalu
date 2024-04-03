package common

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"github.com/google/uuid"
	"magalu.cloud/core"
)

type AccessControlPolicy struct {
	Owner             Owner             `xml:"Owner"`
	AccessControlList AccessControlList `xml:"AccessControlList"`

	XMLName struct{} `xml:"AccessControlPolicy" json:"-"`
}

type AccessControlList struct {
	Grant Grant `xml:"Grant"`
}

type Grant struct {
	Grantee    Grantee `xml:"Grantee"`
	Permission string  `xml:"Permission"`
}

type Grantee struct {
	DisplayName  string `xml:"DisplayName"`
	EmailAddress string `xml:"EmailAddress"`
	ID           string `xml:"ID"`
	URI          string `xml:"URI"`
}

type ACLPermission struct {
	ID string `json:"id" jsonschema:"description=Either a Tenant ID or a User Project ID,example=a4900b57-7dbb-4906-b7e8-efed938e325c"`
}

func (g ACLPermission) Validate() error {
	if g.ID == "" {
		return fmt.Errorf("ID for ACLPermission may not be empty")
	}
	return nil
}

type ACLPermissions struct {
	ACLStandardPermissions `json:",squash"` // nolint
	ACLCannedPermissions   `json:",squash"` // nolint
}

func (p ACLPermissions) IsEmpty() bool {
	return p.ACLStandardPermissions.IsEmpty() && p.ACLCannedPermissions.IsEmpty()
}

func (p ACLPermissions) Validate() error {
	err := p.ACLCannedPermissions.Validate()
	if err != nil {
		return err
	}

	err = p.ACLStandardPermissions.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (p ACLPermissions) SetHeaders(req *http.Request, cfg Config) error {
	err := p.ACLStandardPermissions.SetHeaders(req, cfg)
	if err != nil {
		return err
	}

	p.ACLCannedPermissions.SetHeaders(req)
	return nil
}

type ACLStandardPermissions struct {
	GrantFullControl []ACLPermission `json:"grant_full_control,omitempty" jsonschema:"description=Grantees get FULL_CONTROL" mgc:"hidden"`
	GrantRead        []ACLPermission `json:"grant_read,omitempty" jsonschema:"description=Allows grantees to list the objects in the bucket" mgc:"hidden"`
	GrantWrite       []ACLPermission `json:"grant_write,omitempty" jsonschema:"description=Allows grantees to create objects in the bucket"`
	GrantReadAcp     []ACLPermission `json:"grant_read_acp,omitempty" jsonschema:"description=Allows grantees to read the bucket ACL" mgc:"hidden"`
	GrantWriteAcp    []ACLPermission `json:"grant_write_acp,omitempty" jsonschema:"description=Allows grantees to write the ACL for the applicable bucket" mgc:"hidden"`
}

func (o ACLStandardPermissions) Validate() error {
	allPermissions := [][]ACLPermission{o.GrantFullControl, o.GrantRead, o.GrantWrite, o.GrantReadAcp, o.GrantWriteAcp}
	for _, permissions := range allPermissions {
		for _, permission := range permissions {
			err := permission.Validate()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p ACLStandardPermissions) IsEmpty() bool {
	nonEmptyFields, _ := collectSliceStructNonEmptyFields(&p)
	return len(nonEmptyFields) < 1
}

func (p ACLStandardPermissions) SetHeaders(req *http.Request, cfg Config) (err error) {
	err = p.setHeader(req.Header, "x-amz-grant-full-control", p.GrantFullControl, cfg)
	if err != nil {
		return err
	}

	err = p.setHeader(req.Header, "x-amz-grant-read", p.GrantRead, cfg)
	if err != nil {
		return err
	}

	err = p.setHeader(req.Header, "x-amz-grant-write", p.GrantWrite, cfg)
	if err != nil {
		return err
	}

	err = p.setHeader(req.Header, "x-amz-grant-read-acp", p.GrantReadAcp, cfg)
	if err != nil {
		return err
	}

	err = p.setHeader(req.Header, "x-amz-grant-write-acp", p.GrantWriteAcp, cfg)
	if err != nil {
		return err
	}

	return nil
}

var userProjectRegex = regexp.MustCompile("cloud_(?P<region>[^_]+)_(?P<env>[^_]+).(?P<tenant_id>[^:]+):cloud_(?P<region1>[^_]+)_(?P<env1>[^_]+).(?P<tenant_id1>[^:]+)")

func (p ACLStandardPermissions) tenantIdFromUserProject(userProject string) (string, error) {
	match := userProjectRegex.FindStringSubmatch(userProject)
	for i, substr := range match {
		if userProjectRegex.SubexpNames()[i] == "tenant_id" {
			return substr, nil
		}
	}

	return "", core.UsageError{
		Err: fmt.Errorf("unable to find 'tenant_id' inside 'user_project' ACL ID permission: %q", userProject),
	}
}

func (p ACLStandardPermissions) userProjectFromTenantId(tenantId string, cfg Config) string {
	pattern := fmt.Sprintf("cloud_%s_prod_%s", cfg.translateRegion(), tenantId)
	return fmt.Sprintf("%s:%s", pattern, pattern)
}

func (p ACLStandardPermissions) setHeader(header http.Header, name string, permissions []ACLPermission, cfg Config) error {
	v := ""
	for i, permission := range permissions {
		if i > 0 {
			v += ","
		}

		var tenantId, userProject string

		_, err := uuid.Parse(permission.ID)
		isTenantId := err == nil

		if isTenantId {
			tenantId = permission.ID
			userProject = p.userProjectFromTenantId(permission.ID, cfg)
		} else {
			userProject = permission.ID
			tenantId, err = p.tenantIdFromUserProject(permission.ID)
			if err != nil {
				return err
			}
		}

		v += fmt.Sprintf("id=%s,id=%s", tenantId, userProject)
	}

	if v == "" {
		return nil
	}

	header.Set(name, v)
	return nil
}

type ACLCannedPermissions struct {
	Private           bool `json:"private,omitempty" jsonschema:"description=Owner gets FULL_CONTROL. Delegated users have access. No one else has access rights"`
	PublicRead        bool `json:"public_read,omitempty" jsonschema:"description=Owner gets FULL_CONTROL. Everyone else has READ rights"`
	PublicReadWrite   bool `json:"public_read_write,omitempty" jsonschema:"description=Owner gets FULL_CONTROL. Everyone else has READ and WRITE rights" mgc:"hidden"`
	AuthenticatedRead bool `json:"authenticated_read,omitempty" jsonschema:"description=Owner gets FULL_CONTROL. Authenticated users have READ rights" mgc:"hidden"`
	AwsExecRead       bool `json:"aws_exec_read,omitempty" mgc:"hidden"`
}

func (o ACLCannedPermissions) Validate() error {
	trueFields, err := collectBoolStructTrueFields(&o)
	if err != nil {
		return err
	}

	if len(trueFields) > 1 {
		return core.UsageError{Err: fmt.Errorf("canned ACL cannot have more than one 'true' field: %v", trueFields)}
	}

	return nil
}

func (o ACLCannedPermissions) IsEmpty() bool {
	trueFields, _ := collectBoolStructTrueFields(&o)
	return len(trueFields) < 1
}

func (o ACLCannedPermissions) SetHeaders(req *http.Request) {
	if o.PublicRead {
		req.Header.Set("x-amz-acl", "public-read")
		return
	}

	if o.PublicReadWrite {
		req.Header.Set("x-amz-acl", "public-read-write")
		return
	}

	if o.AuthenticatedRead {
		req.Header.Set("x-amz-acl", "authenticated-read")
		return
	}

	if o.Private {
		req.Header.Set("x-amz-acl", "private")
		return
	}
	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObjectAcl.html
	if o.AwsExecRead {
		req.Header.Set("x-amz-acl", "aws-exec-read")
	}
}

func collectSliceStructNonEmptyFields[T any](o *T) (nonEmptyFields []string, err error) {
	return collectStructFields(o, func(name string, value reflect.Value) (bool, error) {
		kind := value.Kind()
		if kind != reflect.Slice && kind != reflect.Array {
			return false, fmt.Errorf("programming error: slice struct %T cannot have non-slice fields: %q is %s", *o, name, kind.String())
		}

		return value.Len() > 0, nil
	})
}

func collectBoolStructTrueFields[T any](o *T) (trueFields []string, err error) {
	return collectStructFields(o, func(name string, value reflect.Value) (bool, error) {
		kind := value.Kind()
		if kind != reflect.Bool {
			return false, fmt.Errorf("programming error: boolean struct %T cannot have non-boolean fields: %q is %s", *o, name, kind.String())
		}

		return value.Bool(), nil
	})
}

func collectStructFields[T any](o *T, include func(name string, value reflect.Value) (bool, error)) (fields []string, err error) {
	v := reflect.ValueOf(*o)
	t := v.Type()

	if t.Kind() != reflect.Struct {
		err = fmt.Errorf("programming error: 'collectStructFields' only accepts structs")
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := t.Field(i).Name

		var shouldInclude bool
		shouldInclude, err = include(name, field)
		if err != nil {
			return
		}

		if shouldInclude {
			fields = append(fields, name)
		}
	}
	return
}
