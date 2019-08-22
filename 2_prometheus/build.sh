#!bin/bash

# Exit in case of any error
set -e

for filePath in $1/*.yml; do
    fileName="$(expr $filePath \: '/etc/prometheus/\(.*\)')"
    if [ ! "$fileName" == "/prometheus.yml" ]; then
        echo "building file: $fileName" 1>&2
        err="$(promtool check rules $filePath | grep -c SUCCESS)"
    fi
done
