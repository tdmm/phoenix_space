#!/usr/bin/env bash
#!bin/bash
pid=$(ps -ef | grep bin/ubridge | grep -v grep | awk '{print $2}')
count=$(ps -ef | grep bin/ubridge | grep -v grep | wc -l)
if [ $count -ne 0 ];then kill $pid
fi
echo "count:"$count
nohup env GOTRACEBACK=crash ./bin/ubridge.201809051632.f4bd480b4b455f14506115458d923fe36f4649f7 &
tail -f log/server.log