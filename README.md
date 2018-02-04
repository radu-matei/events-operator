Kubernetes External Events Operator
===================================

[![CircleCI](https://circleci.com/gh/radu-matei/events-operator.svg?style=shield&circle-token=14627daadeee06639298258d0a110d360a360d00)](https://circleci.com/gh/radu-matei/events-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/radu-matei/events-operator)](https://goreportcard.com/report/github.com/radu-matei/events-operator)
[![GoDoc](https://godoc.org/github.com/radu-matei/events-operator?status.svg)](https://godoc.org/github.com/radu-matei/events-operator)



What is this?
-------------

This is a [Kubernetes operator][1] that wants to bring external events into Kubernetes. It consists of a [CRD (CustomResourceDefinition)][2] and a controller and its purpose is to **automatically subscribe to various external event providers** (events from cloud providers (storage, database updates), webhooks, pub/sub systems and other event sources) and **provide a consistent way of handling these events**.


Disclaimer
----------

This is not an official Microsoft project.

[1]: https://coreos.com/operators/
[2]: https://kubernetes.io/docs/concepts/api-extension/custom-resources/
