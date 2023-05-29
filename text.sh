#!/bin/bash
rm -rf tmp
mkdir -p tmp\blocks
mkdir -p tmp\wallets
mkdir -p tmp\ref_list

rm main

go build main.go
