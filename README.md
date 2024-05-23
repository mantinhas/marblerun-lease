# Modified-Modified Marblerun

This is a modified version of the Marblerun framework.

## Changes

Adds support for secure revocation. The coordinator can now revoke access to an application by contacting the Marblerun's coordinator. The coordinator will then send a message to the application modules to stop execution. The coordinator then also shuts down.

### Why is this 'secure'

We consider an attacker that controls the cluster network and can intercept and drop messages. This means that the attacker can intercept a revocation message from the provider (owner of the application and the entity that deploys the application on the cluster) to the coordinator, and can also intercept the messages from the coordinator to the application modules.

This means that we can't rely solely on the revocation message to stop the application modules. To achieve this, both the provider starts a simple GRPC server that listens for ping messages from the coordinator. The coordinator repeatedly sends ping messages to this server, with a configurable time interval. If the server doesn't respond after a set timeout, the coordinator will send a revocation message to the application modules and shut down. The coordinator also starts a GRPC server, where the application modules send a ping message to the coordinator. The coordinator will respond to the ping messages, and if it doesn't, the application modules shut down.

This mechanism ensures that the application modules will shut down even if the attacker intercepts the revocation message from the provider to the coordinator, since ping message will be timed out.

## Prerequisites

To build the modified Marblerun, you will need a Go environment set up. If you use Visual Studio Code, you can use the [Dev Container](https://code.visualstudio.com/docs/remote/containers) extension to set up the environment which is defined in the [.devcontainer](.devcontainer) folder. It requires Docker to be installed.

We also need to install [Edgeless RT](https://github.com/edgelesssys/edgelessrt):

```shell
sudo mkdir -p /etc/apt/keyrings
wget -qO- https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | sudo tee /etc/apt/keyrings/intel-sgx-keyring.asc > /dev/null
echo "deb [signed-by=/etc/apt/keyrings/intel-sgx-keyring.asc arch=amd64] https://download.01.org/intel-sgx/sgx_repo/ubuntu $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/intel-sgx.list
sudo apt update
ERT_DEB=edgelessrt_0.4.1_amd64_ubuntu-$(lsb_release -rs).deb
wget https://github.com/edgelesssys/edgelessrt/releases/download/v0.4.1/$ERT_DEB
sudo apt install ./$ERT_DEB build-essential cmake libssl-dev
```

## Build

To build the modified Marblerun, we will execute the commands listed in the [BUILD.md](BUILD.md) file. For convience, the relevant commands are listed below:

```shell
. /opt/edgelessrt/share/openenclave/openenclaverc
mkdir build
cd build
cmake ..
make
```

## How to use

Since this is a modified version of Marblerun, the instructions on how to use it are very similar to the original Marblerun. For simplicity, will show the instructions to execute Marblerun in simulation mode.

*Note:* Always run the following command when you open a new terminal instance:

```shell
. /opt/edgelessrt/share/openenclave/openenclaverc
```

First, let's start the coordinator:

```shell
OE_SIMULATION=1 erthost build/coordinator-enclave.signed
```

*Note:* If you need to reset the state of the coordinator, you can delete the `marblerun-coordinator-data` folder and restart the coordinator.

Now open a new terminal instance and start the provider GRPC ping server (this is the main difference from the original Marblerun):

```shell
./build/marblerun server localhost:4433 -c provider_certificate.crt -k provider_private.key -i
```

This server first contacts the coordinator to get its TLS credentials (and quote if the environment uses SGX). It then starts a GRPC server that listens for ping messages from the coordinator.

Now, open the marble-test-config.json file located in the [build](build) folder. Copy the `SignerID` value and paste it into the `SignerID` of the [manifest](manifest.json) file (line 7). This step is only necessary if you built/rebuilt Marblerun.

Once that is done, let's send the manifest file to the coordinator (you will need to open a new terminal instance for this):

```shell
curl -k --data-binary @manifest.json https://localhost:4433/manifest
```
The pings should now start being sent by the coordinator and received by the provider's server.

Now let's deploy an example application (comes with Marblerun). It consists of a server and a client. Because of the way the deactivation port of each marble is implemented (it is hardcoded to port 50060), you won't be able to start both at the same time (this isn't an issue in Kubernetes, but it is when running locally). So let's start only the server for demonstration purposes:

```shell
. /opt/edgelessrt/share/openenclave/openenclaverc
EDG_MARBLE_TYPE=server EDG_MARBLE_UUID_FILE=$PWD/server_uuid EDG_TEST_ADDR=localhost:8001 OE_SIMULATION=1 erthost build/marble-test-enclave.signed
```

The server is now pinging the coordinator and is ready to receive requests from the client (which we won't be able to start because of the hardcoded port).

## How to see revocation in action

After you execute all the steps in the previous sections, we can now see revocation in action. There are three ways to do this:

1. Send a revocation message to the coordinator. The coordinator will then send a deactivation message to the application modules and shut down.

1. Stop the provider's server. The coordinator pings will timeout and a deactivation message is sent to all the application modules. The coordinator then shuts down.

1. Stop the coordinator. The application modules will timeout and shut down.

For the first option, you will need to open a new terminal instance and execute the following command:

```shell
./build/marblerun deactivate localhost:4433 -c admin_certificate.crt -k admin_private.key -i
```

Both the coordinator and application modules should shutdown after a few seconds.

## Certificates

All communication is secured using TLS. I now explain what each certificate is used for. All certificates are already generated and are located in this folder.

### Provider certificate

This certificate is added to the Marblerun manifest file so the coordinator can verify that it receives ping responses only from the provider. A provider private key is also needed.

### Admin certificate

The manifest file contains a list of admin certificates. These are the certificates of entities that are authorized to send revocation messages to the coordinator. Therefore, these credentials are used when sending revocation requests.

### Coordinator certificate

This certificate is automatically generated when building Marblerun (it's in the [build](build) folder). It is one of the key certificates in Marblerun, used for authenticating the coordinator (see their documentation for more details). It is also requested by the provider when it starts its GRPC ping server.

## Additions to manifest file

The manifest file processing code was modified to support the following fields:

* The `Deactivate` action. It can be assigned to roles and allows this entity to send revocation messages to the coordinator. Uses the existing Marblerun users and roles system.

* The `DeactivationSettings` block. Contains two subblocks with the following fields:
  * `Coordinator`: The deactivation settings for the coordinator.
    * `ConnectionUrl`: The URL of the provider's ping server.
    * `ConnectionCertificate`: The provider certificate used in the ping server.
    * `PingInterval`: The nterval between pings sent to the provider's ping server.
        * `value`: Positive integer
        * `unit`: One of the following: `seconds`, `minutes`, `hours`, `days`
    * `Retries`: The number of ping retries before timing out.
    * `Retry interval`: The interval between retries.
        * `value`: Positive integer
        * `unit`: One of the following: `seconds`, `minutes`, `hours`, `days`
  * `Marbles`: The deactivation settings for the application modules/marbles.
    * `PingInterval`: The interval between pings sent to the coordinator.
        * `value`: Positive integer
        * `unit`: One of the following: `seconds`, `minutes`, `hours`, `days`
    * `Retries`: The number of ping retries before timing out.
    * `Retry interval`: The interval between retries.
        * `value`: Positive integer
        * `unit`: One of the following: `seconds`, `minutes`, `hours`, `days`
