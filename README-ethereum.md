# at2chain

## Add new node 
```
// install golang
curl -OL https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz

// install go-ethereum
git clone https://github.com/ethereum/go-ethereum.git $HOME/go/src/github.com
cd $HOME/go/src/github.com/go-ethereum
make all

// ubuntu
sudo vi /etc/profile
// add line "export PATH=$PATH:/root/go-ethereum/build/bin" 

source /etc/profile

// mac
sudo vi /etc/paths
// add line "export PATH=$PATH:/root/go-ethereum/build/bin" 
reboot


// download at2chain json
git clone https://bitbucket.org/AI-CORE/at2chain.git $HOME/go/src/bitbucket.org


//connect to peer
// mac, replace path to files by your own
geth --datadir /Users/yifenghuang/at2chain init $HOME/go/src/bitbucket.org/at2chain/at2chain.json
geth --datadir /Users/yifenghuang/at2chain --port 3000

// from another console
geth attach ipc:geth.ipc

// add peer
admin.addPeer("enode://c2a82249a2aa89693ff159298f3476fb93db3e03be01e4f67cae8bb1fc3003671703dab557582e5dc566269c1166bcc7a4943c0e082daa9df5b0f544076fdecc@101.132.152.99:3000")

// new an account
personal.newAccount()

// start mining
personal.unlockAccount(eth.coinbase)
eth.defaultAccount = eth.coinbase
miner.start()

```







## installation

### ubuntu
```
// install golang
curl -OL https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz

// install go-ethereum
mkdir at2chain-data
mkdir at2chain-data/node1
mkdir at2chain-data/node2
git clone https://github.com/ethereum/go-ethereum.git $HOME/go/src/github.com
cd $HOME/go/src/github.com/go-ethereum
make all
sudo vi /etc/profile
// add line "export PATH=$PATH:/root/go-ethereum/build/bin" 
source /etc/profile
// installation is done
```

### start nodes
Follow this post to use puppeth initial a private ethereum

https://modalduality.org/posts/puppeth/




## run a ethereum docker

```
docker run -d --name ethereum-node -v /root/at2chain-data:/root \
           -p 8545:8545 -p 30303:30303 -p 3000:3000 \
           ethereum/client-go --fast --cache=512 --rpcaddr 0.0.0.0 \
           --datadir /root/node1 --port 3000
           
```