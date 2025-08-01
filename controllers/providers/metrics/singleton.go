package metrics

/*
Copyright 2021-2025 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"sync"

	"github.com/k8gb-io/k8gb/controllers/resolver"
)

const (
	// DefaultMetricsNamespace provides the default namespace used, when PrometheusMetrics was not initialised
	DefaultMetricsNamespace = "k8gb_default"
)

var (
	o       sync.Once
	metrics = nop()
)

// Metrics public static metrics, providing instance of initialised metrics
func Metrics() *PrometheusMetrics {
	return &metrics
}

// Init always initialise PrometheusMetrics. The initialisation happens only once
func Init(c *resolver.Config) {
	o.Do(func() {
		metrics = *newPrometheusMetrics(*c)
	})
}

// nop provides default PrometheusMetrics in case metrics was not initialised yet.
func nop() PrometheusMetrics {
	return *newPrometheusMetrics(resolver.Config{K8gbNamespace: DefaultMetricsNamespace})
}
