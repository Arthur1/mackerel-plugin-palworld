# mackerel-plugin-palworld

## Description

Mackerel metrics plugin to get informations of Palworld dedicated server.

## Synopsis

```sh
mackerel-plugin-palworld -password "admin_password"
```

## Installation

```sh
sudo mkr plugin install Arthur1/mackerel-plugin-palworld
```

## Setting for mackerel-agent

```
[plugin.metrics.palworld]
command = ["/opt/mackerel-agent/plugins/bin/mackerel-plugin-palworld", "-password", "admin_password"]
```

## Usage

### Options

```
  -host string
    	hostname (default "localhost")
  -password string
    	administrator passowrd
  -port string
    	RCON port (default "25575")
  -timeout string
    	dial timeout seconds (default "5s")
```
