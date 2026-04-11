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

# Commands

## setup

Setup the credential store.

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

# Configuration
