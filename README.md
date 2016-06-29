# Rhobot - a database devops command line tool
[![Build Status](https://travis-ci.org/cfpb/rhobot.svg?branch=develop)](https://travis-ci.org/cfpb/rhobot)

**Description**:  

Rhobot is a tool for generalized stateless automation of databases, leveraging Github, and third party application API's to apply DevOps principles to database development.  It is designed to work with the GoCD continuous delivery tool.

  - **Technology stack**: Rhobot is written in Go, and currently supports interactions with the GoCD application and Postgres databases.
  - **Status**:  Pre-alpha, view the change log here: [CHANGELOG](CHANGELOG.md).

## Dependencies

To fully utilize Rhobot, one must be able to compile the code using Go, and have a working GoCD server and Postgres server to develop against.

## Installation

To build go, run `make`.  See [INSTALL](INSTALL.md) for additional information.

## Configuration

Run rhobot --help to see configuration options.

## Usage

Run rhobot from the command line, usage can be found by running `rhobot` without any options.

## How to test the software

To run tests, run `make test`

## Known issues

WIP

## Getting help

If you have questions, concerns, bug reports, etc, please file an issue in this repository's Issue Tracker.

## Getting involved

To contribute, use the usual git flow of fork -> feature branch -> pull request.  See [CONTRIBUTING](CONTRIBUTING.md) for more details.

----

## Open source licensing info
1. [TERMS](TERMS.md)
2. [LICENSE](LICENSE)
3. [CFPB Source Code Policy](https://github.com/cfpb/source-code-policy/)
