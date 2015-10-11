#!/bin/bash

#
# Utility script for generating a private/public key pair
#

mkdir keys

openssl genrsa -out keys/$1 2048
openssl rsa -in keys/$1 -pubout > keys/$1.pub 
