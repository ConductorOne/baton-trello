![Baton Logo](./baton-logo.png)

# `baton-trello` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-trello.svg)](https://pkg.go.dev/github.com/conductorone/baton-trello) ![main ci](https://github.com/conductorone/baton-trello/actions/workflows/main.yaml/badge.svg)

`baton-trello` is a connector for built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Prerequisites

1. Follow [Atlassian Developer Guide](https://developer.atlassian.com/cloud/trello/guides/power-ups/managing-power-ups/) to create a New Custom Power-Up and generate a valid API key
2. Follow [Atlassian Support Guide](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/#:~:text=variable%20length%20instead.-,Create%20an%20API%20token,-API%20tokens%20with) to create an API token
3. Use the Trello API to get the ID of the organizations you want to sync:
   4. Go to the following URL in your browser (replace API-Key with your API key and Token with your access token):
   `https://api.trello.com/1/organizations/[organization-name]?key=[API-Key]&token=[API-token]`

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-trello
baton-trello
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_API_KEY=apiKey -e BATON_API_TOKEN=apiToken -e BATON_ORGS=trelloOrgs ghcr.io/conductorone/baton-trello:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-trello/cmd/baton-trello@main

baton-trello

baton resources
```

# Data Model

`baton-trello` will pull down information about the following resources:
- Users
- Organizations
- Boards

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-trello` Command Line Usage

```
baton-trello

Usage:
  baton-trello [flags]
  baton-trello [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --api-key string               required: The API key for your Trello account ($BATON_API_KEY)
      --api-token string             required: The API token for your Trello account ($BATON_API_TOKEN)
      --client-id string             The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string         The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                  The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                         help for baton-trello
      --log-format string            The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string             The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --organizations stringArray    required: Limit syncing to specific organizations ($BATON_ORGS)
  -p, --provisioning                 If this connector supports provisioning, this must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --ticketing                    This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                      version for baton-trello

Use "baton-trello [command] --help" for more information about a command.
```
