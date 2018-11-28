#!/bin/bash

peer chaincode invoke -o orderer.at2chain.com:7050  --tls true --cafile ${ORDERER_CA} -C $CHANNEL_NAME -n $CC_NAME -c '{"Args":["get","'$1'"]}'

