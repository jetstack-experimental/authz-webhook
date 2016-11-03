# Bitesize AuthZ Webhook

[![Build Status](https://travis-ci.org/jetstack-experimental/authz-webhook.svg?branch=master)](https://travis-ci.org/jetstack-experimental/authz-webhook)

## Installation

`authz-webhook` is packaged into docker container and can be found at
`jetstackexperimental/authz-webhook:latest`. Currently it does not support HTTPS
termination, so it is advised to run it behind HTTPS proxy.

## Configuration

## Kubernetes configuration

### Environment variables

* `LISTEN_PORT` - Port webhook listens requests on (Default: 8080)
* `RULES_CONFIG` - path to `rules.hcl` file. (Default: rules.hcl in current
directory).

### rules.hcl

Access rules are described in HCL format. Rules file is processed from the top,
and the first rule match found is returned as authorization status. If no match
is found, implicit deny rule is matched at the end.

```
access "allow" {
    user = "admin"
}

access "deny" {
    verb = "create"
}
```

## Changelog

* 0.0.5 - First open-source release. Supports HCL rules.
* 0.0.6 - Adds automated travis builds
* 0.0.7 - Be transparent if config file load fails
