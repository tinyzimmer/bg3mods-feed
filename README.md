# BG3 Mods Feed Generator

This is a simple application for serving RSS/Atom/JSON feeds for mods on the [Baldur's Gate 3 Mod site](https://baldursgate3.game/mods).

## Usage

```bash
bg3mods-feed [options]
```

The following flags are supported:

```
      --api-url string            The API URL to fetch mods from (default "https://embed.modhub.io/v1/games/6715/mods")
      --config string             Path to the configuration file (YAML, JSON, TOML, or HCL)
      --fetch-interval duration   The interval to fetch mods at (default 5m0s)
      --format string             The format to render the feed in (rss, atom, json) (default "atom")
      --listen string             The address to listen on (default ":8080")
      --max-feed-items int        The maximum number of feed items to render (default 100)
      --platform string           Platform to filter mods by (windows, mac, ps5, xboxseriesx)
      --sort string               The field to sort the feed by (default "recent")
      --tags strings              Tags to filter mods by
```

The configuration file is optional and follows the same format as the flags.
An example can be found in the `config.example.yaml` file.

The feed will be available at `/feed` on the listen address.
For example, if the listen address is `:8080`, the feed will be available at `http://localhost:8080/feed`.

Defaults provided via configuration can be overridden per request using query arguments.
The following query arguments are supported:

| Query Argument   | Description                                                  | Example                                                  |
| ---------------- | ------------------------------------------------------------ | -------------------------------------------------------- |
| `max_items`      | The maximum number of feed items to render                   | `http://localhost:8080/feed?max_items=10`                |
| `sort`           | The field to sort the feed by                                | `http://localhost:8080/feed?sort=popular`                |
| `platform`       | Platform to filter mods by                                   | `http://localhost:8080/feed?platform=windows`            |
| `tags`           | Tags to filter mods by                                       | `http://localhost:8080/feed?tags=Classes,Cheats,English` |
| `fetch_interval` | Overrides the fetch interval (how long a response is cached) | `http://localhost:8080/feed?fetch_interval=1h`           |
| `format`         | The format to render the feed in                             | `http://localhost:8080/feed?format=rss`                  |

Sort can be any of the fields returned by the upstream API.
For an exhaustive list, refer to the `json` tags in [this file](internal/mods/types.go).
The following predefined values are supported:

- `recent`: Sort by the most recent mods
- `last_updated`: Sort by the last updated mods
- `trending`: Sort by the trending mods
- `highest_rated`: Sort by the highest rated mods
- `popular`: Sort by the most popular mods
- `subscribers`: Sort by the most subscribed mods
- `alphabetical`: Sort mods by name

## Installation

_TODO_
