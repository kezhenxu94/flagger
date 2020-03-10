/*
Copyright 2018 The Flagger Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"
	"time"

	istiov1alpha3 "github.com/weaveworks/flagger/pkg/apis/istio/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	CanaryKind              = "Canary"
	ProgressDeadlineSeconds = 600
	AnalysisInterval        = 60 * time.Second
	MetricInterval          = "1m"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Canary is the configuration for a canary release,
// which automatically manages the bootstrap, analysis, traffic shifting,
// promotion or rollback of an app revision
type Canary struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CanarySpec   `json:"spec"`
	Status CanaryStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CanaryList is a list of Canary resources
type CanaryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Canary `json:"items"`
}

// CanarySpec is the specification of the desired behavior of the Canary
type CanarySpec struct {
	// Provider overwrites the -mesh-provider flag for this particular canary
	// +optional
	Provider string `json:"provider,omitempty"`

	// MetricsServer overwrites the -metrics-server flag for this particular canary
	// +optional
	MetricsServer string `json:"metricsServer,omitempty"`

	// TargetRef references a target resource
	TargetRef CrossNamespaceObjectReference `json:"targetRef"`

	// AutoscalerRef references an autoscaling resource
	// +optional
	AutoscalerRef *CrossNamespaceObjectReference `json:"autoscalerRef,omitempty"`

	// Reference to NGINX ingress resource
	// +optional
	IngressRef *CrossNamespaceObjectReference `json:"ingressRef,omitempty"`

	// Service defines how ClusterIP services, service mesh or ingress routing objects are generated
	Service CanaryService `json:"service"`

	// Analysis defines the validation process of a release
	Analysis *CanaryAnalysis `json:"analysis,omitempty"`

	// Deprecated: replaced by Analysis
	CanaryAnalysis *CanaryAnalysis `json:"canaryAnalysis,omitempty"`

	// ProgressDeadlineSeconds represents the maximum time in seconds for a
	// canary deployment to make progress before it is considered to be failed
	// +optional
	ProgressDeadlineSeconds *int32 `json:"progressDeadlineSeconds,omitempty"`

	// SkipAnalysis promotes the canary without analysing it
	// +optional
	SkipAnalysis bool `json:"skipAnalysis,omitempty"`
}

// CanaryService defines how ClusterIP services, service mesh or ingress routing objects are generated
type CanaryService struct {
	// Name of the Kubernetes service generated by Flagger
	// Defaults to CanarySpec.TargetRef.Name
	// +optional
	Name string `json:"name,omitempty"`

	// Port of the generated Kubernetes service
	Port int32 `json:"port"`

	// Port name of the generated Kubernetes service
	// Defaults to http
	// +optional
	PortName string `json:"portName,omitempty"`

	// Target port number or name of the generated Kubernetes service
	// Defaults to CanaryService.Port
	// +optional
	TargetPort intstr.IntOrString `json:"targetPort,omitempty"`

	// PortDiscovery adds all container ports to the generated Kubernetes service
	PortDiscovery bool `json:"portDiscovery"`

	// Timeout of the HTTP or gRPC request
	// +optional
	Timeout string `json:"timeout,omitempty"`

	// Gateways attached to the generated Istio virtual service
	// Defaults to the internal mesh gateway
	// +optional
	Gateways []string `json:"gateways,omitempty"`

	// Hosts attached to the generated Istio virtual service
	// Defaults to the service name
	// +optional
	Hosts []string `json:"hosts,omitempty"`

	// TrafficPolicy attached to the generated Istio destination rules
	// +optional
	TrafficPolicy *istiov1alpha3.TrafficPolicy `json:"trafficPolicy,omitempty"`

	// URI match conditions for the generated service
	// +optional
	Match []istiov1alpha3.HTTPMatchRequest `json:"match,omitempty"`

	// Rewrite HTTP URIs for the generated service
	// +optional
	Rewrite *istiov1alpha3.HTTPRewrite `json:"rewrite,omitempty"`

	// Retries policy for the generated virtual service
	// +optional
	Retries *istiov1alpha3.HTTPRetry `json:"retries,omitempty"`

	// Headers operations for the generated Istio virtual service
	// +optional
	Headers *istiov1alpha3.Headers `json:"headers,omitempty"`

	// Cross-Origin Resource Sharing policy for the generated Istio virtual service
	// +optional
	CorsPolicy *istiov1alpha3.CorsPolicy `json:"corsPolicy,omitempty"`

	// Mesh name of the generated App Mesh virtual nodes and virtual service
	// +optional
	MeshName string `json:"meshName,omitempty"`

	// Backends of the generated App Mesh virtual nodes
	// +optional
	Backends []string `json:"backends,omitempty"`
}

// CanaryAnalysis is used to describe how the analysis should be done
type CanaryAnalysis struct {
	// Schedule interval for this canary analysis
	Interval string `json:"interval"`

	// Number of checks to run for A/B Testing and Blue/Green
	// +optional
	Iterations int `json:"iterations,omitempty"`

	//Enable traffic mirroring for Blue/Green
	// +optional
	Mirror bool `json:"mirror,omitempty"`

	// Max traffic percentage routed to canary
	// +optional
	MaxWeight int `json:"maxWeight,omitempty"`

	// Incremental traffic percentage step
	// +optional
	StepWeight int `json:"stepWeight,omitempty"`

	// Max number of failed checks before the canary is terminated
	Threshold int `json:"threshold"`

	// Alert list for this canary analysis
	Alerts []CanaryAlert `json:"alerts,omitempty"`

	// Metric check list for this canary analysis
	// +optional
	Metrics []CanaryMetric `json:"metrics,omitempty"`

	// Webhook list for this canary  analysis
	// +optional
	Webhooks []CanaryWebhook `json:"webhooks,omitempty"`

	// A/B testing HTTP header match conditions
	// +optional
	Match []istiov1alpha3.HTTPMatchRequest `json:"match,omitempty"`
}

// CanaryMetric holds the reference to metrics used for canary analysis
type CanaryMetric struct {
	// Name of the metric
	Name string `json:"name"`

	// Interval represents the windows size
	Interval string `json:"interval,omitempty"`

	// Deprecated: Max value accepted for this metric (replaced by ThresholdRange)
	Threshold float64 `json:"threshold"`

	// Range value accepted for this metric
	// +optional
	ThresholdRange *CanaryThresholdRange `json:"thresholdRange,omitempty"`

	// Deprecated: Prometheus query for this metric (replaced by TemplateRef)
	// +optional
	Query string `json:"query,omitempty"`

	// TemplateRef references a metric template object
	// +optional
	TemplateRef *CrossNamespaceObjectReference `json:"templateRef,omitempty"`
}

// CanaryThresholdRange defines the range used for metrics validation
type CanaryThresholdRange struct {
	// Minimum value
	// +optional
	Min *float64 `json:"min,omitempty"`

	// Maximum value
	// +optional
	Max *float64 `json:"max,omitempty"`
}

// AlertSeverity defines alert filtering based on severity levels
type AlertSeverity string

const (
	SeverityInfo  AlertSeverity = "info"
	SeverityWarn  AlertSeverity = "warn"
	SeverityError AlertSeverity = "error"
)

// CanaryAlert defines an alert for this canary
type CanaryAlert struct {
	// Name of the alert
	Name string `json:"name"`

	// Severity level: info, warn, error (default info)
	Severity AlertSeverity `json:"severity,omitempty"`

	// Alert provider reference
	ProviderRef CrossNamespaceObjectReference `json:"providerRef"`
}

// HookType can be pre, post or during rollout
type HookType string

const (
	// RolloutHook execute webhook during the canary analysis
	RolloutHook HookType = "rollout"
	// PreRolloutHook execute webhook before routing traffic to canary
	PreRolloutHook HookType = "pre-rollout"
	// PreRolloutHook execute webhook after the canary analysis
	PostRolloutHook HookType = "post-rollout"
	// ConfirmRolloutHook halt canary analysis until webhook returns HTTP 200
	ConfirmRolloutHook HookType = "confirm-rollout"
	// ConfirmPromotionHook halt canary promotion until webhook returns HTTP 200
	ConfirmPromotionHook HookType = "confirm-promotion"
	// EventHook dispatches Flagger events to the specified endpoint
	EventHook HookType = "event"
	// RollbackHook rollback canary analysis if webhook returns HTTP 200
	RollbackHook HookType = "rollback"
)

// CanaryWebhook holds the reference to external checks used for canary analysis
type CanaryWebhook struct {
	// Type of this webhook
	Type HookType `json:"type"`

	// Name of this webhook
	Name string `json:"name"`

	// URL address of this webhook
	URL string `json:"url"`

	// Request timeout for this webhook
	Timeout string `json:"timeout"`

	// Metadata (key-value pairs) for this webhook
	// +optional
	Metadata *map[string]string `json:"metadata,omitempty"`
}

// CanaryWebhookPayload holds the deployment info and metadata sent to webhooks
type CanaryWebhookPayload struct {
	// Name of the canary
	Name string `json:"name"`

	// Namespace of the canary
	Namespace string `json:"namespace"`

	// Phase of the canary analysis
	Phase CanaryPhase `json:"phase"`

	// Metadata (key-value pairs) for this webhook
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CrossNamespaceObjectReference contains enough information to let you locate the
// typed referenced object at cluster level
type CrossNamespaceObjectReference struct {
	// API version of the referent
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referent
	// +optional
	Kind string `json:"kind,omitempty"`

	// Name of the referent
	Name string `json:"name"`

	// Namespace of the referent
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// GetServiceNames returns the apex, primary and canary Kubernetes service names
func (c *Canary) GetServiceNames() (apexName, primaryName, canaryName string) {
	apexName = c.Spec.TargetRef.Name
	if c.Spec.Service.Name != "" {
		apexName = c.Spec.Service.Name
	}
	primaryName = fmt.Sprintf("%s-primary", apexName)
	canaryName = fmt.Sprintf("%s-canary", apexName)
	return
}

// GetProgressDeadlineSeconds returns the progress deadline (default 600s)
func (c *Canary) GetProgressDeadlineSeconds() int {
	if c.Spec.ProgressDeadlineSeconds != nil {
		return int(*c.Spec.ProgressDeadlineSeconds)
	}

	return ProgressDeadlineSeconds
}

// GetAnalysis returns the analysis v1beta1 or v1alpha3
// to be removed along with spec.canaryAnalysis in v1
func (c *Canary) GetAnalysis() *CanaryAnalysis {
	if c.Spec.Analysis != nil {
		return c.Spec.Analysis
	}
	return c.Spec.CanaryAnalysis
}

// GetAnalysisInterval returns the canary analysis interval (default 60s)
func (c *Canary) GetAnalysisInterval() time.Duration {
	if c.GetAnalysis().Interval == "" {
		return AnalysisInterval
	}

	interval, err := time.ParseDuration(c.GetAnalysis().Interval)
	if err != nil {
		return AnalysisInterval
	}

	if interval < 10*time.Second {
		return time.Second * 10
	}

	return interval
}

// GetAnalysisThreshold returns the canary threshold (default 1)
func (c *Canary) GetAnalysisThreshold() int {
	if c.GetAnalysis().Threshold > 0 {
		return c.GetAnalysis().Threshold
	}
	return 1
}

// GetMetricInterval returns the metric interval default value (1m)
func (c *Canary) GetMetricInterval() string {
	return MetricInterval
}

// SkipAnalysis returns true if the analysis is nil
// or if spec.SkipAnalysis is true
func (c *Canary) SkipAnalysis() bool {
	if c.Spec.Analysis == nil && c.Spec.CanaryAnalysis == nil {
		return true
	}
	return c.Spec.SkipAnalysis
}