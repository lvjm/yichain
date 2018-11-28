#!/bin/bash

export CHANNEL_NAME=tracec
export CC_NAME=tokencc
export CC_VERSION=1.0.0

echo "========== Installing chaincode on peer0.core =========="
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/users/Admin@core.at2chain.com/msp
export CORE_PEER_ADDRESS=peer0.core.at2chain.com:7051
export CORE_PEER_LOCALMSPID="CoreMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/peers/peer0.core.at2chain.com/tls/ca.crt
peer chaincode install -p bitbucket.org/at2chain/chaincode-token/chaincode -n $CC_NAME -v $CC_VERSION

echo "========== Installing chaincode on peer1.core =========="
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/users/Admin@core.at2chain.com/msp
export CORE_PEER_ADDRESS=peer1.core.at2chain.com:7051
export CORE_PEER_LOCALMSPID="CoreMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/peers/peer1.core.at2chain.com/tls/ca.crt
peer chaincode install -p bitbucket.org/at2chain/chaincode-token/chaincode -n $CC_NAME -v $CC_VERSION

echo "========== Installing chaincode on peer0.support =========="
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/support.at2chain.com/users/Admin@support.at2chain.com/msp
export CORE_PEER_ADDRESS=peer0.support.at2chain.com:7051
export CORE_PEER_LOCALMSPID="SupportMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/support.at2chain.com/peers/peer0.support.at2chain.com/tls/ca.crt
peer chaincode install -p bitbucket.org/at2chain/chaincode-token/chaincode -n $CC_NAME -v $CC_VERSION

echo "========== Installing chaincode on peer1.support =========="
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/support.at2chain.com/users/Admin@support.at2chain.com/msp
export CORE_PEER_ADDRESS=peer1.support.at2chain.com:7051
export CORE_PEER_LOCALMSPID="SupportMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/support.at2chain.com/peers/peer1.support.at2chain.com/tls/ca.crt
peer chaincode install -p bitbucket.org/at2chain/chaincode-token/chaincode -n $CC_NAME -v $CC_VERSION


echo "========== instantiating chaincode on peer0 =========="
peer chaincode instantiate -n $CC_NAME --tls true --cafile ${ORDERER_CA} -c '{"Args":["init"]}' -v $CC_VERSION -C $CHANNEL_NAME
#peer chaincode instantiate -n $CC_NAME --tls true --cafile ${ORDERER_CA} -c '{"Args":["init"]}' -v $CC_VERSION -C $CHANNEL_NAME -P "OR('Core.member', 'Support.member')"