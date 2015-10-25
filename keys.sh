#!/bin/bash

#
# Utility script for generating a private/public key pair
#

if [[ -z "$1" ]]; then
	echo "usage: ./keys.sh [key_basename]"
	exit 1
fi

mkdir keys

openssl genrsa -out keys/$1 2048
openssl rsa -in keys/$1 -pubout > keys/$1.pub 
