package client

import (
	mgcUtils "github.com/MagaluCloud/magalu/mgc/core/utils"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
)

type Client struct {
	sdk *mgcSdk.Sdk
}

var DefaultSdk = mgcUtils.NewLazyLoader(func() *mgcSdk.Sdk {
	return mgcSdk.NewSdk()
})

func (c *Client) Sdk() *mgcSdk.Sdk {
	if c == nil || c.sdk == nil {
		return DefaultSdk()
	}
	return c.sdk
}

// Creates a new Client based on the given SDK.
//
// If sdk is nil, then the DefaultSdk() is used.
func NewClient(sdk *mgcSdk.Sdk) *Client {
	return &Client{
		sdk: sdk,
	}
}

// String is a helper routine returns a pointer to the string value passed in.
func String(s string) *string {
	return &s
}

// Boolean is a helper routine returns a pointer to the bool value passed in.
func Boolean(b bool) *bool {
	return &b
}

// Int is a helper routine returns a pointer to the int value passed in.
func Int(i int) *int {
	return &i
}

// Int64 is a helper routine returns a pointer to the int64 value passed in.
func Int64(i int64) *int64 {
	return &i
}

// Float64 is a helper routine returns a pointer to the float64 value passed in.
func Float64(f float64) *float64 {
	return &f
}

// Float32 is a helper routine returns a pointer to the float32 value passed in.
func Float32(f float32) *float32 {
	return &f
}

// Uint is a helper routine returns a pointer to the uint value passed in.
func Uint(u uint) *uint {
	return &u
}

// Uint64 is a helper routine returns a pointer to the uint64 value passed in.
func Uint64(u uint64) *uint64 {
	return &u
}

// Uint32 is a helper routine returns a pointer to the uint32 value passed in.
func Uint32(u uint32) *uint32 {
	return &u
}

// Uint16 is a helper routine returns a pointer to the uint16 value passed in.
func Uint16(u uint16) *uint16 {
	return &u
}

// Uint8 is a helper routine returns a pointer to the uint8 value passed in.
func Uint8(u uint8) *uint8 {
	return &u
}

// Int32 is a helper routine returns a pointer to the int32 value passed in.
func Int32(i int32) *int32 {
	return &i
}

// Int16 is a helper routine returns a pointer to the int16 value passed in.
func Int16(i int16) *int16 {
	return &i
}

// Int8 is a helper routine returns a pointer to the int8 value passed in.
func Int8(i int8) *int8 {
	return &i
}
