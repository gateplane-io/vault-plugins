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
	"github.com/gateplane-io/vault-plugins/internal/mock"
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

func Backend(c *logical.BackendConfig) *mock.Backend {
	var baseBackend base.BaseBackend
	bFinal := &mock.Backend{BaseBackend: &baseBackend}

	baseBackend.Backend = &framework.Backend{
		BackendType:    logical.TypeCredential,
		Help:           "Vault/OpenBao Plugin for approval-based access to policies",
		RunningVersion: Version,
		PathsSpecial: &logical.Paths{
			// Keeping Claim unauthenticated to match the case
			// of Policy Gate - map to AccessRequest through GrantCode
			Unauthenticated: []string{"claim"},
		},
		Paths: []*framework.Path{
			// Provided by Base package
			base.PathConfig(&baseBackend),

			base.PathRequest(&baseBackend),
			base.PathApprove(&baseBackend),

			// Custom Claim endpoint
			mock.PathClaim(bFinal),
		},
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
