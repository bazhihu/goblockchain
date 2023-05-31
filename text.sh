#!/bin/bash
rm -rf tmp
mkdir -p tmp/blocks
mkdir -p tmp/wallets
mkdir -p tmp/ref_list

rm main

go1.18.7 build main.go

./main createwallet 
./main walletslist
./main createwallet -refname LeoCao
./main walletinfo -refname LeoCao
./main createwallet -refname Krad
./main createwallet -refname Exia
./main createwallet 
./main walletslist
./main createblockchain -refname LeoCao
./main blockchaininfo
./main balance -refname LeoCao
./main sendbyrefname -from LeoCao -to Krad -amount 100
./main balance -refname Krad
./main mine
./main blockchaininfo
./main balance -refname LeoCao
./main balance -refname Krad
./main sendbyrefname -from LeoCao -to Exia -amount 100
./main sendbyrefname -from Krad -to Exia -amount 30
./main mine
./main blockchaininfo
./main balance -refname LeoCao
./main balance -refname Krad
./main balance -refname Exia
./main sendbyrefname -from Exia -to LeoCao -amount 90
./main sendbyrefname -from Exia -to Krad -amount 90
./main mine
./main blockchaininfo
./main balance -refname LeoCao
./main balance -refname Krad
./main balance -refname Exia