FROM golang:1.19.0-bullseye
MAINTAINER Alex Levinson <alexl@uber.com>

ENV GSCRIPT_REVISION master
ENV GO111MODULE=auto
ENV GARBLE_VERSION=v0.9.0
ENV GOPATH='/root/go'
ENV GOROOT='/usr/local/go'

RUN mkdir -p $GOROOT && \
    git clone https://github.com/susMdT/gscript $GOPATH/src/github.com/susMdT/gscript && \
    cd $GOPATH/src/github.com/susMdT/gscript && \
    go get ./... && \
    cd /usr/local/go/src/ && \
    go get github.com/susMdT/gscript && \
    cd - && \
    cd cmd/gscript && \
    go mod tidy && \
    go install -a

WORKDIR /root

RUN DEBIAN_FRONTEND=noninteractive \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
      file build-essential vim nano wget curl sudo net-tools pwgen locales \
      git-core logrotate software-properties-common && \
    locale-gen en_US en_US.UTF-8 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN go install mvdan.cc/garble@${GARBLE_VERSION}

ADD ps1.sh /etc/profile.d/Z1_PS1.sh

RUN chmod +x /etc/profile.d/Z1_PS1.sh && \
    echo "" >> /root/.bashrc && \
    echo "GSCRIPT_REVISION=$GSCRIPT_REVISION" >> /root/.bashrc && \
    echo "source /etc/profile.d/Z1_PS1.sh" >> /root/.bashrc

RUN git clone https://github.com/scopatz/nanorc.git /usr/share/nano-syntax-highlighting/

ADD nanorc /etc/nanorc

RUN git clone --depth=1 https://github.com/amix/vimrc.git /opt/vim

ADD vimrc /root/.vimrc

RUN echo "autocmd BufNewFile,BufRead *.gs set syntax=javascript" >> /opt/vim/my_configs.vim
RUN echo "set tabstop=4" >> /opt/vim/my_configs.vim
RUN echo "set shiftwidth=4" >> /opt/vim/my_configs.vim
RUN echo "set expandtab" >> /opt/vim/my_configs.vim

VOLUME /root/share

ENV HOME /root
WORKDIR /root
