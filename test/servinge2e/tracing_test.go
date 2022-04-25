package servinge2e

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	jaegermodel "github.com/jaegertracing/jaeger/model"
	jaegerapi "github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"github.com/openshift-knative/serverless-operator/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"knative.dev/pkg/test/spoof"
)

const (
	testNamespace    = "serverless-tests"
	requestCount     = 100
	tracingNamespace = "istio-system"
	jaegerName       = "jaeger"
)

func TestTraceStartedAtActivator(t *testing.T) {
	tracingTest(t, true /* activatorInPath */)
}

func TestTraceStartedAtQueueProxy(t *testing.T) {
	tracingTest(t, false /* activatorInPath */)
}

func tracingTest(t *testing.T, activatorInPath bool) {
	ctx := test.SetupClusterAdmin(t)
	if test.IsServiceMeshInstalled(ctx) {
		// Traces look different when ServiceMesh is installed.
		t.Skip("ServiceMesh installed, skipping tracing test.")
	}
	test.CleanupOnInterrupt(t, func() { test.CleanupAll(t, ctx) })
	defer test.CleanupAll(t, ctx)
	name := strings.ToLower(t.Name())
	annotations := map[string]string{
		"autoscaling.knative.dev/targetBurstCapacity": "0",
	}
	if activatorInPath {
		annotations = nil
	}
	ksvc := test.WithServiceReadyOrFail(ctx, test.Service(name, testNamespace, image, annotations))

	WaitForRouteServingText(t, ctx, ksvc.Status.URL.URL(), helloworldText)

	doHelloWorldRequests(ctx, ksvc.Status.URL.URL(), requestCount)

	serviceNamePrefixes := []string{name}
	if activatorInPath {
		serviceNamePrefixes = append(serviceNamePrefixes, "activator-service")
	}
	var err error
	// Verify all the traces of our service also contain spans from the activator.
	// Tracing is asynchronous, retry on failures until timeout is reached.
	if waitErr := wait.PollImmediate(time.Second, 30*time.Second, func() (bool, error) {
		err = verifyServicesArePresentInAllJaegerTraces(ctx, "/", name, serviceNamePrefixes...)
		return err == nil, nil
	}); waitErr != nil {
		t.Fatalf("Error verifying traces: %v: %v", waitErr, err)
	}
}

// Do `count` requests, expect the helloworld response.
func doHelloWorldRequests(ctx *test.Context, url *url.URL, count int) {
	client, err := MakeSpoofingClient(ctx, url)
	if err != nil {
		ctx.T.Fatal("Failed to create client:", err)
	}

	for i := 0; i < count; i++ {
		resp, err := HTTPGetAsStringWithClient(client, url.String())
		if err != nil {
			ctx.T.Errorf("Error GETing %s: %v", url, err)
		}

		if strings.TrimSpace(resp) != helloworldText {
			ctx.T.Errorf("Unexpected response: %s", resp)
		}
	}
}

func HTTPGetAsStringWithClient(client *spoof.SpoofingClient, url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP GET %s returned %d", url, resp.StatusCode)
	}
	return string(resp.Body), err
}

// Queries Jaeger for the services to find the one with the given prefix.
// Lists all traces containing the service name and operation name.
// Verifies all the trace spans cover all the services given in varargs (or at least their prefixes).
// Verifies there is no span in the traces that doesn't match any of the serviceNamePrefixes.
// Returns the number of traces found.
func verifyServicesArePresentInAllJaegerTraces(ctx *test.Context,
	traceOperationName string,
	traceServiceNamePrefix string,
	serviceNamePrefixes ...string) error {
	podList, err := ctx.Clients.Kube.CoreV1().Pods(tracingNamespace).List(
		context.Background(),
		metav1.ListOptions{LabelSelector: "app=" + jaegerName})
	if err != nil {
		return fmt.Errorf("error listing app=%s pods: %w", jaegerName, err)
	}

	if len(podList.Items) != 1 {
		return fmt.Errorf("expecting exactly 1 jaeger pod, got %d", len(podList.Items))
	}

	portForward, err := test.PortForward(podList.Items[0], 16685)
	if err != nil {
		return fmt.Errorf("error creating port-forward: %w", err)
	}
	defer portForward.Close()

	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", portForward.LocalPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error dialing grpc to 127.0.0.1:%d: %w", portForward.LocalPort, err)
	}
	defer conn.Close()

	queryClient := jaegerapi.NewQueryServiceClient(conn)

	// First list all services so that we can find the one matching the prefix we know (as we only know the ksvc name).
	getServicesResponse, err := queryClient.GetServices(context.Background(), &jaegerapi.GetServicesRequest{})
	if err != nil {
		return fmt.Errorf("error getting services: %w", err)
	}

	var serviceName string
	for _, service := range getServicesResponse.Services {
		ctx.T.Logf("service: %s", service)
		if strings.HasPrefix(service, traceServiceNamePrefix) {
			if serviceName != "" {
				return fmt.Errorf("service: %s: found more than one service with %q prefix in Jaeger", service, traceServiceNamePrefix)
			}
			serviceName = service
		}
	}

	if serviceName == "" {
		return fmt.Errorf("didn't find any services with %q prefix in Jaeger", traceServiceNamePrefix)
	}

	// Find all traces matching our service name and the given operation name.
	traceClient, err := queryClient.FindTraces(context.Background(), &jaegerapi.FindTracesRequest{
		Query: &jaegerapi.TraceQueryParameters{
			OperationName: traceOperationName,
			ServiceName:   serviceName,
		},
	})
	if err != nil {
		return fmt.Errorf("error getting FindTraces client: %w", err)
	}

	traces := make(map[string][]jaegermodel.Span)

	// Spans from a single trace can be in different chunks, so gather all traces together first.
	for {
		chunk, err := traceClient.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("error recv traces: %w", err)
		}

		for _, span := range chunk.Spans {
			traces[span.TraceID.String()] = append(traces[span.TraceID.String()], span)
		}
	}

	// We did requestCount requests (+a few during waitForNo503OrFail).
	if len(traces) < requestCount {
		return fmt.Errorf("expected at least %d traces, got %d", requestCount, len(traces))
	}

	for traceID, spans := range traces {
		// All the serviceNamePrefixes should be covered by some traces, so we'll note the matched ones in a boolean array.
		found := make([]bool, len(serviceNamePrefixes))

		ctx.T.Logf("Trace %s:", traceID)

		for _, span := range spans {
			ctx.T.Logf("  %s(%s)", span.Process.ServiceName, span.OperationName)

			matchesAny := false
			for i, serviceNamePrefix := range serviceNamePrefixes {
				if strings.HasPrefix(span.Process.ServiceName, serviceNamePrefix) {
					matchesAny = true
					found[i] = true
				}
			}

			// Verify there is no span in the traces that doesn't match any of the serviceNamePrefixes.
			if !matchesAny {
				return fmt.Errorf("span %s(%s) doesn't match any of the expected prefixes (%v)", span.Process.ServiceName, span.OperationName, serviceNamePrefixes)
			}
		}

		// Verify trace spans cover all services in serviceNamePrefixes.
		for i, serviceNamePrefix := range serviceNamePrefixes {
			if !found[i] {
				return fmt.Errorf("Trace does not contain a span matching serviceName prefix %q", serviceNamePrefix)
			}
		}
	}

	return nil
}
