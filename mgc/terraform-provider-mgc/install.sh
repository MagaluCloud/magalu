#!/bin/bash

config_file=~/.terraformrc
registry=registry.terraform.io
if [[ "$MGC_OPENTF" ]]; then
    config_file=~/.tofurc
    registry=registry.opentofu.org
fi

cat > $config_file <<- EOM
provider_installation {
  dev_overrides {
    "$registry/magalucloud/mgc" = "$PWD"
  }

  direct {}
}
EOM

echo "File $config_file writen to $HOME:"
cat $config_file
