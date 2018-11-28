#!/bin/bash

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="tracec"}
: ${TIMEOUT:="60"}
COUNTER=1
MAX_RETRY=5

echo "sleep 20 sec..."
sleep 20s
echo "Channel name : "$CHANNEL_NAME
echo "ORDERER_CA: "${ORDERER_CA}

## Create channel
echo "Creating channel..."
peer channel create -o orderer.at2chain.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile ${ORDERER_CA}

## Join Channel
echo "Join peers into the channel..."

echo "Join peer0.core into the channel..."
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/users/Admin@core.at2chain.com/msp
export CORE_PEER_ADDRESS=peer0.core.at2chain.com:7051
export CORE_PEER_LOCALMSPID="CoreMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/peers/peer0.core.at2chain.com/tls/ca.crt

peer channel join -b $CHANNEL_NAME.block

## Update anchor
echo "Update anchor"
peer channel update -o orderer.at2chain.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls true --cafile ${ORDERER_CA}

echo "Join peer1.core into the channel..."
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/users/Admin@core.at2chain.com/msp
export CORE_PEER_ADDRESS=peer1.core.at2chain.com:7051
export CORE_PEER_LOCALMSPID="CoreMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/peers/peer1.core.at2chain.com/tls/ca.crt
peer channel join -b $CHANNEL_NAME.block

echo "Join peer0.support into the channel..."
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/support.at2chain.com/users/Admin@support.at2chain.com/msp
export CORE_PEER_ADDRESS=peer0.support.at2chain.com:7051
export CORE_PEER_LOCALMSPID="SupportMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/support.at2chain.com/peers/peer0.support.at2chain.com/tls/ca.crt
peer channel join -b $CHANNEL_NAME.block

## Update anchor
echo "Update anchor"
peer channel update -o orderer.at2chain.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls true --cafile ${ORDERER_CA}

echo "Join peer1.support into the channel..."
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/support.at2chain.com/users/Admin@support.at2chain.com/msp
export CORE_PEER_ADDRESS=peer1.support.at2chain.com:7051
export CORE_PEER_LOCALMSPID="SupportMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/support.at2chain.com/peers/peer1.support.at2chain.com/tls/ca.crt
peer channel join -b $CHANNEL_NAME.block


## Install chaincode 
echo "Installing chaincode on peer0..."
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/users/Admin@core.at2chain.com/msp
export CORE_PEER_ADDRESS=peer0.core.at2chain.com:7051
export CORE_PEER_LOCALMSPID="CoreMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/core.at2chain.com/peers/peer0.core.at2chain.com/tls/ca.crt
peer chaincode install -p bitbucket.org/at2chain/chaincode-traceability -n at2chain -v 0
peer chaincode install -p bitbucket.org/at2chain/chaincode-digital-asset -n asset -v 0
peer chaincode install -p bitbucket.org/at2chain/chaincode-chinastirling -n chinastirling -v 1.0.0

#Instantiate chaincode
echo "Instantiating chaincode on peer0..."
peer chaincode instantiate -n at2chain --tls true --cafile ${ORDERER_CA} -c '{"Args":["a","100"]}' -v 0 -C $CHANNEL_NAME
peer chaincode instantiate -n asset --tls true --cafile ${ORDERER_CA} -c '{"Args":["init"]}' -v 0 -C $CHANNEL_NAME
peer chaincode instantiate -n chinastirling --tls true --cafile ${ORDERER_CA} -c '{"Args":["init"]}' -v 1.0.0 -C $CHANNEL_NAME

exit 0
