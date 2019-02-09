relay
=========

[![Build Status](https://img.shields.io/travis/com/akerl/relay.svg)](https://travis-ci.com/akerl/relay)
[![GitHub release](https://img.shields.io/github/release/akerl/relay.svg)](https://github.com/akerl/relay/releases)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-green.svg)](https://tldrlegal.com/license/mit-license)

A simple AWS Lambda webhook relay. It takes a config file with webhook IDs, each with a configured list of targets. When it receives a request for an ID, it sends a request to all listed targets for that ID.

## Usage

1. Create a Lambda with the payload.zip generated in the "Installation" section below.
2. In the Environment Variables for the Lambda, set `S3_BUCKET` and `S3_KEY` to refer to an S3 bucket and S3 key where you will store the configuration file. The bucket/file must be readable by the Lambda.
3. Create a file at that bucket/key with the configuration:

```
---
webhooks:
  c4eb2149-be1c-4c7c-bdc8-9615ce17c6cf:
    targets:
      - url: https://example.org
        method: POST
      - url: https://example.com/path
  156f5b8e-8c8f-4baf-ac29-e9ffbccfafe5:
    targets:
      - url: https://www.example.org
```

You can list as many targets as you'd like for each hook. The `method` setting is optional, and defaults to `GET`.

## Installation

The methods below describe how to create a payload.zip that can be used for AWS Lambdas.

### Official build process

This requires that you have Docker installed and running. It will launch a Docker b
uild container, build the binary, and create a zip file for loading into AWS Lambda
. The zip file can be found at `./pkg/payload.zip`.

```
make
```

### Local pkgforge build

This doesn't require Docker but does require that you have [the pkgforge gem](https://github.com/akerl/pkgforge) installed. It builds a zip file at `./pkg/payload.zip`

```
pkgforge build
```

### Local manual build

This method has no deps other than golang, make, and zip. You have to manually create the zip file.

```
make local
cp ./bin/relay_linux ./main
zip payload.zip ./main
```

## License

relay is released under the MIT License. See the bundled LICENSE file for details.
