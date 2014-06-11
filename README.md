# meeko-collector-circleci #

[![Build
Status](https://drone.io/github.com/meeko-contrib/meeko-collector-circleci/status.png)](https://drone.io/github.com/meeko-contrib/meeko-collector-circleci/latest)

Meeko collector for CircleCI webhooks

## Meeko Variables ##

* `LISTEN_ADDRESS` - the TCP network address to listen on; format [HOST]:PORT
* `ACCESS_TOKEN` - Token to be used for for webhook authentication. The token
  is expected to be set via a query parameter `token`, e.g. `https://example.com?token=secret`.

## Meeko Interface ##

This collector accepts CircleCI webhooks (HTTP POST requests) and forwards the
payload as `circleci.build` event.

## License ##

MIT, see the `LICENSE` file.
