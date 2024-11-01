package controllers

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
	"context"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
)

func (r *GslbReconciler) finalizeGslb(gslb *k8gbv1beta1.Gslb) (err error) {
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	err = r.DNSProvider.Finalize(gslb, r.Client)
	if err != nil {
		log.Err(err).
			Str("gslb", gslb.Name).
			Msg("Can't finalize GSLB")
		return
	}
	log.Info().
		Str("gslb", gslb.Name).
		Msg("Successfully finalized Gslb")
	return
}

func (r *GslbReconciler) addFinalizer(gslb *k8gbv1beta1.Gslb) error {
	log.Info().
		Str("gslb", gslb.Name).
		Msg("Adding Finalizer for the Gslb")
	gslb.SetFinalizers(append(gslb.GetFinalizers(), gslbFinalizer))

	// Update CR
	err := r.Update(context.TODO(), gslb)
	if err != nil {
		log.Err(err).
			Str("gslb", gslb.Name).
			Msg("Failed to update Gslb with finalizer")
		return err
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}
