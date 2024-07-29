// Package utils contains helper functions and framework for k8gb terratest
package utils

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
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DefaultRetries = 120

// GetIngressIPs returns slice of IP's related to ingress
func GetIngressIPs(t *testing.T, options *k8s.KubectlOptions, ingressName string) []string {
	var ingressIPs []string
	ingress := k8s.GetIngress(t, options, ingressName)
	for _, lb := range ingress.Status.LoadBalancer.Ingress {
		if len(lb.IP) > 0 {
			ingressIPs = append(ingressIPs, lb.IP)
		} else if len(lb.Hostname) > 0 {
			digLbHostnameIPs, _ := Dig(t, "1.1.1.1", 53, lb.Hostname)
			log.Printf("Digging LB hostname %s, got %v", lb.Hostname, digLbHostnameIPs)
			ingressIPs = append(ingressIPs, digLbHostnameIPs...)
		}
	}
	return ingressIPs
}

// Dig gets sorted slice of records related to dnsName
func Dig(t *testing.T, dnsServer string, dnsPort int, dnsName string, additionalArgs ...string) ([]string, error) {
	port := fmt.Sprintf("-p%d", dnsPort)
	dnsServer = fmt.Sprintf("@%s", dnsServer)

	digApp := shell.Command{
		Command: "dig",
		Args:    append([]string{port, dnsServer, dnsName, "+short"}, additionalArgs...),
	}

	digAppOut := shell.RunCommandAndGetOutput(t, digApp)
	digAppSlice := strings.Split(digAppOut, "\n")

	sort.Strings(digAppSlice)

	return digAppSlice, nil
}

// DoWithRetryWaitingForValueE Concept is borrowed from terratest/modules/retry and extended to our use case
func DoWithRetryWaitingForValueE(t *testing.T, actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, action func() ([]string, error), expectedResult []string) ([]string, error) {
	var output []string
	var err error

	for i := 0; i <= maxRetries; i++ {

		output, err = action()
		if err != nil {
			t.Logf("%s returned an error: %s. Sleeping for %s and will try again.", actionDescription, err.Error(), sleepBetweenRetries)
			return output, nil
		}

		if EqualStringSlices(output, expectedResult) {
			t.Logf("%s found match: Expected:(%s). Actual:(%s)", actionDescription, expectedResult, output)
			return output, err
		}

		t.Logf("%s does not match expected result. Expected:(%s). Actual:(%s). Sleeping for %s and will try again.", actionDescription, expectedResult, output, sleepBetweenRetries)
		time.Sleep(sleepBetweenRetries)
	}

	return output, retry.MaxRetriesExceeded{Description: actionDescription, MaxRetries: maxRetries}
}

func CreateGslb(t *testing.T, options *k8s.KubectlOptions, settings TestSettings, kubeResourcePath string) {
	k8sManifestBytes, err := os.ReadFile(kubeResourcePath)
	if err != nil {
		log.Fatal(err)
	}

	zoneReplacer := strings.NewReplacer("cloud.example.com", settings.DNSZone,
		"primaryGeoTag: \"eu\"", fmt.Sprintf("primaryGeoTag: \"%s\"", settings.PrimaryGeoTag),
		"primaryGeoTag: \"us\"", fmt.Sprintf("primaryGeoTag: \"%s\"", settings.SecondaryGeoTag),
		"k8gb.io/primary-geotag: \"eu\"", fmt.Sprintf("k8gb.io/primary-geotag: \"%s\"", settings.PrimaryGeoTag))

	k8sManifestString := zoneReplacer.Replace(string(k8sManifestBytes))

	k8s.KubectlApplyFromString(t, options, k8sManifestString)
}

func InstallPodinfo(t *testing.T, options *k8s.KubectlOptions, settings TestSettings) {
	helmRepoAdd := shell.Command{
		Command: "helm",
		Args:    []string{"repo", "add", "--force-update", "podinfo", "https://stefanprodan.github.io/podinfo"},
	}

	helmRepoUpdate := shell.Command{
		Command: "helm",
		Args:    []string{"repo", "update"},
	}
	shell.RunCommand(t, helmRepoAdd)
	shell.RunCommand(t, helmRepoUpdate)
	helmOptions := helm.Options{
		KubectlOptions: options,
		Version:        "5.2.0",
		SetValues: map[string]string{
			"image.repository": settings.PodinfoImage,
		},
	}
	helm.Install(t, &helmOptions, "podinfo/podinfo", "frontend")

	testAppFilter := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=frontend-podinfo",
	}

	k8s.WaitUntilNumPodsCreated(t, options, testAppFilter, 1, DefaultRetries, 1*time.Second)

	var testAppPods []corev1.Pod

	testAppPods = k8s.ListPods(t, options, testAppFilter)

	for _, pod := range testAppPods {
		k8s.WaitUntilPodAvailable(t, options, pod.Name, DefaultRetries, 1*time.Second)
	}

	k8s.WaitUntilServiceAvailable(t, options, "frontend-podinfo", DefaultRetries, 1*time.Second)

}

