#!/bin/bash

for (( i = 0; i < 1000; ++i ))
do
peer chaincode invoke -o orderer.at2chain.com:7050  --tls true --cafile ${ORDERER_CA}  -C $CHANNEL_NAME -n $CC_NAME -c '{"Args":["update","'$1'","'$2'","'$3'"]}'
done
