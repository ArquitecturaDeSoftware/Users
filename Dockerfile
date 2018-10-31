FROM golang

RUN mkdir -p /go/src/github.com/jscastelblancoh/users_service

ADD . /go/src/github.com/jscastelblancoh/users_service

RUN curl https://glide.sh/get | sh
  #RUN go get  github.com/canthefason/go-watcher
  #RUN go install github.com/canthefason/go-watcher/cmd/watcher

RUN cd /go/src/github.com/jscastelblancoh/users_service && glide install
RUN ls /go/src/github.com/jscastelblancoh/users_service/user/cmd


  #ENTRYPOINT  watcher -run github.com/jscastelblancoh/users_service/user/cmd -watch github.com/jscastelblancoh/users_service/user
ENTRYPOINT go run /go/src/github.com/jscastelblancoh/users_service/user/cmd/main.go