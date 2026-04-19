---
title: gcredstash
date: 2026-04-11
lang: en
abstract: |
  **gcredstash** is a reimplementation of [credstash](https://github.com/fugue/credstash) in Go.

  Like the original, it allows for the management of versioned secrets using AWS [Key Management Service](https://docs.aws.amazon.com/kms/) (KMS) for encryption/decryption and [DynamoDB](https://docs.aws.amazon.com/dynamodb/) for storage.

  Unlike the original, gcredstash is a standalone binary, making its deployment simpler.

  This particular version of gcredstash is a fork of [the original by Genki Sugawara](https://github.com/winebarrel/gcredstash), aiming to modernise it, increase its compatibility with the original credstash, and fix bugs.

  With the advent of [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/), both credstash and gcredstash are less useful these days, but if you're still a heavy user of either, this fork might be for you.
---

# Quickstart

If you have the Go toolchain already configured, you can `go install` the binary:

```console
go install github.com/kgaughan/goreleaser/cmd/@latest
```

Alternatively, you can [download the latest release](https://github.com/kgaughan/gcredstash/releases/latest) where there are Linux (ARM64 and x86-64), macOS (ARM64 and x86-64), and Windows (ARM64 and x86-64) builds, or pull the latest container image and use that with:

```console
$ docker pull ghcr.io/kgaughan/gcredstash:latest
```

# Configuration

## Environment

Ensure you have `AWS_PROFILE` and `AWS_REGION` set correctly first if you're running this from your own machine. If you're running this from an AWS instance. Alternatively, you can ensure that `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` are exported within your environment.

You can specify the DynamoDB table to query using either the `--table` flag or the `GCREDSTASH_TABLE` or `CREDSTASH_DEFAULT_TABLE` environment variables. The default table name is `credential-store`.

The `get`, `getall`, `put`, and `template` subcommands expect you to specify a key if you're not using one with the default name/alias of `alias/credstash`. You can specify it with the `-k`/`--key` flag or the `GCREDSTASH_KMS_KEY` or `CREDSTASH_KMS_KEY` environment variables.

## IAM policies

Assuming your region is `us-east-1` and you're using the default table name of `credential-table`...

### Setup

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "dynamodb:CreateTable",
        "dynamodb:DescribeTable"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:dynamodb:us-west-1:<ACCOUNT NUMBER>:table/credential-store"
    },
    {
      "Action": [
        "dynamodb:ListTables"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```

### Reader

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "kms:Decrypt"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:kms:us-east-1:AWSACCOUNTID:key/KEY-GUID"
    },
    {
      "Action": [
        "dynamodb:GetItem",
        "dynamodb:Query",
        "dynamodb:Scan"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:dynamodb:us-east-1:AWSACCOUNTID:table/credential-store"
    }
  ]
}
```

### Writer

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "kms:GenerateDataKey"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:kms:us-east-1:AWSACCOUNTID:key/KEY-GUID"
    },
    {
      "Action": [
        "dynamodb:PutItem"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:dynamodb:us-east-1:AWSACCOUNTID:table/credential-store"
    }
  ]
}
```

# Commands

## setup

Set up the credential store.

## completion

Generate the autocompletion script for the specified shell.

## list

List credentials and their versions.

## get

Get a credential from the store.

## getall

Get all credentials from the store.

## put

Put a credential into the store.

## template

Generate a file from a template, embedding credentials into it.

## help

Show help about a command.

# Colophon

The site was built using [pandoc](http://pandoc.org/) using a custom theme inspired by that of [mdBook](https://rust-lang.github.io/mdBook/).
