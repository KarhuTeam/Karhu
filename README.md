# Karhu
Karhu is a deployment management and infrastructure tool, designed to be simple to configure & use.

## Features

* Application deployment (custom, services, Docker)
* Build and deployment history
* Application configuration
* Server registration (custom, EC2, DigitalOcean)
* Monitoring (Collectd)
* Alerts (nagios)
* Logs searching (Logstash / ElasticSearch)

## Requirements

* Docker 1.9.1+
* Any modern Linux distribution that supports Docker 1.9.1+.
* RAM: 1GB+

## Running Karhu with Docker

/!\ WARNING this is a very early preview of Karhu, database may change, future release may be incompatible with previous ones.
We are currently working on testing this, and figure out what is OK and what have to change. Feel free to open an issue :) <3 :love:

docker-compose.yml
```
mongo:
   hostname: mongo
   image: mongo
   restart: always

elasticsearch:
    hostname: elasticsearch
    image: elasticsearch
    restart: always

logstash:
    hostname: logstash
    image: maxwayt/karhu-logstash
    restart: always
    ports:
    - "5044:5044"
    - "25826:25826/udp"
    links:
    - "elasticsearch:elasticsearch"

karhu:
    hostname: karhu
    image: maxwayt/karhu
    restart: always
    links:
    - "mongo:mongo"
    - "elasticsearch:elasticsearch"
    volumes_from:
    - "logstash"
    ports:
    - "80:8080"
    environment:
    - MGO_HOSTS=mongo
    - MGO_DB=karhu
    - STORAGE_DRIVER=fs
    - STORAGE_PATH=/data
    - PUBLIC_HOST=http://your-karhu-host.com
    - INFLUXDB_COLLECTD_HOST=your-karhu-host.com
    - INFLUXDB_COLLECTD_PORT=25826
    - LOGSTASH_IP=you-karhu-ip
    - ES_HOST=http://elasticsearch:9200
    - EMAIL_PROVIDER=mailgun
    - MAILGUN_DOMAIN=your-domain
    - MAILGUN_APIKEY=your-key
```

## License

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
