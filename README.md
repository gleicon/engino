# eNGINo X

## Description

Engino is a NGINX vhost manager. It works by getting data either from *consul*, *etcd* or *redis hash*. It merges data from these databases into a predefined configuration template, updates NGINX configuration files and SIGHUP it. The expected config layout is Debian/Ubuntu /etc/nginx/sites-available / sites-enabled.

A couple of security measures were embedded into *engino* to avoid downtime. One of them is a configurable throttle for template changes. If you happen to signal engino many changes within a short period of time, it will batch them, wait and schedule a config file change for the future. The other measure is to run nginx -f to check for syntax errors and reverting to the latest safe change.

Engino is to be coupled with [Habitat](https://github.com/gleicon/habitat) to allow for flexible application provisioning.

## Usage

Engino depends on a folder on your service discovery database or a hash in redis named *engino/appname* after your application. This folder must have the template name and key/values that will be substituted in the template. There will be two timestamps, one for creation data and another for last change and a boolean flag to disable config generation.

Every config change will trigger a graceful restart on NGINX, but the frequency of triggering is monitored to avoid loops and minimize downtime.

Engino's user and group need to have priviliges over nginx. Run with

$ engino -c <db://host:port/> -a appname -t /opt/engino/templates [-r <max restarts per minute throttle, default 2>]

## Build
	$ cd src ; make

## Options
	-c connection string. redis://localhost:6379, etcd://localhost:4001, consul://localhost:8500
	-a application name. This folder will have a subfolder called engino
	-t template directory
	-r throttling, number of restarts per minute
	-n nginx config dir, debian/ubuntu layout

## Examples
	- etcd
		$ curl http://127.0.0.1:4001/v1/keys/myapp/engino/template -d value="myapp.conf"
		$ curl http://127.0.0.1:4001/v1/keys/myapp/engino/backend -d value="127.0.0.1:10000"
		$ curl http://127.0.0.1:4001/v1/keys/myapp/engino/servername -d value="myapp.localhost"
		$ engino -c etcd://127.0.0.1:4001 -t /opt/engino/templates -a myapp

	- redis
		$ redis-cli hset env db newdb
		$ redis-cli hset env cache newcache
		$ redis-cli hset env queue newqueue
		$ habitat -r 127.0.0.1:6379 env

	- consul
		$ curl -X PUT -d 'newdb' http://localhost:8500/v1/kv/env/db
		$ curl -X PUT -d 'newcache' http://localhost:8500/v1/kv/env/cache
		$ curl -X PUT -d 'newqueue' http://localhost:8500/v1/kv/env/queue
		$ habitat -c 127.0.0.1:8500 env

	- consul with a different app name
		$ curl -X PUT -d 'newdb' http://localhost:8500/v1/kv/newapp/db
		$ curl -X PUT -d 'newcache' http://localhost:8500/v1/kv/newapp/cache
		$ curl -X PUT -d 'newqueue' http://localhost:8500/v1/kv/newapp/queue
		$ habitat -k newapp -c 127.0.0.1:8500 env

	You can mix data coming from all sources too.

## Vagrant

There is a Vagrantfile and an install.sh provided to create a testing linux environment. The project directory will be at /vagrant and consul, etcd, redis and Golang are installed.


## Authors
	Gleicon <gleicon@gmail.com>

## License MIT
