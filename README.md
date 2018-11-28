# at2chain v1_3砹链开发环境部署指南

## install from scratch

### install docker, docker-compose, golang
```
git clone https://github.com/sh777/abc
sh abc/docker-install.sh
sh abc/docker-compose-install.sh
sh abc/golang-install.sh
```
re-login is required here

### prepare blockchain environment

```
cd $GOPATH/src
mkdir -p bitbucket.org
cd bitbucket.org
git clone https://bitbucket.org/at2plus/at2chain-v1_3.git
cd at2chain-v1_3
sh scripts/download-cmds.sh
sh scripts/download-images.sh
```

### start docker-compose, channel and chaincode will be installed and instantiated ###
```
docker-compose up -d
# or start with explorer
docker-compose -f docker-compose-with-explorer.yaml up -d

```

### test the chain ###
```
docker exec -it cli bash
peer chaincode invoke -n at2chain --tls true --cafile $ORDERER_CA -c '{"Args":["uploaddoc","15","30a492e0c5335b8032747507c42f622f47ff4a2f","hyf","Desc###123"]}' -C tracec
peer chaincode invoke -n at2chain --tls true --cafile $ORDERER_CA -c '{"Args":["querydoc","fba9d88164f3e2d9109ee770223212a0"]}' -C tracec
```

### blockchain explorer ###

open http://host-name:8080 to explore the chain


## trouble shooting ##
### 安装路径 ###
at2chain必须严格按照要求安装在: $GOPATH/src/bitbucket.org/at2chain

### 在国内服务器上部署时的问题总结 ###
* golang如果由于网络问题安装失败, 使用这个链接提供的第三方工具安装方案: https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/01.1.md
* docker如果由于网络问题安装失败, 使用这个链接提供的方案: https://blog.csdn.net/u011365831/article/details/78851663
* download-cmds.sh如果由于超时执行失败, 替换使用download-cmds-from-mirror.sh
* 在第一次安装智能合约的时候, 由于需要从海外获取大量文件来创建新的智能合约镜像, 创建过程比较慢, 可能需要半个小时左右的时间, 直到使用docker ps命令可以看到类似于“dev-peer0.core.at2chain.com-at2chain-...”这样的镜像运行中, 智能合约才能正常被调用, 在此之前, 调用智能合约会得到“chaincodeid not exist”之类的错误信息

### 停止开发环境 ###
命令行切换到at2chain目录, 执行:
```
docker-compose down
```

### 重新启动开发环境 ###
重新启动开发环境有时必须先清除镜像缓存, 否则会导致不可预料的异常, 目前已知的清除过程为:
```
# stop the compose if not yet
docker-compose down

# remove all containers
docker rm $(docker ps -a -q)

# remove all images
docker rmi -f $(docker images -q)

sh scripts/download-images.sh

# remove and download the repo again
cd ..
rm -rf at2chain
git clone https://bitbucket.org/AI-CORE/at2chain

# reboot
sudo reboot
```

### 重新生成创世块 ###
如果由于配置修改必须重新生成创世块, 需要严格按照以下顺序执行
```
# remove existing artifacts
cp scripts/reset-codebase.sh ../
cd ..
sh reset-codebase.sh
cd at2chain

# initial the chain
sh scripts/start-at2chain.sh

# commit new artifacts
git add . -A
git commit -m "update artifacts"
git push

# start the chain
docker-compose up -d
```



