![Zeep Xanflorp and his Miniverse ](/assets/zeep_and_miniverse.png)

Zeep
=====

Worker for workflow aware serverless platform.


Table of Contents
=================

- [Zeep](#zeep)
- [Table of Contents](#table-of-contents)
- [Architecture](#architecture)
  - [Overview](#overview)
  - [Zeep](#zeep-1)
  - [Kyle](#kyle)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Build Kyle](#build-kyle)
  - [Install Zeep and Runtime](#install-zeep-and-runtime)
  - [Run Zeep](#run-zeep)
  - [Add a function](#add-a-function)
  - [Invoke a function](#invoke-a-function)
  - [Invoke a workflow](#invoke-a-workflow)
  - [Destroy Zeep](#destroy-zeep)

Architecture
============

## Overview

![Architecture Overview](/assets/overview.png)
*<p align="center">Architecture Overview</p>*
 
_Zeep_ is a local agent that manages containers, storages and handles function invoke requests on a worker node.
_Kyle_ is an agent delegate that manages the runtimes takes over some of _Zeep_'s roles.

Consider executing a workflow (also a function) that consists of sequential functions using _Zeep_ CLI. _Zeep_ takes the request and assign task to _Kyle_ (It might not be _Kyle_. See details on [_Kyle_ section](#kyle)) and _Kyle_ sends invoke request for first function of the workflow to a runtime. After the execution of the first function is completes, the second function of the workflow is invoked by another runtime but in the same container by _Kyle_.

For a workflow that consists of parallel functions, those functions are must be executed in different containers for performance isolation. At this time, _Kyle_ sends invoke request for the first parallel function of the workflow to it's runtime and sends invoke request for rest of the parallel functions of the workflow to _Zeep_. In this case, _Zeep_ knows that the requested functions are part of the workflow, so bind same storage that bound to the container that execute the first function to the other containers that are going to execute rest of the functions.

As a result, all the functions in the same workflow are deployed in the same host so the latency between the execution of functions is very small (about 150ms in OpenWhisk, about 1~15ms in ours). In addition, if there is data dependency among functions, the data can be shared using host storage, thereby saving the cost of external storage.


## Zeep

_Zeep_ provides _Invokee_ and _Invoker_ services via gRPC. _Invokee_ service handles the requests related to the worker that responsible to execute a function. _Invoker_ service handles the requests related to manage function resources and invoke a function.

When a container managed by _Zeep_ is created in the pool, _Kyle_ run in the container connect to _Invokee_ service and _Zeep_ pauses the container. _Zeep_ unpauses this container when the invoke a function is requested by _Invoker_ service, copies function resources into the container, and then assigns invoke task to _Kyle_ that is connected through _Invokee_ service. After the execution is completes, the container is paused again.


## Kyle

_Kyle_ also provides _Invokee_ and _Invoker_ service as same as _Zeep_ provides. This is why we call _Kyle_ an agent delegate. When a runtime managed by _Kyle_ is created, the runtime connect to _Invokee_ service provided by _Kyle_. The _Invoker_ service provided by _Kyle_ is only used by its runtime that executes a workflow. When the invoke a sequential function is requested, _Kyle_ spawns a new runtime process, waits for connection from runtime through _Invokee_ servie, and then assigns invoke task to the new runtime. If the function is a parallel function, _Kyle_ simply forward the request to _Zeep_.

At the cold start, _Kyle_ checks if the function is a workflow. If the function is not a workflow, _Kyle_ requests for handover to _Invokee_ service provided by _Zeep_ and sends message to handover to runtime connected through its _Invokee_ service. After the handover procedure, the runtime is connected to _Invokee_ service provided by _Zeep_ directly and the _Kyle_ is shut down.

Getting Started
===============


## Prerequisites

- [Go](https://golang.org/doc/install) >= 1.11
- [Docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/) >= 1.39


## Build Kyle

```sh
# Download Kyle.
$ git clone https://github.com/mee6aas/kyle.git

# Build Kyle.
$ cd ./kyle
$ ./scripts/build.sh
```

We are looking for a way to distribute Kyle as a package. Currently only manual builds are supported.


## Install Zeep and Runtime

```sh
# Pull Zeep image from Docker Hub.
$ docker pull mee6aas/zeep

# Pull Runtime image from Docker Hub.
$ docker pull mee6ass/runtime-nodejs

# Unstall Zeep CLI
$ go get github.com/mee6aas/zeep/cmd/zeep

# Make sure the Zeep CLI is installed.
$ zeep --help
```
Currently only the runtime for Node.js is supported.

## Run Zeep

```sh
$ zeep agent serve
```

When Zeep starts, you can see that one Zeep container and one runtime container are created.

## Add a function

```sh
# Download sample functions.
$ git clone https://github.com/mee6aas/samples.git

# Add echo function.
$ cd ./samples/nodejs
$ zeep act add echo

# List added functions.
$ zeep act ls
> name:"echo" runtime:"mee6aas/runtime-nodejs:latest" added:"2019-12-24 07:03:51.534119732 +0000 UTC m=+370.172112759"
```

## Invoke a function

```sh
# Invoke echo with argument "foo".
$ zeep act invoke echo foo
> foo
```

## Invoke a workflow
```sh
# Add functions.
$ cd ./pingpong
$ zeep act add pinger
$ zeep act add ponger

# Invoke pinger with argument "ponger".
# The function pinger will invoke the function ponger.
# The ponger returns string "ponged".
$ zeep act invoke pinger ponger
> ping to ponger...pong from ponger: ponged.
```

## Destroy Zeep

```
$ zeep agent destroy
```
