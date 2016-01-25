# Karhu
Karhu is a deployment management and infrastructure tool, designed to be simple to configure & use.

## Planned features

Non exhaustive of planned feature:
* Build registration for a project
* Click to deploy a build
* Link applications and services
* Scripts / Configurations
* Infrastructure declaration
* Infrastructure monitoring
* Servers and applications logs aggregation
* Alerts based on logs or monitoring

## Implemented

* Applications creation
* Services creation
* Applications and services linking
* Build registration for applications (Only binary supported atm)
* Server node registration
* Application deployment based on tags (deployed on matching tags servers)
* Deployment streaming through websocket

## Running Karhu with Docker

/!\ WARNING this is a very early preview of Karhu, database may change, future release may be incompatible with previous ones.
We are currently working on testing this, and figure out what is OK and what have to change. Feel free to open an issue :) <3 :love:

MongoDB database
```
docker run --name=karhu-mongo --restart=always -d mongo
```

Karhu
```
docker run --name=karhu --restart=always -p 80:8080 --link karhu-mongo:mongo -e PUBLIC_HOST=http://you-public-host.com -d maxwayt/karhu
```

Available env options:
* `DEBUG`: debug mode, default `1`
* `MGO_HOSTS`: mongo database hosts, default `mongo`
* `MGO_DB`: mongo database name, default `karhu`
* `PUBLIC_HOST`: Karhu public host, default `http://127.0.0.1:8080`
* `STORAGE_DRIVER`: Karhu file storage driver, only `fs` is supported for now, default `fs`
* `STORAGE_PATH`: storage directory, default `/data`
