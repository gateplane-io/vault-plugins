#!/bin/bash

# Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
# SPDX-License-Identifier: Elastic-2.0
#
# Licensed under the Elastic License 2.0.
# You may obtain a copy of the license at:
# https://www.elastic.co/licensing/elastic-license
#
# Use, modification, and redistribution permitted under the terms of the license,
# except for providing this software as a commercial service or product.

# Fetches the 2 latest *-release Git Tags and updates
# them in GoReleaser as env vars:
# https://goreleaser.com/cookbooks/set-a-custom-git-tag/

release_tags=$(git for-each-ref refs/tags/*-release \
	--sort=-refname \
	--format='%(refname)' \
	--count=2 \
| cut -d '/' -f 3 \
| sort -r \
| tr '\n' ' ')


read GORELEASER_CURRENT_TAG GORELEASER_PREVIOUS_TAG < <(echo $release_tags)

echo "
export GORELEASER_CURRENT_TAG=$GORELEASER_CURRENT_TAG
export GORELEASER_PREVIOUS_TAG=$GORELEASER_PREVIOUS_TAG
" | tee /dev/stderr
