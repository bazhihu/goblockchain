#!/bin/bash
rm -rf tmp
mkdir -r tmp\blocks
mkdir -r tmp\wallets
mkdir -r tmp\ref_list

rm main

go build main.go
