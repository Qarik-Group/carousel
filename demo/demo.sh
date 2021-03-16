#! /bin/bash

vars_store=$(mktemp)

faketime -f "-340d" bosh int --vars-store=$vars_store ./demo/manifest.yml 2>/dev/null

name_base="/$(bosh curl /info | jq -r '.name')/$(bosh int ./demo/manifest.yml --path /name)"

vars_store_get() {
    bosh int $vars_store --path "$1"
}

credhub set --name "${name_base}/carousel_demo_ca" \
	--type certificate \
	--root <(vars_store_get "/carousel_demo_ca/ca") \
	--certificate <(vars_store_get "/carousel_demo_ca/certificate") \
	--private <(vars_store_get "/carousel_demo_ca/private_key")

credhub set --name "${name_base}/carousel_demo_leaf" \
	--type certificate \
	--ca-name "${name_base}/carousel_demo_ca" \
	--certificate <(vars_store_get "/carousel_demo_leaf/certificate") \
	--private <(vars_store_get "/carousel_demo_leaf/private_key")

rm $vars_store

bosh -n deploy -d carousel-demo ./demo/manifest.yml
