# Next Smart Chain GraphQL API Server

GraphQL API server for NEXT Smart Chain powered blockchain network.

## Releases
Please check the [release tags](https://github.com/NextSmartChain/next-api-graphql/tags) to get more details and to download previous releases.

#### Version 0.3.1
This version connects with the Orion protocol on the NEXT Smart Chain blockchain. The SFC contract ABI bundled with the API is version 3.3.0.

## Requirements

### MongoDB installation

Persistent data are stored in a MongoDB database. Going through the installation and
configuration process of MongoDB is out of scope here, please consult
[MongoDB manual](https://docs.mongodb.com/manual/) to install and configure appropriate
MongoDB environment for your deployment of the API server.

## Building the source

Building `apiserver` requires a Go (version 1.13 or later). You can install
it using your favourite package manager. Once the dependencies are installed, run

```shell
make
```

The build output is ```build/apiserver``` executable.

You don't need to clone the project into $GOPATH, due to use of Go Modules you can
use any location.

## Running the API server

To run the API Server you need access to a RPC interface of a full NEXT Smart Chain node. Alternatively you can obtain access to a remotely running instance of NEXT on rpc.nextsmartchain.com.

We recommend using local IPC channel for communication between a NEXT Smart Chain node and the
API Server for performance and security reasons. Please consider security implications
of opening NEXT RPC to outside access, especially if you enable "personal" commands
on your node while keeping your account keys in the NEXT key store.