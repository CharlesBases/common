#!/bin/zsh

ip=http://www.baidu.com         # IP
threads=10                      # 线程数量
duration=1m                     # 压测时间
connections=1000                # TCP连接数量

wrk -t $threads -d $duration -c $connections --latency $ip