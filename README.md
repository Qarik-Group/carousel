# Carousel

A BOSH aware cli tool for managing the rotation of credentials stored in CredHub.

## Installation

As with all go applications you can create the binary by:

* `git clone https://github.com/starkandwayne/carousel`
* `cd carousel`
* `go build`

you may then run it using `./carousel -help` or move it under your path i.e. `cp carousel /usr/local/bin`

## Usage

To be able to talk to BOSH and CredHub the following environment variables need to set:

```
# BOSH
export BOSH_ENVIRONMENT=https://{bosh_director_ip}:25555
export BOSH_CLIENT={bosh_director_uaa_client}
export BOSH_CLIENT_SECRET={bosh_director_uaa_client_secret}
export BOSH_CA_CERT="-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----"

# CredHub
export CREDHUB_SERVER=https://{bosh_director_ip}:8844
export CREDHUB_CLIENT={credhub_uaa_client}
export CREDHUB_SECRET={credhub_uaa_client_secret}
export CREDHUB_CA_CERT="-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----"
```

When using [bosh-bootloader](https://github.com/cloudfoundry/bosh-bootloader) the above
can be achieved by running `eval "$(bbl print-env)"` in your terminal.

When using [BOSH Genesis Kit](https://github.com/genesis-community/bosh-genesis-kit) the same can be achieved by running `eval "$(genesis do environment-name-file.yml -- print-env)"`

### Browse

To make it easier to debug credential, and in particular certificate issues, carousel
provides an interactive terminal UI. Which gives the user an simpel way of browsing
trought certificate signing chains in a tree like fashion.

```
carousel browse
```

### List

List CredHub credentials augmented with information from the BOSH director:
* update_mode: looked up from runtime configs and deployment manifest 'variables:' sections
* deployments: list of deployment names which use this version of the credential

```
carousel list [flags]

Flags:
  -d, --deployments strings   filter by deployment names (comma separated)
  -h, --help                  help for list
	  --include-all           also show unused credential versions
	  --signing               only show Certificates used to sign
  -t, --types strings         filter by credential type (comma sperated) (default [certificate,ssh,rsa,password,user,value,json])
```

### Update Transitional

TODO

### Regenerate

TODO

### Remove Unused

TODO

## carousel-concourse

Write a description of the resource here.

## Source Configuration

* `a`: *Required.* This is a required setting.

* `b`: *Optional.* This is an optional setting.

* `c`: *Optional. Default `true`* This is an optional setting with a default value.

### Example

```yaml
resource_types:
- name: carousel
  type: registry-image
  source:
	repository: starkandwayne/carousel-concourse

resources:
- name: carousel
  type: carousel
  check_every: 5m
  source:
	log_level: debug

jobs:
- name: do-it
  plan:
  - get: carousel
	trigger: true
  - put: carousel
	params:
	  version_path: carousel/version
```

## Behavior

### `check`: Check for something

Write a description of what is checked here.

### `in`: Fetch something

Write a description of what is fetched here.

#### Parameters

* `a`: *Required.* This is a required parameter.

* `b`: *Optional.* This is an optional parameter.

### `out`: Put something somewhere

Write a description of what is being put somewhere.

#### Parameters

* `a`: *Required.* This is a required parameter.

* `b`: *Optional. Default `true`* This is an optional parameter with a default value.

## Development

### Prerequisites

* golang is *required* - version 1.11.x or higher is required.
* docker is *required* - version 17.05.x or higher is required.
* make is *required* - version 4.1 of GNU make is tested.

### Running the tests

The Makefile includes a `test` target, and tests are also run inside the Docker build.

Run the tests with the following command:

```sh
make test
```

### Building and publishing the image

The Makefile includes targets for building and publishing the docker image. Each of these
takes an optional `VERSION` argument, which will tag and/or push the docker image with
the given version.

```sh
make VERSION=1.2.3
make publish VERSION=1.2.3
```
