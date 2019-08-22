#!/bin/bash
set -e
set -x

export LOCAL_IP=$(ip route get 8.8.8.8 | grep -oE 'src ([0-9\.]+)' | cut -d ' ' -f 2)
if [ "$SERVER_NAME" == "" ]; then
    SERVER_NAME=$LOCAL_IP
fi

hostip=$(echo -n $(ip addr show eth0 | grep 'inet' | cut -d " " -f 6 | cut -d "/" -f 1 | awk '{ print $1 }'))
echo $hostip
if [ "$REGISTRY_ETCD_URL" != "" ]; then
    echo "Will register this generator instance. --etcd-url=$REGISTRY_ETCD_URL --etcd-base=$REGISTRY_ETCD_BASE --service=$REGISTRY_SERVICE --name=$(hostname):9090 --ttl=$REGISTRY_TTL"
    etcd-registrar \
        --loglevel=info \
        --etcd-url=$REGISTRY_ETCD_URL \
        --etcd-base=$REGISTRY_ETCD_BASE \
        --service=$REGISTRY_SERVICE \
        --name="$hostip:3000" \
        --ttl=$REGISTRY_TTL&
fi

echo "Starting the almighty Metrics Generator Tabajara..."
metrics-generator-tabajara \
    --server-name=${SERVER_NAME} \
    --component-name=${COMPONENT_NAME} \
    --component-version=${COMPONENT_VERSION} \
    --accident-resource="${ACCIDENT_RESOURCE}" \
    --accident-ratio="${ACCIDENT_RATIO}" \
    --accident-type="$ACCIDENT_TYPE"
