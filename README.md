# RVault <!-- omit in toc -->
[![codecov](https://codecov.io/gh/kir4h/rvault/branch/master/graph/badge.svg)](https://codecov.io/gh/kir4h/rvault)
![Build](https://github.com/kir4h/rvault/workflows/Build/badge.svg)
[![GitHub license](https://img.shields.io/github/license/kir4h/rvault)](https://github.com/kir4h/rvault/blob/master/LICENSE)
![GitHub top language](https://img.shields.io/github/languages/top/kir4h/rvault)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/kir4h/rvault?sort=semver)
![Docker Pulls](https://img.shields.io/docker/pulls/kir4h/rvault)
[![Github All Releases](https://img.shields.io/github/downloads/kotlin-graphics/kotlin-unsigned/total.svg)]()


* [Summary](#summary)
* [Motivation](#motivation)
* [Pre-requirements](#pre-requirements)
* [Installation](#installation)
* [Available commands](#available-commands)
  * [List](#list)
  * [Read](#read)
    * [File output](#file-output)
    * [JSON output](#json-output)
    * [YAML output](#yaml-output)
* [Filtering results](#filtering-results)
* [Configuration](#configuration)
* [Running in Docker](#running-in-docker)
* [Detecting KV version](#detecting-kv-version)

## Summary

RVault (standing for 'recursive vault') is a small cli utility aiming to perform some recursive actions over
[Hashicorp's Vault](https://www.vaultproject.io/).

The supported actions are:

* Listing secrets from a `kv` engine
* Reading secrets from a `kv` engine

Read secrets can be saved as files (each key/value will be stored to a separate file), or written to stdout in
`json`/`yaml` format.

## Motivation

On a scenario where I had to download multiple secrets from Vault, I realized this is currently not supported by `vault`
cli. There is an [open issue](https://github.com/hashicorp/vault/issues/5275) to request recursive key listing.
According to the discussion, recursivity is complex as current list operations are not optimized for recursion. So this
is just a workaround until a proper approach is implemented in Vault.

## Pre-requirements

In order for this tool to work you will need:

* A Vault instance to read from (kind of useless without it). I have tested it with Vault `1.4.2` but I guess it should
work with any Vault `1.x` version.
* A Token with enough priviledges to list the secrets from the selected engine. Additionally, the tool relies on
listing current mounts (through the `v1/sys/mounts` endpoint) to determine the kv version for the engine. If the token
doesn't have enough priviledges the version value should be provided either via the configuration or as cli argument.

## Installation

You can download the binary from the releases and uncompress it somewhere in your path.

For instance, to download the latest release to `/usr/local/bin` on Linux/MacOS:

```console
version=$(curl -s https://api.github.com/repos/kir4h/rvault/releases/latest | jq -r .tag_name)
platform=$(uname | tr '[:upper:]' '[:lower:]')
curl -L -s https://github.com/kir4h/rvault/releases/latest/download/rvault-${version}-${platform}-amd64.tar.gz \
| sudo tar xz -C /usr/local/bin
```

You can also simply run the tool from the `kir4h/rvault` Docker image.

## Available commands

### List

The `list` allows recursively listing secrets from a given path.

```console
$ rvault list --help

Recursively list secrets for a given path

Usage:
  rvault list <engine> [flags]

Flags:
  -h, --help                help for list
  -k, --kv-version string   KV Version
  -p, --path string         Path to look for secrets (default "/")

Global Flags:
  -a, --address string          Vault address
      --alsologtostderr         log to standard error as well as files
  -c, --concurrency uint32      Maximum number of concurrent queries to Vault (default 20)
      --config string           config file (default is $HOME/.config/rvault/config.yaml)
  -e, --exclude-paths strings   KV paths to be excluded (Applied on 'include-paths' output
  -i, --include-paths strings   KV paths to be included (default [*])
      --insecure                Enables or disables SSL verification
      --log_dir string          If non-empty, write log files in this directory
      --log_file string         If non-empty, use this log file
      --logtostderr             log to standard error instead of files (default true)
  -t, --token string            Vault token
  -v, --v Level                 number for the log level verbosity

```

For instance, in order to list secrets under `v2` `secret` engine (default one when Vault is launched in dev mode), for
the `spain` subpath:

```console
$ rvault list secret -k 2 -p spain -v=0
/spain/central/ssh.key
/spain/south/passwd.conf
```

### Read

The `read` command allows recursive read of secrets from a given path. Read secrets can be written to files or to
stdout in `json`/`yaml` format.

```console
$ rvault read --help

Recursively read secrets for a given path

Usage:
  rvault read <engine> [flags]

Flags:
      --file-permission uint32     Permissions for created secret files (file format only) (default 384)
      --folder-permission uint32   Permissions for newly created folders (file format only) (default 448)
  -f, --format string              Output format ('file', 'yaml', 'json') (default "file")
  -h, --help                       help for read
  -k, --kv-version string          KV Version
  -o, --output string              Output folder for 'file' format  (default ".")
  -w, --overwrite                  Overwrite existing files (file format only)
  -p, --path string                Path to look for secrets (default "/")

Global Flags:
  -a, --address string          Vault address
      --alsologtostderr         log to standard error as well as files
  -c, --concurrency uint32      Maximum number of concurrent queries to Vault (default 20)
      --config string           config file (default is $HOME/.config/rvault/config.yaml)
  -e, --exclude-paths strings   KV paths to be excluded (Applied on 'include-paths' output
  -i, --include-paths strings   KV paths to be included (default [*])
      --insecure                Enables or disables SSL verification
      --log_dir string          If non-empty, write log files in this directory
      --log_file string         If non-empty, use this log file
      --logtostderr             log to standard error instead of files (default true)
  -t, --token string            Vault token
  -v, --v Level                 number for the log level verbosity
```

#### File output

* Reading all secrets under `v2` `secret` engine for the `spain` subpath and store them as files under the `/tmp/secret`
folder

    ```console
    $ rvault read secret -f file -k 2 -o /tmp/secret -p spain
    I0712 19:04:03.791067   27410 root.go:84] Using config file: '/home/user/.config/rvault/config.toml'
    I0712 19:04:03.795243   27410 secrets.go:193] Secrets written: 2

    $ find /tmp/secret -type f
    /tmp/secret/spain/north/ssh.key/key
    /tmp/secret/spain/south/passwd.conf/key
    ```

    Please note that in above example, `key` is just the `key` name of the only `key`/`value` in the secret. If more
    than one `key`/`value` is found in the secret multiple files are created, each one having as name the corresponding
    `key` name.

#### JSON output

* Reading all secrets under `v2` `secret` engine for the `spain` subpath, returning them as `json`

    Running

    ```console
    rvault read secret -f json -k 2 -p spain 2>/dev/null | jq .
    ```

    Will produce the output

    ```json
    {
        "spain/north/ssh.key": {
            "key": "This is north's secret key"
        },
        "spain/south/passwd.conf": {
            "key": "This is south's secret key"
        }
    }
    ```

#### YAML output

* Reading all secrets under `v2` `secret` engine for the `spain` subpath, returning them as `yaml`

    Running

    ```console
    rvault read secret -f json -k 2 -p spain 2>/dev/null | jq .
    ```

    Will produce the output

    ```yaml
    spain/north/ssh.key:
        value: "This is north's secret key"
    spain/south/passwd.conf:
        value: "This is south's secret key"
    ```

## Filtering results

The `include-paths` and `exclude-paths` flags can be used in both `list` and `read` commands to return only secrets
matching a set of given patterns through wildcard patterns.

Let's say we have a secret engine with the following secrets:

```console
$ rvault list secret -k 2 -v=0
/france/paris/id_rsa
/france/vault.kdbx
/spain/north/id_rsa
/spain/south/shadow
/uk/london/bigben.crypt
```

We could return only secrets with the `id_rsa` name:

```console
$ rvault list secret -k 2 -v=0 -i */id_rsa
/france/paris/id_rsa
/spain/north/id_rsa
```

Or secrets in `france` and `uk`:

```console
$ rvault list secret -k 2 -v=0 -i /france/*,/uk/*
/france/paris/id_rsa
/france/vault.kdbx
/uk/london/bigben.crypt
```

The `include` filter can also be expressed as a repeteable flag:

```console
$ rvault list secret -k 2 -v=0 -i /france/* -i /uk/*
/france/paris/id_rsa
/france/vault.kdbx
/uk/london/bigben.crypt
```

We can use the `exclude-paths` to perform a second filter:

```console
$ rvault list secret -k 2 -v=0 -i /france/* -i /uk/* -e **/*.kdbx
/france/paris/id_rsa
/uk/london/bigben.crypt
```

The pattern matching is done through the [gobwas/go](https://github.com/gobwas/glob) globbing library, so you can
refer to its documentation for further reference and valid syntax.

## Configuration

Arguments can be provided as cli arguments but also through means of a configuration file.
[Viper](https://github.com/spf13/viper) is used for configuration management so configuration can be stored in any of
the formats supported by it (JSON, TOML, YAML, HCL, envfile and Java properties config files).

The default path for the configuration file is `$HOME/.config/rvault/config.<ext>`, where `<ext>` represents the
extension for the configuration format (`json`, `toml`, `yaml`, ...)

A sample of the configuration in `TOML` format:

```toml
[global]
# Vault address
address = "http://127.0.0.1:8200"
# Vault token
token = "devtoken"
# Log verbosity. 0 (less verbose) to 5 (more verbose)
verbosity = 2
# Maximum number of concurrent queries to Vault. '0' for unlimited (use with care)
concurrency = 20
# List of path patterns to return
include_paths = ["*"]
# List of paths to be excluded from the selected 'include_paths'
exclude_paths = []
# Default kv version to use
kv_version = ""
# Enables or disables SSL Verification
insecure = false

[list]
# Default path to use for listing
path = "/"

[read]
# Default path to use for reading
path = "/"
# Whether to overwrite existing secrets or not when using 'file' format
overwrite = false
# Permissions for newly created files when using 'file' format
file_permission = 0o0600
# Permissions for newly created folders when using 'file' format
folder_permission = 0o0700

#[engines]
## Properties for an engine named 'secret'
#[engines.secret]
## Set engine 'secret' as kv version '2'
## Please note that this takes precedence over 'global.kv_version'
#kv_version=2
## Properties for an engine named 'secret'
#[engines.secretkv1]
## Set engine 'secret' as kv version '1'
## Please note that this takes precedence over 'global.kv_version'
#kv_version=1
```

Additionally, environment variables `VAULT_ADDR` and `VAULT_TOKEN` are bound to `global.address` and `global.token`
configuration variables for convenience, since they might be already in place under a shell using vault cli.

## Running in Docker

A Docker image is built with this repository and pushed to Docker Hub
(see [rvault Docker image](https://hub.docker.com/repository/docker/kir4h/rvault))

There is nothing special about running it in docker other than bind mounting your configuration if desired, or the
output folder if using the `read` command with the `file` format.

Listing secrets for the `secret` engine with our custom configuration and the `spain` path:

```console
$ docker run -it --rm -v ${HOME}/.config/rvault/:/config kir4h/rvault  --config /config/config.toml list secret -p spain
I0715 22:19:36.175957       1 root.go:85] Using config file: '/config/config.toml'
/spain/central/ssh.key
/spain/south/passwd.conf
```

Reading those secrets and dumping them into files:

```console
$ docker run -it --rm -v ${HOME}/.config/rvault/:/config -v $(pwd)/out:/out kir4h/rvault --config /config/config.toml read secret -o /out
I0715 22:21:46.676216       1 root.go:85] Using config file: '/config/config.toml'
I0715 22:21:46.683996       1 output.go:77] Secrets written: 2
```

## Detecting KV version

If KV version is not provided, rvault will try to list the mounts by using the `v1/sys/mounts` endpoint. In order to do
so, the provided token must have enough priviledges for that query.

There is an alternative unauthenticated endpoint at `/sys/internal/ui/mounts` that could be used instead to retrieve
KV version, but since it's an internal endpoint it's not exposed by vault library and its usage is discouraged as there
is no guarantee on its backwards compatibility.
