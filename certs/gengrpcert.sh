#!/bin/bash

#clean any previous certs
find . -type f -name server.\* -exec rm -f {} \;
find . -type f -name ca.\* -exec rm -f {} \;
# generate new
openssl genrsa -out ca.key 4096
openssl req -new -x509 -key ca.key -sha256 -subj "//C=PO\ST=NB\O=Test" -days 365 -out ca.crt
openssl genrsa -out server.key 4096
openssl req -new -key server.key -out server.csr -config certificate.conf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256 -extfile certificate.conf -extensions req_ext
openssl pkcs8 -topk8 -nocrypt -in server.key -out server.pem