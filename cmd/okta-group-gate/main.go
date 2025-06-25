// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin"

	"github.com/gateplane-io/vault-plugins/internal/base"
	okta_gate "github.com/gateplane-io/vault-plugins/internal/okta-group-gate"
)

// set at buildtime with "-ldflags -X main.Version=..."
var Version = "0.0.0"

func BackendFactory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(c)
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend(c *logical.BackendConfig) *okta_gate.Backend {
	var baseBackend base.BaseBackend
	bFinal := &okta_gate.Backend{BaseBackend: &baseBackend}

	baseBackend.Backend = &framework.Backend{
		BackendType:    logical.TypeCredential,
		Help:           "Vault/OpenBao Plugin for approval-based access to Okta Groups",
		RunningVersion: Version,
		Paths: []*framework.Path{
			// Provided by Base package
			base.PathRequest(&baseBackend),
			base.PathApprove(&baseBackend),

			// Custom Claim endpoint
			okta_gate.PathClaim(bFinal),
			okta_gate.PathConfig(bFinal),
			okta_gate.PathOktaApiConfig(bFinal),
			okta_gate.PathOktaApiTest(bFinal),
		},
		// Used to Ensure Okta Users are removed from Groups
		PeriodicFunc: bFinal.CleanExpiredMemberships,
		// Clean: bFinal.CleanMemberships,
	}

	bFinal.Logger().Debug("Plugin initialized")
	return bFinal
}

func main() {

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	if err := plugin.ServeMultiplex(&plugin.ServeOpts{
		BackendFactoryFunc: BackendFactory,
		TLSProviderFunc:    tlsProviderFunc,
	}); err != nil {
		log.Fatal(err)
	}
}
