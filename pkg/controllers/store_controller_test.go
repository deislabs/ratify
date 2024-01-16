/*
Copyright The Ratify Authors.
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

package controllers

import (
	"context"
	"testing"

	configv1beta1 "github.com/deislabs/ratify/api/v1beta1"
	"github.com/deislabs/ratify/pkg/referrerstore"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestStoreAdd_EmptyParameter(t *testing.T) {
	resetStoreMap()
	var testStoreSpec = configv1beta1.StoreSpec{
		Name: "oras",
	}

	if err := storeAddOrReplace(testStoreSpec, "oras"); err != nil {
		t.Fatalf("storeAddOrReplace() expected no error, actual %v", err)
	}
	if len(StoreMap) != 1 {
		t.Fatalf("Store map expected size 1, actual %v", len(StoreMap))
	}
}

func TestStoreAdd_WithParameters(t *testing.T) {
	resetStoreMap()
	if len(StoreMap) != 0 {
		t.Fatalf("Store map expected size 0, actual %v", len(StoreMap))
	}

	var testStoreSpec = getOrasStoreSpec()

	if err := storeAddOrReplace(testStoreSpec, "testObject"); err != nil {
		t.Fatalf("storeAddOrReplace() expected no error, actual %v", err)
	}
	if len(StoreMap) != 1 {
		t.Fatalf("Store map expected size 1, actual %v", len(StoreMap))
	}
}

func TestWriteStoreStatus(t *testing.T) {
	logger := logrus.WithContext(context.Background())
	testCases := []struct {
		name       string
		isSuccess  bool
		store      *configv1beta1.Store
		errString  string
		reconciler client.StatusClient
	}{
		{
			name:       "success status",
			isSuccess:  true,
			store:      &configv1beta1.Store{},
			reconciler: &mockStatusClient{},
		},
		{
			name:       "error status",
			isSuccess:  false,
			store:      &configv1beta1.Store{},
			errString:  "a long error string that exceeds the max length of 30 characters",
			reconciler: &mockStatusClient{},
		},
		{
			name:      "status update failed",
			isSuccess: true,
			store:     &configv1beta1.Store{},
			reconciler: &mockStatusClient{
				updateFailed: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			writeStoreStatus(context.Background(), tc.reconciler, tc.store, logger, tc.isSuccess, tc.errString)
		})
	}
}
func TestStore_UpdateAndDelete(t *testing.T) {
	resetStoreMap()
	// add a Store

	var resource = "oras"

	var testStoreSpec = getOrasStoreSpec()

	if err := storeAddOrReplace(testStoreSpec, resource); err != nil {
		t.Fatalf("storeAddOrReplace() expected no error, actual %v", err)
	}
	if len(StoreMap) != 1 {
		t.Fatalf("Store map expected size 1, actual %v", len(StoreMap))
	}

	// modify the Store
	var updatedSpec = configv1beta1.StoreSpec{
		Name: "oras",
	}

	if err := storeAddOrReplace(updatedSpec, resource); err != nil {
		t.Fatalf("storeAddOrReplace() expected no error, actual %v", err)
	}

	// validate no Store has been added
	if len(StoreMap) != 1 {
		t.Fatalf("Store map should be 1 after replacement, actual %v", len(StoreMap))
	}

	storeRemove(resource)

	if len(StoreMap) != 0 {
		t.Fatalf("Store map should be 0 after deletion, actual %v", len(StoreMap))
	}
}

func resetStoreMap() {
	StoreMap = map[string]referrerstore.ReferrerStore{}
}

func getOrasStoreSpec() configv1beta1.StoreSpec {
	var parametersString = "{\"authProvider\":{\"name\":\"k8Secrets\",\"secrets\":[{\"secretName\":\"myregistrykey\"}]},\"cosignEnabled\":false,\"useHttp\":false}"
	var storeParameters = []byte(parametersString)

	return configv1beta1.StoreSpec{
		Name: "oras",
		Parameters: runtime.RawExtension{
			Raw: storeParameters,
		},
	}
}
