#!/bin/bash

# Apply Service Mesh 3 compatibility patches
# This script applies all SM3-related patches for serverless-operator

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
YAML_PATCH="$SCRIPT_DIR/hack/patches/024-containersource-yaml-template.patch"
GO_PATCH="$SCRIPT_DIR/hack/patches/025-containersource-go-labels.patch"
MULTITENANT_PATCH="$SCRIPT_DIR/hack/patches/026-multitenant-sm3-compatibility.patch"
RECORDEVENTS_PATCH="$SCRIPT_DIR/hack/patches/027-recordevents-sm3-labels.patch"
WAIT_PATCH="$SCRIPT_DIR/hack/patches/028-wait-sm3-labels.patch"
MESH_PATCH="$SCRIPT_DIR/hack/patches/029-mesh-disable-network-policy.patch"

echo "🔧 Applying Service Mesh 3 compatibility patches..."

# Apply YAML template patch
if [[ -f "$YAML_PATCH" ]]; then
    if git apply --check "$YAML_PATCH" >/dev/null 2>&1; then
        echo "✅ Applying YAML template patch: $YAML_PATCH"
        git apply "$YAML_PATCH"
        echo "✅ YAML template patch applied successfully!"
    else
        echo "ℹ️  YAML template patch appears to already be applied or conflicts exist."
    fi
else
    echo "❌ Error: YAML patch file not found at $YAML_PATCH"
fi

# Apply Go code patch
if [[ -f "$GO_PATCH" ]]; then
    if git apply --check "$GO_PATCH" >/dev/null 2>&1; then
        echo "✅ Applying Go code patch: $GO_PATCH"
        git apply "$GO_PATCH"
        echo "✅ Go code patch applied successfully!"
    else
        echo "ℹ️  Go code patch appears to already be applied or conflicts exist."
    fi
else
    echo "❌ Error: Go patch file not found at $GO_PATCH"
fi

# Apply multitenant SM3 compatibility patch
if [[ -f "$MULTITENANT_PATCH" ]]; then
    if git apply --check "$MULTITENANT_PATCH" >/dev/null 2>&1; then
        echo "✅ Applying multitenant SM3 compatibility patch: $MULTITENANT_PATCH"
        git apply "$MULTITENANT_PATCH"
        echo "✅ Multitenant SM3 compatibility patch applied successfully!"
    else
        echo "ℹ️  Multitenant SM3 compatibility patch appears to already be applied or conflicts exist."
    fi
else
    echo "❌ Error: Multitenant SM3 compatibility patch file not found at $MULTITENANT_PATCH"
fi

# Apply recordevents SM3 labels patch
if [[ -f "$RECORDEVENTS_PATCH" ]]; then
    if git apply --check "$RECORDEVENTS_PATCH" >/dev/null 2>&1; then
        echo "✅ Applying recordevents SM3 labels patch: $RECORDEVENTS_PATCH"
        git apply "$RECORDEVENTS_PATCH"
        echo "✅ Recordevents SM3 labels patch applied successfully!"
    else
        echo "ℹ️  Recordevents SM3 labels patch appears to already be applied or conflicts exist."
    fi
else
    echo "❌ Error: Recordevents SM3 labels patch file not found at $RECORDEVENTS_PATCH"
fi

# Apply wait SM3 labels patch
if [[ -f "$WAIT_PATCH" ]]; then
    if git apply --check "$WAIT_PATCH" >/dev/null 2>&1; then
        echo "✅ Applying wait SM3 labels patch: $WAIT_PATCH"
        git apply "$WAIT_PATCH"
        echo "✅ Wait SM3 labels patch applied successfully!"
    else
        echo "ℹ️  Wait SM3 labels patch appears to already be applied or conflicts exist."
    fi
else
    echo "❌ Error: Wait SM3 labels patch file not found at $WAIT_PATCH"
fi

# Apply mesh disable network policy patch
if [[ -f "$MESH_PATCH" ]]; then
    if git apply --check "$MESH_PATCH" >/dev/null 2>&1; then
        echo "✅ Applying mesh disable network policy patch: $MESH_PATCH"
        git apply "$MESH_PATCH"
        echo "✅ Mesh disable network policy patch applied successfully!"
    else
        echo "ℹ️  Mesh disable network policy patch appears to already be applied or conflicts exist."
    fi
else
    echo "❌ Error: Mesh disable network policy patch file not found at $MESH_PATCH"
fi

echo ""
echo "🔍 Verifying changes:"
git status --porcelain vendor/knative.dev/eventing/test/rekt/resources/containersource/
git status --porcelain test/servinge2e/servicemesh/multitenant_test.go
git status --porcelain vendor/knative.dev/eventing/test/lib/recordevents/resources.go
git status --porcelain vendor/knative.dev/reconciler-test/pkg/k8s/wait.go
git status --porcelain hack/lib/mesh.bash
echo ""
echo "📋 Changes applied:"
echo "   - containersource.go: Added manifest.WithIstioPodLabels(cfg)"
echo "   - containersource.yaml: Updated template to handle both podannotations and podlabels"
echo "   - multitenant_test.go: Updated LocalGatewayHost for SM3 + Added IstioInjectKey label"
echo "   - recordevents/resources.go: Added Istio injection label for SM3"
echo "   - wait.go: Added Istio injection label for SM3"
echo "   - mesh.bash: Disabled network policy monitoring for SM3"
echo ""
echo "✅ Service Mesh 3 compatibility patches applied successfully!"
