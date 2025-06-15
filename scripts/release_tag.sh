#!/bin/sh

tag=`date +"%-Y.%-m.%-d"`-release
echo "Tagging with '$tag'"

git tag --message "Release tag on '$(date)'" "$tag"
