/*
Package: config

# Summary

# Manage CLI Configuration values

# Description

Configuration values are available to be set so that they persist between
different executions of the MgcSDK. They reside in a YAML file when set.
Config values may also be loaded via Environment Variables. Any Config available
(see 'list') may be exported as an env variable in uppercase with the 'MGC_' prefix

import "magalu.cloud/lib/products/config"
*/
package config
