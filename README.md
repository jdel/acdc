# ACDC: A Continuous Docker Compose

## Goal

acdc aims at creating a way to use docker-compose as a lightweight continuous deployment endpoint.

## How

acdc provides a REST API to trigger Docker Compose actions as well as dedicated webhook receivers for, but not limited to, slack and docker hub.

For ease of development and deployment alike, acdc is written in Go, and is shipped as a single statically linked binary.

## Important considerations

acdc has not been develpped with maximum security in mind. In its current state, I do not recomment to run in a production environment.

acdc does not offer SSL support, you should always run it behind an SSL termination proxy, like nginx with letsencrypt certificates for example.

# Concepts

acdc introduces a few concepts that I will describe here for clarity:

  - Master Key: Is the key used for creating other API Keys, it is important that this keys remains very much secret.
  - API Key: In acdc, an API Key can be local or remote. Remote means it is tied to a git repository.

# Install instructions

## Binaries

Nightly builds from master are available on [Bintray](https://bintray.com/jdel/acdc/master).

Tagged builds are available in the [Releases](https://github.com/jdel/acdc/releases) page.

Unless you absolutely need a patch that has not been released yet, you should stick to tagged builds.

## Build from source

You will require go 1.8+ (untested with previous versions), and glide to handle dependencies.

```bash
mkdir -p $GOPATH/src/github.com/jdel/acdc/
git clone https://github.com/jdel/acdc.git $GOPATH/src/github.com/jdel/acdc/
cd $GOPATH/src/github.com/jdel/acdc/
glide install -v
go build
```

This will generate the `acdc` binary in $GOPATH/src/github.com/jdel/acdc

## Run with Docker

This repository is linked to Automated builds on [Docker Hub](https://hub.docker.com/r/jdel/acdc/tags/).

Tagged builds and master are generated automatically. Latest is tagged manually once I decide a version is stable enough.

The docker image can work only if it is given access to `/var/run/docker.sock`, this is why we __need to__ bind-mount it.

```bash
docker run -d --name acdc \
  -p 8080:8080 \
  -v ~/acdc/:/home/user/acdc/:rw \
  -v /var/run/docker.sock:/var/run/docker.sock:rw \
  --group-add 50 \
  jdel/acdc:latest
```

The `--group-add` bit is important, as this is what will grant access to `docker.sock`. 50 is the required value to work with Docker for Mac.

In order to find which value works for you, run:

```bash
docker run -ti --rm --name acdc \
  -v /var/run/docker.sock:/var/run/docker.sock:rw \
  --user root \
  jdel/acdc:latest ls -alh /var/run/docker.sock
```

This will output permissions of `/var/run/docker.sock` as seen by the container.

Of course, this is not meant for production, so you could also just run it with `--user root` instead.

# Usage

```
A Continuous Docker Compose provides a docker-compose REST API and hooks for Slack, Docker Hub, Github and more.

Usage:
  acdc [flags]
  acdc [command]

Available Commands:
  api-key     Make operations on api-keys
  help        Help about any command
  status      Get the status of acdc
  version     Get the version of acdc

Flags:
      --compose-dir string   compose directory (default is $HOME/acdc/compose/) (default "compose")
  -C, --config string        config file (default is $HOME/acdc/config.yml)
  -h, --help                 help for acdc
  -H, --home string          acdc home (default is $HOME/acdc/
  -l, --log-level string     log level [Error,Warn,Info,Debug] (default "Error")
  -m, --master-key string    Master API key
  -p, --port int             port to listen to (default is 8080) (default 8080)
      --static string        prefix to serve static images (defaults to /static/) (default "static")
      --static-dir string    static directory (default is $HOME/acdc/static/) (default "static")

Use "acdc [command] --help" for more information about a command.
```

## First launch

Upon firt launch, acdc will generate its own config in ~/acdc/acdc.yml file based on command line parameters it received.

acdc will generate a master key in this config file, it's up to you to change it if needed.

Command line parameters will always override config values that have been set in the configuration file.

## Routes

| Route                       | Method | Auth       | Description                                         |
| --------------------------- | ------ | ---------- | --------------------------------------------------- |
| /about                      | GET    | N/A        | Shows acdc version                                  |
| /slack                      | POST   | API Key    | Receives hooks from Slack                           |
| /dockerhub/{apiKey}/{tag}   | POST   | API Key    | Receives Docker hub hooks                           |
| /github                     | POST   | API Key    | Receives Github hooks                               |
| /v1/key/new                 | POST   | Master Key | Creates a new API Key                               |
| /v1/key/{apiKey}            | DELETE | Master Key | Deletes the API Key                                 |
| /v1/key                     | GET    | Master Key | Lists all API Keys                                  |
| /v1/key/pull                | GET    | API Key    | Git pulls the repository attached to the remote key |
| /v1/{hookName}/up           | GET    | API Key    | Executes docker-compose up                          |
| /v1/{hookName}/down         | GET    | API Key    | Executes docker-compose down                        |
| /v1/{hookName}/start        | GET    | API Key    | Executes docker-compose start                       |
| /v1/{hookName}/stop         | GET    | API Key    | Executes docker-compose stop                        |
| /v1/{hookName}/restart      | GET    | API Key    | Executes docker-compose restart                     |
| /v1/{hookName}/logs         | GET    | API Key    | Executes docker-compose logs                        |
| /v1/{hookName}/pull         | GET    | API Key    | Executes docker-compose pull                        |
| /v1/{hookName}              | GET    | API Key    | Executes docker-compose ps                          |
| /v1/{hookName}              | POST   | API Key    | Uploads a new hook                                  |
| /v1/{hookName}              | DELETE | API Key    | Deletes an existing hook                            |
| /v1/                        | GET    | API Key    | Lists all hooks                                     |

## Generate API Keys

### New API Key rfom the command line

This command generates a local key:

```bash
$ acdc api-key new 
k68sVV7pBvwYR3n0
```

In order to use it, hooks need to be uploaded first

This command creates an api-key linked to a git repository so hooks are managed remotely:

```bash
$ acdc api-key new -r https://github.com/jdel/acdc-recipes
WOZVO5xRfx0Zm4sh 	 https://github.com/jdel/acdc-recipes
```

This command creates an api-key with a known unique `GSukJLa3LYR4ypks1nowEHrX`:

```bash
$ acdc api-key new -u GSukJLa3LYR4ypks1nowEHrX
GSukJLa3LYR4ypks1nowEHrX
```

This is useful for generating api-keys to work with Slack's auto generated hooks:

### New API Key from the API

The examples below assume you are running acdc behind a SSL termination proxy and that `JkCilNGK-yGgVNRtdQHZyg==` is the master key.

Feel free to replace with `localhost:8080` for testing purpose.

The same commands as above can be executed from the API using the Master Key:

Local API Key:

```bash
$ curl -XPOST -u api-key:JkCilNGK-yGgVNRtdQHZyg== https://acdc.yourdomain.net/v1/key/new
{"message":["Created key"],"key-unique":"W_TGCBY7DowX4vjI"}
```

Remote API Key:

```bash
curl -XPOST -u api-key:JkCilNGK-yGgVNRtdQHZyg== https://acdc.yourdomain.net/v1/key/new -F 'remote=https://github.com/jdel/acdc-recipes'
{"message":["Created key"],"key-unique":"URPvGI5qrqPRxAqZ"}
```

Remote API Key with known unique:

```bash
curl -XPOST -u api-key:JkCilNGK-yGgVNRtdQHZyg== https://acdc.yourdomain.net/v1/key/new -F 'remote=https://github.com/jdel/acdc-recipes' -F 'unique=GSukJLa3LYR4ypks1nowEHrX'
{"message":["Created key"],"key-unique":"GSukJLa3LYR4ypks1nowEHrX"}
```

## Getting started scenario

In order to get you started quickly, let's create a remote key linked to a git repository with docker-compose files in it:

```bash
curl -XPOST -u api-key:JkCilNGK-yGgVNRtdQHZyg== https://acdc.yourdomain.net/v1/key/new -F 'remote=https://github.com/jdel/acdc-recipes'
{"message":["Created key"],"key-unique":"URPvGI5qrqPRxAqZ"}
```

Let's use that new key to start the redis hook (be patient, docker is probably pulling the redis image !) :

```bash
curl -u api-key:URPvGI5qrqPRxAqZ https://acdc.yourdomain.net/v1/redis/up
{
  "message": [
    "Creating network \"urpvgi5qrqprxaqzredis_default\" with the default driver",
    "Creating urpvgi5qrqprxaqzredis_redis_1 ... \r",
    "Creating urpvgi5qrqprxaqzredis_redis_1",
    "\u001b[1A\u001b[2K\rCreating urpvgi5qrqprxaqzredis_redis_1 ... \u001b[32mdone\u001b[0m\r\u001b[1B"
  ],
  "hook-name": "redis"
}
```

Now, let's check the status of the redis hook:

```bash
curl -u api-key:URPvGI5qrqPRxAqZ https://acdc.yourdomain.net/v1/redis
{
  "message": [
    "            Name                           Command               State    Ports   ",
    "---------------------------------------------------------------------------------",
    "urpvgi5qrqprxaqzredis_redis_1   docker-entrypoint.sh redis ...   Up      6379/tcp ",
    ""
  ],
  "hook-name": "redis"
}
```

## Slack Hooks

In Slack, create a new Outgoing WebHooks integration. Fill in all the fields, and in the URLs section, add:

```
https://acdc.yourdomain.net/slack
```

Slack will have generated a token for you, and unfortunately, it cannot be overridden.

We will need to create a known unique API Key in acdc matching the Slack generated token:

```bash
curl -XPOST -u api-key:JkCilNGK-yGgVNRtdQHZyg== https://acdc.yourdomain.net/v1/key/new -F 'remote=https://github.com/jdel/acdc-recipes' -F 'unique=GSukJ4asLYPOy3kh1nlwEHrX'
{"message":["Created key"],"key-unique":"GSukJ4asLYPOy3kh1nlwEHrX"}
```

In the example above, `JkCilNGK-yGgVNRtdQHZyg==` is the acdc master key, and `GSukJ4asLYPOy3kh1nlwEHrX` is the Slack token.

Assuming your channel is #acdc and your trigger word is `acdc`, type the following in the #acdc channel:

```
acdc redis up
acdc redis logs
acdc redis down
```

# Backup and recovery

As everything is stored as plain files, you can use your favorite backup solution to keep your API Keys and hooks safe.

# Known caveats and limitations

  - The output is not pretty and not anonymized
  - Probably the best looking / best coded API
  - There are no tests
  - Documentation is not complete

All the points above will be remediated at some point

# Why am I not using libcompose ?

I tried to use [libcompose](https://github.com/docker/libcompose) and I got lost in dependency hell and couldn't get anything to work. Somemthing to look at, but not until they release 0.5.0.