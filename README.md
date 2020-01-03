# Maimok

A web interface for virtual machines running on KVM/libvirt written in Go and Vue.js

## Features

- Create virtual machines from base images supporting cloud-init (currently, only Ubuntu based images are tested and supported)
- Configure RAM and disk space on creation
- Setup network inside VM: IP address, netmask, gateway
- Create user `sysadmin` with sudo permissions (without password)
- Inject ssh keys into VM (`sysadmin` user)
- Set hostname of VM
- Install all updates on first bootup
- Install all apt upgrades and reboot VM at 04:00 if necessary (unattended-upgrades)

Tested with base images from https://cloud-images.ubuntu.com/.

## Setup

- On the virtual machine host, install KVM and libvirtd.
- Download a base image (e.g. https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img) and move it to `/var/lib/libvirt/images/`
- Clone this repo to the machine that should run maimok (can be the virtual machine host)
- Install docker and docker-compose
- Build docker image: `docker build . -t maimok`
- Copy and modify one of the docker-compose example files in the `examples/` directories
- Start the docker container with `docker-compose up -d`

## Development

### Backend

The REST API backend is written in Go using the chi router.
To start hacking on the backend:

- Install Go
- Create a `config.toml` file from the provided `example-config.toml` file
- Run the backend with `go run main.go`

### Frontend

The web frontend is a TypeScript SPA using Vue.js and the Vuetify.js component frameowrk.

- Install a current version of node.js, npm and yarn
- Install dependencies: `yarn`
- Run the development server: `yarn serve`
