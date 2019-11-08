#!/bin/bash
geth='geth'
coinbase='7e4f8dfad4de68ddbb4917ddce69263bbe4fd1bb'
bootnode1='enode://7cd5df57476c68e4b8dcfc15dbc2295b3fd3dc133780f02587a5029597c141be3215e7ded125139a2bbbb97aec85946f1508c5ac340e55d744c99b4f04c57ea3@139.199.180.239:7003'
bootnode2='enode://e5f709620c65e2d281f93088cc36f814896393d395d9bad90c91648cdcfcef00148d75ffd112096c22df2b6db770550412025e77b6f7df9a643fa769ef13a5b6@192.168.0.11:30303'
dir=$(pwd)


pid=$(ps -ef | grep $geth | grep -v grep | awk '{print $2}')
count=$(ps -ef | grep $geth | grep -v grep | wc -l)
if [ $count -ne 0 ];then kill $pid
fi
echo "count:"$count

nohup $dir/$geth \
    --datadir $dir/data/ \
    --rpc --rpcapi "db,eth,net,web3,personal"  --rpccorsdomain '*'  \
    --networkid 123456789 \
    --mine --minerthreads 1 --etherbase $coinbase \
    --bootnodes $bootnode1 \
    &

echo '\n --done--'
sleep 1;
ps -ef |grep geth
echo '---done---'