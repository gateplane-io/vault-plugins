#!/bin/bash

plugins="$(ls cmd/ | tr '\n' ' ')"

echo "# Plugin codebases: [$plugins]"
base_version=$(git describe --tags --match "base/*" --abbrev=0 | cut -d/ -f2)
echo "# Base - base/${base_version}"
for PLUGIN in $plugins; do
	plugin_version=$(git describe --tags --match "$PLUGIN/*" --abbrev=0  2>/dev/null| cut -d/ -f2)

	if [ ! -n "$plugin_version" ]; then
		echo "# Plugin '${PLUGIN}' does not have a tag of ${PLUGIN}/x.x.x"
		continue
	fi
	tag="v$plugin_version-base.$base_version"
	echo "# ${PLUGIN} - ${PLUGIN}/${plugin_version}"
	echo "export VERSION_${PLUGIN//-/_}=$tag" | tee /dev/stderr

done;
