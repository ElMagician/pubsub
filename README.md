# pubsub [![PkgGoDev][godoc_img]]() [![Coverage Status][coverage_img]][coverage] [![Build Status][status_img]][status] [![Go Report Card][report_img]][report]

This package aims to simplify pubsub management pubsub implementation in go programs. It allows to regroup any pubsub
client through a single interface lessening the burdening of provider switching.

## Usage

- `go get github.com/elmagician/pubsub`

## Implementations

- Mock: testify mock implementation for unit testing
- GCP: google pubsub implementation


[//]: <> (Badges links and images)
[coverage]: https://pkg.go.dev/github.com/elmagician/pubsub?tab=overview
[coverage_img]: https://coveralls.io/repos/github/elmagician/pubsub/badge.svg?branch=main

[status]: https://github.com/elmagician/pubsub/actions
[status_img]: https://github.com/elmagician/pubsub/workflows/CI/badge.svg

[report]: https://goreportcard.com/report/github.com/elmagician/pubsub
[report_img]: https://goreportcard.com/badge/github.com/elmagician/pubsub

[godoc_img]: https://pkg.go.dev/badge/github.com/elmagician/pubsub?tab=overview