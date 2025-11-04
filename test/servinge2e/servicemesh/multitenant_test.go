package servicemesh

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	// needed to serialize the Kubernetes object
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	"github.com/openshift-knative/serverless-operator/test"
	"github.com/openshift-knative/serverless-operator/test/servinge2e"
	"knative.dev/networking/pkg/apis/networking"
	pkgTest "knative.dev/pkg/test"
	"knative.dev/pkg/test/spoof"
	"knative.dev/serving/pkg/apis/autoscaling"
	"knative.dev/serving/pkg/apis/serving"
	servingTest "knative.dev/serving/test"
)

const (
	Tenant1          = "tenant-1"
	Tenant2          = "tenant-2"
	LocalGatewayHost = "knative-local-gateway.knative-serving-ingress.svc.cluster.local"
)

var ExpectStatusForbidden = func(resp *spoof.Response) (bool, error) {
	if resp.StatusCode != 403 {
		// Returning (false, nil) causes SpoofingClient.Poll to retry.
		return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return true, nil
}

func TestMultiTenancyWithServiceMesh(t *testing.T) {
	tests := []testCase{
		//{
		//	name: "same-tenant-directly",
		//	annotations: map[string]string{
		//		autoscaling.TargetBurstCapacityKey: "0",
		//		autoscaling.MinScaleAnnotationKey:  "1",
		//	},
		//	sourceNamespace:   Tenant1,
		//	targetNamespace:   Tenant1,
		//	usePrivateService: true,
		//	checkResponseFunc: spoof.MatchesBody(servinge2e.HelloworldGoText),
		//},
		//{
		//	name: "cross-tenant-directly",
		//	annotations: map[string]string{
		//		autoscaling.TargetBurstCapacityKey: "0",
		//		autoscaling.MinScaleAnnotationKey:  "1",
		//	},
		//	sourceNamespace:   Tenant1,
		//	targetNamespace:   Tenant2,
		//	usePrivateService: true,
		//	checkResponseFunc: ExpectStatusForbidden,
		//},
		//{
		//	name: "same-tenant-via-activator",
		//	annotations: map[string]string{
		//		autoscaling.TargetBurstCapacityKey: "-1",
		//	},
		//	sourceNamespace:   Tenant1,
		//	targetNamespace:   Tenant1,
		//	checkResponseFunc: spoof.MatchesBody(servinge2e.HelloworldGoText),
		//},
		//{
		//	name: "cross-tenant-via-activator",
		//	annotations: map[string]string{
		//		autoscaling.TargetBurstCapacityKey: "-1",
		//	},
		//	sourceNamespace:   Tenant1,
		//	targetNamespace:   Tenant2,
		//	checkResponseFunc: ExpectStatusForbidden,
		//},
		//{
		//	name: "same-tenant-via-ingress-via-activator",
		//	annotations: map[string]string{
		//		autoscaling.TargetBurstCapacityKey: "-1",
		//	},
		//	sourceNamespace:   Tenant1,
		//	targetNamespace:   Tenant1,
		//	checkResponseFunc: spoof.MatchesBody(servinge2e.HelloworldGoText),
		//	gateway:           LocalGatewayHost,
		//},
		//{
		//	name: "cross-tenant-via-ingress-via-activator",
		//	annotations: map[string]string{
		//		autoscaling.TargetBurstCapacityKey: "-1",
		//	},
		//	sourceNamespace:   Tenant1,
		//	targetNamespace:   Tenant2,
		//	checkResponseFunc: ExpectStatusForbidden,
		//	gateway:           LocalGatewayHost,
		//},
		//{
		//	name: "same-tenant-via-ingress-no-activator",
		//	annotations: map[string]string{
		//		autoscaling.TargetBurstCapacityKey: "0",
		//		autoscaling.MinScaleAnnotationKey:  "1",
		//	},
		//	sourceNamespace:   Tenant1,
		//	targetNamespace:   Tenant1,
		//	checkResponseFunc: spoof.MatchesBody(servinge2e.HelloworldGoText),
		//	gateway:           LocalGatewayHost,
		//},
		{
			name: "cross-tenant-via-ingress-no-activator",
			annotations: map[string]string{
				autoscaling.TargetBurstCapacityKey: "0",
				autoscaling.MinScaleAnnotationKey:  "1",
			},
			sourceNamespace:   Tenant1,
			targetNamespace:   Tenant2,
			checkResponseFunc: ExpectStatusForbidden,
			gateway:           LocalGatewayHost,
		}}

	for _, tc := range tests {
		tc := tc

		tc.annotations[IstioRewriteProbersKey] = "true"

		// Always use cluster-local service.
		tc.labels = map[string]string{
			networking.VisibilityLabelKey: serving.VisibilityClusterLocal,
			IstioInjectKey:                "true",
		}
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := test.SetupClusterAdmin(t)
			test.CleanupOnInterrupt(t, func() { test.CleanupAll(t, ctx) })
			defer test.CleanupAll(t, ctx)

			service := test.Service(tc.name, tc.targetNamespace, pkgTest.ImagePath(test.HelloworldGoImg), map[string]string{
				ServingEnablePassthroughKey: "true",
			}, tc.annotations)
			service.ObjectMeta.Labels = tc.labels
			// ==========================================================
			// 👇 START OF DEBUG CODE TO ADD 👇
			// ==========================================================
			// Assuming 't' is your *testing.T object from the test function argument

			// 1. Create a serializer (using JSON is often easier for quick logging)
			// You may need to adapt this based on the project's scheme.
			scheme := runtime.NewScheme()
			// Add the Knative Serving scheme to allow serialization
			// servingv1.AddToScheme(scheme)

			// 2. Serialize the service object to JSON bytes
			serializer := json.NewSerializer(json.DefaultMetaFactory, scheme, scheme, true)
			var buffer bytes.Buffer
			if err := serializer.Encode(service, &buffer); err == nil {
				// 3. Print the YAML/JSON payload to the test log
				t.Logf("--- DEBUG: KService YAML Payload for %s ---\n%s\n--- END KSERVICE PAYLOAD ---", tc.name, buffer.String())
			} else {
				t.Logf("Failed to serialize KService for debugging: %v", err)
			}
			// ==========================================================
			// 👆 END OF DEBUG CODE TO ADD 👆
			// ==========================================================
			service = test.WithServiceReadyOrFail(ctx, service)

			gateway := ""
			targetHost := service.Status.URL.Host
			if tc.usePrivateService {
				targetHost = fmt.Sprintf("%s-00001-private.%s.svc.cluster.local", service.Name, tc.targetNamespace)
			} else if tc.gateway != "" {
				gateway = tc.gateway
				targetHost = service.Status.URL.Host
			}

			httpProxy := test.WithServiceReadyOrFail(ctx, servinge2e.HTTPProxyService(tc.name+"-proxy", tc.sourceNamespace, gateway, targetHost, "", map[string]string{
				ServingEnablePassthroughKey: "true",
			}, tc.annotations))

			if _, err := pkgTest.CheckEndpointState(
				context.Background(),
				ctx.Clients.Kube,
				t.Logf,
				httpProxy.Status.URL.URL(),
				tc.checkResponseFunc,
				"CheckResponse",
				true,
				servingTest.AddRootCAtoTransport(context.Background(), t.Logf, &servingTest.Clients{KubeClient: ctx.Clients.Kube}, true),
			); err != nil {
				t.Fatalf("Unexpected state for %s :%v", httpProxy.Status.URL.URL(), err)
			}
		})
	}
}
