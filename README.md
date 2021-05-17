# twitter-stream-rule-manager
manager for GET or POST /2/tweets/search/stream/rules

### Homebrew (macOS and Linux)

```console
$ brew install mashiike/tap/twitter-stream-rule-manager
```

### Binary packages

[Releases](https://github.com/mashiike/twitter-stream-rule-manager/releases)

## Usage

```console
Usage: twitter-stream-rule-manager <flags> <subcommand> <subcommand args>

Subcommands:
        deploy           deploy
        diff             diff

Subcommands for help:
        commands         list all command names
        flags            describe all known top-level flags
        help             describe subcommands and their syntax


Use "twitter-stream-rule-manager flags" for a list of top-level flags
```

## Quick Start

twitter-stream-rule-manager  can easily manage for twitter API /2/tweets/search/stream/rules state by json file.

create json rule file as below.
```json
{
  "rules":[
    {
      "tag": "hoge not RT",
      "value": "hoge -RT"
    }
  ]
}
```

and, check diff

```console
$ twitter-stream-rule-manager diff --rules example.json --bearer $TWITTER_BEARER_TOKEN
```

deploy 

```console
$ twitter-stream-rule-manager deploy --rules example.json --bearer $TWITTER_BEARER_TOKEN
```
