#!/bin/bash
set -e

CLUSTER_NAME="radius-sandbox"

if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    echo "Cluster ${CLUSTER_NAME} already exists"
    exit 0
fi

echo "Creating cluster ${CLUSTER_NAME}..."
kind create cluster --name "${CLUSTER_NAME}" --config kind-config.yaml --quiet

echo "Configuring CA certificates..."
docker exec "${CLUSTER_NAME}-control-plane" update-ca-certificates 2>/dev/null || true
docker exec "${CLUSTER_NAME}-control-plane" systemctl restart containerd 2>/dev/null || true

kubectl cluster-info --context "kind-${CLUSTER_NAME}" >/dev/null
echo "Cluster ${CLUSTER_NAME} created"
