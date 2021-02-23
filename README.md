# Carousel

A BOSH aware cli tool for managing the rotation of credentials stored in CredHub.

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

### Broswe

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
  -d, --deployments strings   filter by deployment names (comma sperated)
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
