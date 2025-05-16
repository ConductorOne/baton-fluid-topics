![Baton Logo](./baton-logo.png)

# `baton-fluid-topics` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-fluid-topics.svg)](https://pkg.go.dev/github.com/conductorone/baton-fluid-topics) ![main ci](https://github.com/conductorone/baton-fluid-topics/actions/workflows/main.yaml/badge.svg)

`baton-fluid-topics` is a connector for [Fluid-Topics](https://www.fluidtopics.com/) built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

## Prerequisites

In order to use this connector, you need an API key with ADMIN permissions, which is indicated by the `--bearer-token` flag and a domain with the `--domain` flag.
Example:
For connecting to https://example.fluidtopics.net you should do:

  ```
  baton-fluid-topics --bearer-token abcdefghij1234567890 --domain example
  ```

## Where can I find my API Key?
1- Log in [Fluid-Topics](https://www.fluidtopics.com/), then in the top right corner of the main page of your fluid topics page, click on administration.
2- In the menu that opens, click on integrations.
3- On the integrations page, below the list of api keys, a section will appear to create an api-key by adding a name and clicking create&add.
4- When you click on create&add it will open a menu of options to customize the apikey. 
5- After configuring your api key click on ok, and when the menu closes, in the integrations page with the apikey list click on save in the lower right corner.
   
Note: documentation of api keys: [Fluid-topics-APIKEY](https://doc.fluidtopics.com/r/Fluid-Topics-Configuration-and-Administration-Guide/Configure-a-Fluid-Topics-tenant/Integrations/API-keys)

# Connector capabilities
- Sync Users and Roles.
- Account provisioning:
    When you creating and new account, the following fields are required:
        - Name: The full display name of the user.
               Example: Name Example 
        - Email Address: The user email address. 
               Example: email@example.com
- Entitlements provisioning
- User usage

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-fluid-topics
baton-fluid-topics
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-fluid-topics:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-fluid-topics/cmd/baton-fluid-topics@main

baton-fluid-topics

baton resources
```

# Data Model

`baton-fluid-topics` will pull down information about the following resources:
- Users

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-fluid-topics` Command Line Usage

```
baton-fluid-topics

Usage:
  baton-fluid-topics [flags]
  baton-fluid-topics [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --bearer-token string          REQUIRED: The client secret token used to authenticate with ConductorOne
      --domain string                REQUIRED: Fluid topics account domain 
      --client-id string             The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string         The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                  The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                         help for baton-fluid-topics
      --log-format string            The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string             The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning                 If this connector supports provisioning, this must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --ticketing                    This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                      version for baton-fluid-topics

Use "baton-fluid-topics [command] --help" for more information about a command.
```
