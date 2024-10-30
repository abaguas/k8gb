package istiovirtualservice

/*
Copyright 2022 The k8gb Contributors.

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
	"net"
	"testing"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetServers(t *testing.T) {
	var tests = []struct {
		name               string
		virtualServiceFile string
		expectedServers    []*k8gbv1beta1.Server
	}{
		{
			name:               "single host and route",
			virtualServiceFile: "../testdata/istio_virtualservice.yaml",
			expectedServers: []*k8gbv1beta1.Server{
				{
					Host: "istio.cloud.example.com",
					Services: []*k8gbv1beta1.NamespacedName{
						{
							Name:      "istio",
							Namespace: "test-gslb",
						},
					},
				},
			},
		},
		{
			name:               "multiple hosts",
			virtualServiceFile: "./testdata/istio_virtualservice_multiple_hosts.yaml",
			expectedServers: []*k8gbv1beta1.Server{
				{
					Host: "istio1.cloud.example.com",
					Services: []*k8gbv1beta1.NamespacedName{
						{
							Name:      "istio",
							Namespace: "test-gslb",
						},
					},
				},
				{
					Host: "istio2.cloud.example.com",
					Services: []*k8gbv1beta1.NamespacedName{
						{
							Name:      "istio",
							Namespace: "test-gslb",
						},
					},
				},
			},
		},
		{
			name:               "multiple routes",
			virtualServiceFile: "./testdata/istio_virtualservice_multiple_routes.yaml",
			expectedServers: []*k8gbv1beta1.Server{
				{
					Host: "istio.cloud.example.com",
					Services: []*k8gbv1beta1.NamespacedName{
						{
							Name:      "istio1",
							Namespace: "test-gslb",
						},
						{
							Name:      "istio2",
							Namespace: "test-gslb",
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// arrange
			vs := utils.FileToIstioVirtualService(test.virtualServiceFile)
			resolver := ReferenceResolver{
				virtualService: vs,
			}

			// act
			servers, err := resolver.GetServers()
			assert.NoError(t, err)

			// assert
			assert.Equal(t, test.expectedServers, servers)
		})
	}
}

func TestGetGslbExposedIPs(t *testing.T) {
	var tests = []struct {
		name        string
		serviceYaml string
		expectedIPs []net.IP
	}{
		{
			name:        "no exposed IPs",
			serviceYaml: "./testdata/istio_service_no_ips.yaml",
			expectedIPs: []net.IP{},
		},
		{
			name:        "single exposed IP",
			serviceYaml: "../testdata/istio_service.yaml",
			expectedIPs: []net.IP{net.ParseIP("10.0.0.1")},
		},
		{
			name:        "multiple exposed IPs",
			serviceYaml: "./testdata/istio_service_multiple_ips.yaml",
			expectedIPs: []net.IP{net.ParseIP("10.0.0.1"), net.ParseIP("10.0.0.2")},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// arrange
			svc := utils.FileToService(test.serviceYaml)
			resolver := ReferenceResolver{
				lbService: svc,
			}

			// act
			IPs, err := resolver.GetGslbExposedIPs([]utils.DNSServer{})
			assert.NoError(t, err)

			// assert
			assert.Equal(t, test.expectedIPs, IPs)
		})
	}
}
