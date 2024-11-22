package longrunning

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/openshift-knative/serverless-operator/test"
	"github.com/openshift-knative/serverless-operator/test/servinge2e/servicemesh"
	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgTest "knative.dev/pkg/test"
	"knative.dev/pkg/test/spoof"
	"knative.dev/serving/pkg/apis/autoscaling"
	servingTest "knative.dev/serving/test"
)

const (
	sleepTime = 120000
)

/*
		  Requires the following config in Subscription:
			  spec:
			    config:
			      env:
			        - name: ROUTE_HAPROXY_TIMEOUT
			          value: '180'
	    Requires this setting in ServiceMeshControlPlane:
	      spec:
	        techPreview:
						meshConfig:
							defaultConfig:
								terminationDrainDuration: 190s
*/
func TestTimeoutForLongRunningRequests(t *testing.T) {
	ctx := test.SetupClusterAdmin(t)
	test.CleanupOnInterrupt(t, func() { test.CleanupAll(t, ctx) })
	defer test.CleanupAll(t, ctx)

	service := test.Service("longrunning", test.Namespace, pkgTest.ImagePath(test.AutoscaleImg), map[string]string{
		servicemesh.ServingEnablePassthroughKey: "true",
	}, map[string]string{
		autoscaling.TargetBurstCapacityKey: "-1",
	})
	service = test.WithServiceReadyOrFail(ctx, service)
	serviceURL := service.Status.URL.URL()
	serviceURL.RawQuery = fmt.Sprintf("sleep=%d", sleepTime)

	g := errgroup.Group{}

	for i := 0; i != 300; i++ {
		g.Go(func() error {
			if _, err := pkgTest.WaitForEndpointStateWithTimeout(
				context.Background(),
				ctx.Clients.Kube,
				t.Logf,
				serviceURL,
				spoof.MatchesBody("Slept"),
				"CheckResponse",
				true,
				time.Second*180,
				servingTest.AddRootCAtoTransport(context.Background(), t.Logf, &servingTest.Clients{KubeClient: ctx.Clients.Kube}, true),
			); err != nil {
				return fmt.Errorf("unexpected state for %s :%w", serviceURL, err)
			}
			return nil
		})
	}

	time.Sleep(10 * time.Second)

	// Delete random activator pod
	podList, err := ctx.Clients.Kube.CoreV1().Pods(test.ServingNamespace).List(context.Background(), metav1.ListOptions{LabelSelector: "app=activator"})
	if err != nil {
		t.Fatalf("Error listing activator pods: %v", err)
	}
	if err := ctx.Clients.Kube.CoreV1().Pods(test.ServingNamespace).
		Delete(context.Background(), podList.Items[0].Name, metav1.DeleteOptions{}); err != nil {
		t.Fatalf("Error deleting pod: %v", err)
	}

	time.Sleep(5 * time.Second)

	// Send short requests while activator is terminating.
	serviceURL.RawQuery = fmt.Sprintf("sleep=%d", 0)
	for i := 0; i != 100; i++ {
		g.Go(func() error {
			if _, err := pkgTest.CheckEndpointState(
				context.Background(),
				ctx.Clients.Kube,
				t.Logf,
				serviceURL,
				spoof.MatchesBody("Slept"),
				"CheckResponse",
				true,
				servingTest.AddRootCAtoTransport(context.Background(), t.Logf, &servingTest.Clients{KubeClient: ctx.Clients.Kube}, true),
			); err != nil {
				return fmt.Errorf("unexpected state for %s :%w", serviceURL, err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		t.Error("Something went wrong with the request:", err)
	}
}
