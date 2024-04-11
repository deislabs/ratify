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
	cs "github.com/deislabs/ratify/pkg/customresources/certificatestores"
	"github.com/deislabs/ratify/pkg/customresources/policies"
	rs "github.com/deislabs/ratify/pkg/customresources/referrerstores"
	"github.com/deislabs/ratify/pkg/customresources/verifiers"
)

var (
	VerifierMap = verifiers.NewActiveVerifiers()

	// ActivePolicy is the active policy generated from CRD. There would be exactly
	// one active policy belonging to a namespace at any given time.
	ActivePolicies = policies.NewActivePolicies()

	// a map to track active stores
	StoreMap = rs.NewActiveStores()

	// a map between CertificateStore name to array of x509 certificates
	CertificatesMap = cs.NewActiveCertStores()
)
