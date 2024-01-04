package collector

import (
	"context"
	"fmt"
	"github.com/DataDog/KubeHound/pkg/config"
	"github.com/DataDog/KubeHound/pkg/telemetry/log"
	"github.com/DataDog/KubeHound/pkg/telemetry/metric"
	"github.com/DataDog/KubeHound/pkg/telemetry/span"
	"github.com/DataDog/KubeHound/pkg/telemetry/statsd"
	"github.com/DataDog/KubeHound/pkg/telemetry/tag"
	"go.uber.org/ratelimit"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/pager"
	ctrl "sigs.k8s.io/controller-runtime"

	routev1 "github.com/openshift/api/route/v1"
	routev1Clientset "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
)

type openShiftAPICollector struct {
	*k8sAPICollector
	routeClientset routev1Clientset.RouteV1Client
}

const (
	OpenShiftAPICollectorName = "openshift-api-collector"
)

func checkOpenShiftAPICollectorConfig(collectorType string) error {
	if collectorType != config.CollectorTypeOpenShiftAPI {
		return fmt.Errorf("invalid collector type in config: %s", collectorType)
	}

	return nil
}

// NewOpenShiftAPICollector creates a new instance of the OpenShit live API collector from the provided application config.
func NewOpenShiftAPICollector(ctx context.Context, cfg *config.KubehoundConfig) (OpenShiftCollectorClient, error) {
	tags := tag.BaseTags
	tags = append(tags, tag.Collector(OpenShiftAPICollectorName))
	l := log.Trace(ctx, log.WithComponent(OpenShiftAPICollectorName))

	err := checkOpenShiftAPICollectorConfig(cfg.Collector.Type)
	if err != nil {
		return nil, err
	}

	kubeConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("building kubernetes config: %w", err)
	}

	kubeConfig.UserAgent = CollectorUserAgent

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("getting kubernetes config: %w", err)
	}

	return &openShiftAPICollector{
		k8sAPICollector: &k8sAPICollector{
			cfg:       cfg.Collector.Live,
			clientset: clientset,
			log:       l,
			rl:        ratelimit.New(cfg.Collector.Live.RateLimitPerSecond), // per second
			tags:      tags,
		},
		routeClientset: *routev1Clientset.NewForConfigOrDie(kubeConfig),
	}, nil
}

func (c *openShiftAPICollector) Name() string {
	return OpenShiftAPICollectorName
}

// streamRoutesNamespace streams the endpoint slice objects corresponding to a cluster namespace.
func (c *openShiftAPICollector) streamRoutesNamespace(ctx context.Context, namespace string, ingestor RouteIngestor) error {
	err := c.checkNamespaceExists(ctx, namespace)
	if err != nil {
		return err
	}

	opts := tunedListOptions()
	pager := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
		entries, err := c.routeClientset.Routes(namespace).List(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("getting K8s endpoint slices for namespace %s: %w", namespace, err)
		}

		return entries, err
	}))

	c.setPagerConfig(pager)

	return pager.EachListItem(ctx, opts, func(obj runtime.Object) error {
		_ = statsd.Incr(metric.CollectorCount, append(c.tags, tag.Entity(tag.EntityRoutes)), 1)
		c.rl.Take()
		item, ok := obj.(*routev1.Route)
		if !ok {
			return fmt.Errorf("endpoint stream type conversion error: %T", obj)
		}

		err := ingestor.IngestRoute(ctx, item)
		if err != nil {
			return fmt.Errorf("processing K8s endpoint slice %s for namespace %s: %w", item.Name, namespace, err)
		}

		return nil
	})
}

func (c *openShiftAPICollector) StreamRoutes(ctx context.Context, ingestor RouteIngestor) error {
	span, ctx := tracer.StartSpanFromContext(ctx, span.CollectorStream, tracer.Measured())
	span.SetTag(tag.EntityTag, tag.EntityRoutes)
	defer span.Finish()

	// passing an empty namespace will collect all namespaces
	err := c.streamRoutesNamespace(ctx, "", ingestor)
	if err != nil {
		return err
	}

	return ingestor.Complete(ctx)
}
