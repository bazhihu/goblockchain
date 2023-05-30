#!/bin/bash
rm -rf tmp
mkdir -p tmp/blocks
mkdir -p tmp/wallets
mkdir -p tmp/ref_list

rm main

go build main.go

./main createblockchain -address LeoCao
./main blockchaininfo
./main balance -address LeoCao
./main send -from LeoCao -to Krad -amount 100
./main balance -address Krad
./main mine
./main blockchaininfo
./main balance -address LeoCao
./main balance -address Krad
./main send -from LeoCao -to Exia -amount 100
./main send -from Krad -to Exia -amount 30
./main mine
./main blockchaininfo
./main balance -address LeoCao
./main balance -address Krad
./main balance -address Exia