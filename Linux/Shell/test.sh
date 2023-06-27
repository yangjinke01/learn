#!/bin/bash

# [INFO]
# component_name=client
# component_version=v3.3.4-sp4
# component_type=client
# os_arch=aarch64
# system=linux
# os_release=
# cdm_version=v3.3.4-sp4
# config_path=obclient-linux/etc/config.ini
# config_template_path=config.template

function pre_check() {
    echo "checking package ..."
    os_arch=$(awk -F = '/os_arch/ {print $2}' ./manifest)
    if [ $os_arch == $(uname -m) ]; then
        echo "architecture is right, checking os release ..."
        echo "current os release is $(uname -r)"
    else
        echo "package architecture is not right"
        exit 1
    fi

}

pre_check



