#! /bin/bash

set -e

# -------------------- wget -------------------- #
yum -y install wget

# -------------------- yum -------------------- #
rpm -e --nodeps yum
rm -rf /etc/yum.repos.d/CentOS-Base.repo
curl http://mirrors.aliyun.com/repo/Centos-7.repo > /etc/yum.repos.d/CentOS-Base.repo

#wget https://mirrors.aliyun.com/centos/7/os/x86_64/Packages/python-iniparse-0.4-9.el7.noarch.rpm
wget https://mirrors.aliyun.com/centos/7/os/x86_64/Packages/yum-3.4.3-161.el7.centos.noarch.rpm
wget https://mirrors.aliyun.com/centos/7/os/x86_64/Packages/yum-metadata-parser-1.1.4-10.el7.x86_64.rpm
wget https://mirrors.aliyun.com/centos/7/os/x86_64/Packages/yum-plugin-fastestmirror-1.1.31-50.el7.noarch.rpm
rpm -ivh --force --nodeps yum-*
rm -rf yum-*

yum clean all
yum makecache

# -------------------- git -------------------- #
yum install git -y

# -------------------- oh my zsh -------------------- #
yum install zsh -y
chsh -s /bin/zsh

reboot

sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"

# .zshrc
rm -rf /root/.zshrc
echo '
export ZSH="/root/.oh-my-zsh"

ZSH_THEME="ys"

# Command execution time stamp shown in the history command output.
HIST_STAMPS="yyyy-mm-dd"

source $ZSH/oh-my-zsh.sh

' > /root/.zshrc
source /root/.zshrc