#! /bin/zsh

# -------------------- scp -------------------- #
user=root
ipaddr=127.0.0.1

path=`pwd`
scp -rp -C ${path} ${user}@${ipaddr}:/root/