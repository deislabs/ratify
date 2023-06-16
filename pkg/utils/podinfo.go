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

package utils

import "os"

// GetNamespace returns the namespace.
func GetNamespace() string {
	ns, found := os.LookupEnv("RATIFY_NAMESPACE")
	if !found {
		return "gatekeeper-system"
	}
	return ns
}

// GetServiceName returns the service name.
func GetServiceName() string {
	name, found := os.LookupEnv("RATIFY_NAME")
	if !found {
		return "ratify"
	}
	return name
}
