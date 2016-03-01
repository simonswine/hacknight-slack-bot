# Slack Quote Bot

Written for February 2016's [West London Hack Night](http://www.meetup.com/West-London-Hack-Night/).

## Overview

This golang project tries to respond all incoming Slack messages with an
fitting quote. To be able to repsond accordingly it indexes a quote database
first. The indexing/matching is done using elasticsearch

## End of Hacknight state

* We used [tray.io](https://tray.io) for handling the Slack part
* Indexing was done with golang bleve library (was pretty slow, did not work 100% in the end)

## Example usage

![Screenshot](/screenshots/example_chat.png?raw=true "Screenshot")

## Build and run it locally

### Prerequisites

* You'll need [docker](https://docker.io) and [docker-compose](https://docs.docker.com/compose/overview/) installed locally.
* You need to download this [quote database](http://thewebminer.com/do) (Go to the very bottom of the page, then save to this directory with the name `quotes_all.csv`

### Build go and images

```
$ docker-compose build
elasticsearch uses an image, skipping
Building slackbot
Step 1 : FROM golang:1.5
 ---> 32854d8f16c4
Step 2 : RUN go get github.com/tools/godep
 ---> Using cache
 ---> 485da879b0cf
Step 3 : RUN mkdir /slackbot && chown nobody /slackbot
 ---> Using cache
 ---> 86de35465fc6
Step 4 : WORKDIR /slackbot
 ---> Using cache
 ---> dcb368fff2a0
Step 5 : USER nobody
 ---> Using cache
 ---> 286a290cf9fb
Step 6 : ADD quotes_all.csv /slackbot/
 ---> Using cache
 ---> aa042e742875
Step 7 : ADD main.go /slackbot/main.go
 ---> Using cache
 ---> 45127529c7a9
Step 8 : ADD Godeps/ /slackbot/Godeps/
 ---> Using cache
 ---> 0bb450f8b28f
Step 9 : RUN godep go build
 ---> Using cache
 ---> b1c4d94a42ca
Step 10 : CMD /slackbot/slackbot
 ---> Using cache
 ---> 5f8993406a52
Successfully built 5f8993406a52
```

### Run app and elasticsearch

```
$ docker-compose up
Creating hacknightslackbot_elasticsearch_1
Creating hacknightslackbot_slackbot_1
Attaching to hacknightslackbot_elasticsearch_1, hacknightslackbot_slackbot_1
slackbot_1      | time="2016-02-23T22:45:05Z" level=info msg="Initializing slack quote bot..." 
slackbot_1      | time="2016-02-23T22:45:05Z" level=info msg="Use slack api token: xoxb-22757066566-LtgYbLxQcDpedIKOe11PPFFH" 
slackbot_1      | time="2016-02-23T22:45:05Z" level=info msg="Use elasticsearch url: http://172.17.0.2:9200" 
slackbot_1      | time="2016-02-23T22:45:05Z" level=info msg="Elasticsearch http://172.17.0.2:9200 is not ready yet (1. try)" 
elasticsearch_1 | [2016-02-23 22:45:06,003][INFO ][node                     ] [Thinker] version[2.2.0], pid[1], build[8ff36d1/2016-01-27T13:32:39Z]
elasticsearch_1 | [2016-02-23 22:45:06,005][INFO ][node                     ] [Thinker] initializing ...
slackbot_1      | time="2016-02-23T22:45:06Z" level=info msg="Elasticsearch http://172.17.0.2:9200 is not ready yet (2. try)" 
elasticsearch_1 | [2016-02-23 22:45:06,664][INFO ][plugins                  ] [Thinker] modules [lang-expression, lang-groovy], plugins [], sites []
elasticsearch_1 | [2016-02-23 22:45:06,702][INFO ][env                      ] [Thinker] using [1] data paths, mounts [[/usr/share/elasticsearch/data (/dev/disk/by-uuid/85d441f8-cbd3-4f25-87f1-50cf5ec1d94f)]], net usable_space [11.7gb], net total_space [69.1gb], spins? [possibly], types [ext4]
elasticsearch_1 | [2016-02-23 22:45:06,702][INFO ][env                      ] [Thinker] heap size [990.7mb], compressed ordinary object pointers [true]
slackbot_1      | time="2016-02-23T22:45:07Z" level=info msg="Elasticsearch http://172.17.0.2:9200 is not ready yet (3. try)" 
slackbot_1      | time="2016-02-23T22:45:08Z" level=info msg="Elasticsearch http://172.17.0.2:9200 is not ready yet (4. try)" 
slackbot_1      | time="2016-02-23T22:45:09Z" level=info msg="Elasticsearch http://172.17.0.2:9200 is not ready yet (5. try)" 
elasticsearch_1 | [2016-02-23 22:45:10,008][INFO ][node                     ] [Thinker] initialized
elasticsearch_1 | [2016-02-23 22:45:10,008][INFO ][node                     ] [Thinker] starting ...
elasticsearch_1 | [2016-02-23 22:45:10,166][INFO ][transport                ] [Thinker] publish_address {172.17.0.2:9300}, bound_addresses {[::]:9300}
elasticsearch_1 | [2016-02-23 22:45:10,179][INFO ][discovery                ] [Thinker] elasticsearch/NVlrQSgDTK6nXcajZfaWTQ
slackbot_1      | time="2016-02-23T22:45:10Z" level=info msg="Elasticsearch http://172.17.0.2:9200 is not ready yet (6. try)" 
slackbot_1      | time="2016-02-23T22:45:11Z" level=info msg="Elasticsearch http://172.17.0.2:9200 is not ready yet (7. try)" 
slackbot_1      | time="2016-02-23T22:45:12Z" level=info msg="Elasticsearch http://172.17.0.2:9200 is not ready yet (8. try)" 
slackbot_1      | time="2016-02-23T22:45:13Z" level=info msg="Elasticsearch http://172.17.0.2:9200 is not ready yet (9. try)" 
elasticsearch_1 | [2016-02-23 22:45:13,261][INFO ][cluster.service          ] [Thinker] new_master {Thinker}{NVlrQSgDTK6nXcajZfaWTQ}{172.17.0.2}{172.17.0.2:9300}, reason: zen-disco-join(elected_as_master, [0] joins received)
elasticsearch_1 | [2016-02-23 22:45:13,311][INFO ][http                     ] [Thinker] publish_address {172.17.0.2:9200}, bound_addresses {[::]:9200}
elasticsearch_1 | [2016-02-23 22:45:13,312][INFO ][node                     ] [Thinker] started
elasticsearch_1 | [2016-02-23 22:45:13,336][INFO ][gateway                  ] [Thinker] recovered [0] indices into cluster_state
slackbot_1      | time="2016-02-23T22:45:14Z" level=info msg="Established connection to elasticsearch" 
elasticsearch_1 | [2016-02-23 22:45:14,583][INFO ][cluster.metadata         ] [Thinker] [quotes] creating index, cause [api], templates [], shards [5]/[0], mappings []
slackbot_1      | time="2016-02-23T22:45:14Z" level=info msg="Created index 'quotes' in elasticsearch" 
elasticsearch_1 | [2016-02-23 22:45:15,294][INFO ][cluster.routing.allocation] [Thinker] Cluster health status changed from [RED] to [GREEN] (reason: [shards started [[quotes][4], [quotes][4]] ...]).
slackbot_1      | time="2016-02-23T22:45:15Z" level=info msg="Indexing bulk 0/759 quotes 0 - 99" 
elasticsearch_1 | [2016-02-23 22:45:15,865][INFO ][cluster.metadata         ] [Thinker] [quotes] create_mapping [quote]
slackbot_1      | time="2016-02-23T22:45:16Z" level=info msg="Indexing bulk 1/759 quotes 100 - 199" 
[...]
slackbot_1      | time="2016-02-23T22:45:53Z" level=info msg="Indexing bulk 759/759 quotes 75900 - 75965" 
slackbot_1      | time="2016-02-23T22:45:53Z" level=info msg="Connecting to slack" 
slackbot_1      | time="2016-02-23T22:46:09Z" level=info msg="Incoming message 'what?'" 
slackbot_1      | time="2016-02-23T22:46:09Z" level=info msg="query elasticsearch for 'what?' hits=6500" 
```

To integrate it with Slack, we built a (tray.io)[https://tray.io] flow that uses Slack outgoing webhook on a specific channel. 

1. Each messages posted to the channel is sent to the flow's trigger. 
2. After checking if it's a valid user message (not from a bot), an http request is sent to golang part containing the full text of the message. 
3. If the response is positive, the flow finally posts it (the quote) in Slack.

![Screenshot](/screenshots/tray_flow.png?raw=true "Screenshot")
