package client

import (
	mgcUtils "magalu.cloud/core/utils"
	mgcSdk "magalu.cloud/sdk"
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
