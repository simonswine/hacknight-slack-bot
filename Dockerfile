FROM golang:1.5

# install godep
RUN go get github.com/tools/godep

# ensure right directory rights
RUN mkdir /slackbot && chown nobody /slackbot

WORKDIR /slackbot
USER nobody

ADD quotes_all.csv /slackbot/
ADD main.go /slackbot/main.go
ADD Godeps/ /slackbot/Godeps/

RUN godep go build

CMD /slackbot/slackbot
