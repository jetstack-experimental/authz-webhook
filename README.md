# Bitesize AuthZ Webhook

## Installation

`bitesize-authz-webhook` is packaged into docker container and can be found at
`geribatai/bitesize-authz-webhook:0.0.5`.

## Configuration

### Environment variables

* `LISTEN_PORT`
* `RULES_CONFIG` - path to `rules.hcl` file. Default - rules.hcl in current
directory.

### rules.hcl

Access rules are described in HCL format. 

## Changelog

* 0.0.5 - First open-source release. Supports HCL rules.