func CreateGslbWithHealthyApp(t *testing.T, options *k8s.KubectlOptions, settings TestSettings, kubeResourcePath, gslbName, hostName string) {

	CreateGslb(t, options, settings, kubeResourcePath)

	k8s.WaitUntilIngressAvailable(t, options, gslbName, DefaultRetries, 1*time.Second)
	ingress := k8s.GetIngress(t, options, gslbName)
	require.Equal(t, ingress.Name, gslbName)

	InstallPodinfo(t, options, settings)

	serviceHealthStatus := fmt.Sprintf("%s:Healthy", hostName)
	AssertGslbStatus(t, options, gslbName, serviceHealthStatus)
}

func AssertGslbStatus(t *testing.T, options *k8s.KubectlOptions, gslbName, serviceStatus string) {

	t.Helper()

	actualHealthStatus := func() ([]string, error) {
		k8gbServiceHealth, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "gslb", gslbName, "-o",
			"custom-columns=SERVICESTATUS:.status.serviceHealth", "--no-headers")
		if err != nil {
			t.Logf("Failed to get k8gb status with kubectl (%s)", err)
		}
		return []string{k8gbServiceHealth}, nil
	}
	expectedHealthStatus := []string{fmt.Sprintf("map[%s]", serviceStatus)}
	_, err := DoWithRetryWaitingForValueE(
		t,
		"Wait for expected ServiceHealth status...",
		DefaultRetries,
		1*time.Second,
		actualHealthStatus,
		expectedHealthStatus)
	require.NoError(t, err)
}

func AssertGslbSpec(t *testing.T, options *k8s.KubectlOptions, gslbName, specPath, expectedValue string) {
	t.Helper()
	actualValue, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "gslb", gslbName, "-o", fmt.Sprintf("custom-columns=SERVICESTATUS:%s", specPath), "--no-headers")
	require.NoError(t, err)
	assert.Equal(t, expectedValue, actualValue)
}

func AssertDNSEndpointLabel(t *testing.T, options *k8s.KubectlOptions, label string) {
	t.Helper()
	k8s.RunKubectl(t, options, "get", "dnsendpoint", "-l", label)
}

func AssertGslbDeleted(t *testing.T, options *k8s.KubectlOptions, gslbName string) {
	t.Helper()
	deletionExpected := []string{fmt.Sprintf("Error from server (NotFound): gslbs.k8gb.absa.oss \"%s\" not found", gslbName)}
	deletionActual, err := DoWithRetryWaitingForValueE(
		t,
		"Waiting for Gslb CR to be deleted...",
		300,
		1*time.Second,
		func() ([]string, error) {
			out, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "gslb", gslbName)
			return []string{out}, err
		},
		deletionExpected)
	require.NoError(t, err)

	assert.Equal(t, deletionExpected, deletionActual)
}

func WaitForLocalGSLB(t *testing.T, dnsServer string, dnsPort int, settings TestSettings, host string, expectedResult []string) (output []string, err error) {
	var additionalArgs []string
	if !settings.DigUsingUDP {
		additionalArgs = append(additionalArgs, "+tcp")
	}
	return DoWithRetryWaitingForValueE(
		t,
		"Wait for failover to happen and coredns to pickup new values...",
		300,
		time.Second*1,
		func() ([]string, error) { return Dig(t, dnsServer, dnsPort, host+settings.DNSZone, additionalArgs...) },
		expectedResult)
}

// EqualStringSlices tells whether slice a and b contain the same elements in unsorted slices.
// A nil argument is equivalent to an empty slice.
func EqualStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	x := make([]string, len(a))
	y := make([]string, len(b))
	copy(x, a)
	copy(y, b)
	sort.Strings(x)
	sort.Strings(y)
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}

// RunBusyBoxCommand the command argument is executed inside the busybox pod. It can be for example an HTTP request etc.
func RunBusyBoxCommand(t *testing.T, options *k8s.KubectlOptions, dns string, command []string) (out string, err error) {
	dnsOverride := fmt.Sprintf("{\"spec\":{\"dnsConfig\":{\"nameservers\":[\"%s\"]},\"dnsPolicy\": \"None\"}}", dns)
	args := []string{}
	kubectlCtx := []string{"--context", options.ContextName, "-n", options.Namespace}
	containerArgs := []string{"run", "-i", "--rm", "busybox", "--restart", "Never", "--image", "busybox"}
	containerDNS := []string{"--overrides", dnsOverride}
	appArgs := append([]string{"--"}, command...)
	args = append(kubectlCtx, containerArgs...)
	args = append(args, containerDNS...)
	args = append(args, appArgs...)
	cmd := shell.Command{
		Command: "kubectl",
		Args:    args,
		Env:     options.Env,
	}
	return shell.RunCommandAndGetOutputE(t, cmd)
}

// EqualAnnotations compares annotations of two maps
// TODO: redundant to k8gb utils internal function.
func EqualAnnotations(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if b[k] != a[k] {
			return false
		}
	}
	return true
}
