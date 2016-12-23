#!/bin/bash

if [ "$1" == "onlp" ]; then
	cd pluginManager/onlp/libs
	ln -sf libonlp.so.1 libonlp.so
	cd -
fi
