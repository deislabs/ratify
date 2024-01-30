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

package factory

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/deislabs/ratify/pkg/featureflag"
	"github.com/deislabs/ratify/pkg/referrerstore"
	"github.com/deislabs/ratify/pkg/referrerstore/config"
	"github.com/deislabs/ratify/pkg/referrerstore/mocks"
	"github.com/deislabs/ratify/pkg/referrerstore/plugin"
)

type TestStoreFactory struct{}

func (f *TestStoreFactory) Create(_ string, _ config.StorePluginConfig) (referrerstore.ReferrerStore, error) {
	return &mocks.TestStore{}, nil
}

func TestCreateStoresFromConfig_BuiltInStores_ReturnsExpected(t *testing.T) {
	builtInStores = map[string]StoreFactory{
		"testStore": &TestStoreFactory{},
	}

	storeConfig := map[string]interface{}{
		"name": "testStore",
	}
	storesConfig := config.StoresConfig{
		Stores: []config.StorePluginConfig{storeConfig},
	}

	stores, err := CreateStoresFromConfig(storesConfig, getReferrerstorePluginsDir())

	if err != nil {
		t.Fatalf("create stores failed with err %v", err)
	}

	if len(stores) != 1 {
		t.Fatalf("expected to have %d stores, actual count %d", 1, len(stores))
	}

	if stores[0].Name() != "testStore" {
		t.Fatalf("expected to create test store")
	}

	if _, ok := stores[0].(*plugin.StorePlugin); ok {
		t.Fatalf("type assertion failed expected a built in store")
	}
}

func TestCreateStoresFromConfig_PluginStores_ReturnsExpected(t *testing.T) {
	storeConfig := map[string]interface{}{
		"name": "sample",
	}
	storesConfig := config.StoresConfig{
		Stores: []config.StorePluginConfig{storeConfig},
	}

	stores, err := CreateStoresFromConfig(storesConfig, getReferrerstorePluginsDir())

	if err != nil {
		t.Fatalf("create stores failed with err %v", err)
	}

	if len(stores) != 1 {
		t.Fatalf("expected to have %d stores, actual count %d", 1, len(stores))
	}

	if stores[0].Name() != "sample" {
		t.Fatalf("expected to create plugin store")
	}

	if _, ok := stores[0].(*plugin.StorePlugin); !ok {
		t.Fatalf("type assertion failed expected a plugin store")
	}
}

func TestCreateStoresFromConfig_DynamicPluginStores_ReturnsExpected(t *testing.T) {
	os.Setenv("RATIFY_EXPERIMENTAL_DYNAMIC_PLUGINS", "1")
	featureflag.InitFeatureFlagsFromEnv()

	testCases := []struct {
		name     string
		artifact string
	}{
		{
			name:     "image specified by tag",
			artifact: "wabbitnetworks.azurecr.io/test/sample-store-plugin:v1",
		},
		{
			name:     "image specified by digest",
			artifact: "wabbitnetworks.azurecr.io/test/sample-store-plugin@sha256:96ba9f9636cde32df87d62dcad4e430d055e708b9f173475c5d7468b732d6566",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeConfig := map[string]interface{}{
				"name": "plugin-store",
				"source": map[string]interface{}{
					"artifact": tc.artifact,
				},
			}

			storesConfig := config.StoresConfig{
				Stores: []config.StorePluginConfig{storeConfig},
			}
			stores, err := CreateStoresFromConfig(storesConfig, getReferrerstorePluginsDir())

			if err != nil {
				t.Fatalf("create stores failed with err %v", err)
			}

			if len(stores) != 1 {
				t.Fatalf("expected to have %d stores, actual count %d", 1, len(stores))
			}

			if stores[0].Name() != "plugin-store" {
				t.Fatalf("expected to create plugin store")
			}

			if _, ok := stores[0].(*plugin.StorePlugin); !ok {
				t.Fatalf("type assertion failed expected a plugin store")
			}

			pluginPath := path.Join(stores[0].GetConfig().PluginBinDirs[0], stores[0].Name())
			if _, err := os.Stat(pluginPath); errors.Is(err, os.ErrNotExist) {
				t.Fatalf("downloaded plugin not found in path")
			}
		})
	}
}

func getReferrerstorePluginsDir() string {
	workingDir, _ := os.Getwd()
	pluginDir := filepath.Clean(filepath.Join(workingDir, "../../../", "./bin/plugins/referrerstore/"))
	return pluginDir
}
