# 

<!-- PROJECT LOGO -->
<br />
<div align="center">

<h3 align="center">snet-matrix-framework</h3>

  <p align="center">
    ai-bots based on the Matrix protocol and snet ecosystem
    <br />
    <a href="https://github.com/tensved/snet-matrix-framework"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/tensved/snet-matrix-framework">View Demo</a>
    ·
    <a href="https://github.com/tensved/snet-matrix-framework/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/tensved/snet-matrix-framework/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>

## About The Project

The snet-matrix-framework is designed to create bots that will connect the messenger user on the Matrix protocol and AI services in the snet ecosystem

## Getting Started

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/tensved/snet-matrix-framework.git
   ```
2. Install dependencies
   ```sh
   go mod download
   ```
3. Create an `.env` file based on `example.env` and add data to it
4. Generate `registry_and_mpe.go` file
   ```sh
   go generate ./...
   ```

## Usage

The minimal example is located at the path `pkg/lib/examples/snet/main.go`
