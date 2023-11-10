This folder contains example build and install files:

# Dockerfile

A `Dockerfile` is used to build a docker image that is handy for deployments and local tests.
In order to use it in your publisher project - rename the binary inside it, or provide the binary name in the arguments.

# Makefile

A `Makefile` is used to build, test, update the code in a repeatable manner.
In order to use it in your publisher project - rename the binary name inside it.

# install-service.sh

`install-service.sh` is a shell executable file that helps installing a systemd service for your publisher.
You can provide a publisher name to it that will be also used as a binary name, however, by default, the binary path is still
`/usr/bin`.
