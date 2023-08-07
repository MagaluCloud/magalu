#!/bin/bash
set -xe

MGC_CLI=${MGC_CLI:-./mgc}

# 1. Login
$MGC_CLI auth login

# 2. Create Keypair for your SSH, if more than one pub key it will fail
SSH_KEY_NAME="example-key";
ssh-keygen -t ed25519 -N "" -f /tmp/$SSH_KEY_NAME
$MGC_CLI virtual-machine keypairs create --name=$SSH_KEY_NAME --public-key="$(cat /tmp/$SSH_KEY_NAME.pub)"
rm /tmp/$SSH_KEY_NAME

# 3. Create VM
IMAGE="cloud-debian-11 LTS"
TYPE="cloud-bs1.xsmall"
INSTANCE_NAME="vm-example-1"
read VM_ID < <($MGC_CLI virtual-machine instances create \
    --image="$IMAGE" \
    --type=$TYPE \
    --key_name=$SSH_KEY_NAME \
    --name=$INSTANCE_NAME -o jsonpath='$.id')

# 4. Create Disk
DESCRIPTION="example-volume"
DISK_NAME="example-volume";
DISK_TYPE="cloud_nvme";
DISK_SIZE=1;
read DISK_ID < <($MGC_CLI block-storage volume create \
    --description=$DESCRIPTION \
    --name=$DISK_NAME \
    --volume-type=$DISK_TYPE \
    --size=$DISK_SIZE -o jsonpath='$.id')

# 5. Wait for the VM to transition to active
CUR_STATUS=""
DESIRED_STATUS="ACTIVE"
while [ [${CUR_STATUS}] != [\"${DESIRED_STATUS}\"] ]
do
    CUR_STATUS=$($MGC_CLI virtual-machine instances get --id=$VM_ID 2>/dev/null -o jsonpath='$.status')
    sleep 1
done

# 6. Attach Disk to VM - may fail if VM is in Pending status
$MGC_CLI block-storage volume attach \
    --id=$DISK_ID \
    --virtual-machine-id=$VM_ID

# 7. Shutoff VM
$MGC_CLI virtual-machine instances update --id=$VM_ID --status="shutoff"
