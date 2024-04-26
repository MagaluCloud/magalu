module magalu.cloud/core

go 1.21.4

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/PaesslerAG/gval v1.2.2
	github.com/PaesslerAG/jsonpath v0.1.1
	github.com/dustin/go-humanize v1.0.1
	github.com/getkin/kin-openapi v0.118.0
	github.com/go-openapi/jsonpointer v0.19.5
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0
	github.com/invopop/jsonschema v0.7.0
	github.com/invopop/yaml v0.1.0
	github.com/magiconair/properties v1.8.7
	github.com/mitchellh/mapstructure v1.5.0
	github.com/spf13/afero v1.11.0
	github.com/spf13/viper v1.16.0
	go.uber.org/zap v1.25.0
	golang.org/x/exp v0.0.0-20230817173708-d852ddb80c63
	golang.org/x/sync v0.5.0
	magalu.cloud/testing/fs_test_helper v0.0.0-00010101000000-000000000000
)

require (
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/perimeterx/marshmallow v1.1.4 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace magalu.cloud/testing/fs_test_helper => ../testing/fs_test_helper
