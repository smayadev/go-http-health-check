# Go HTTP Health Check

## Description

This is a Go-based HTTP health check script that reads a YAML configuration file containing endpoint data and will make requests to each one to determine, on a per-domain basis, if the endpoint is up and logs the availability percentage.

## Prerequisites

You will need to ensure Go is installed on your system:

[Install Go](https://go.dev/doc/install)

Verify your Go installation by running:

```bash
go version
```

## Installation

Close this repository to your local computer:

```bash
git clone https://github.com/smayadev/go-http-health-check.git
```

Alternatively, you may download a zip of this project to your local computer and open it with the archive application of your choice.

Once cloned or downloaded, navigate into the project directory:

```bash
cd go-http-health-check
```

Then, build the binary:

```bash
go build -o health_check main.go
```

## Usage

Output of `./health_check -h`:

```bash
Usage of ./health_check:
  -f string
        Full path to the YAML containing URL data, i.e. /home/user/url_data.yaml
```

Once the program is built, you may run it as follows:

```bash
./health_check -f /home/user/url_data.yaml
```

Pass the full path of your YAML configuration file to the program using the `-f` flag as shown, replacing `/home/user/url_data.yaml` with the actual path to the file on your computer.

See below for an example YAML configuration file.

## Example Configuration

```yaml
- headers:
    user-agent: test-google-response
  method: GET
  name: Google homepage
  url: https://google.com/
- headers:
    user-agent: test-google-response
  method: GET
  name: Google careers page
  url: https://google.com/careers
- body: '{"foo":"bar"}'
  headers:
    content-type: application/json
    user-agent: test-google-response
  method: POST
  name: fake endpoint
  url: https://google.com/some/post/endpoint
- name: GitHub homepage
  url: https://www.github.com/
```

- `name` (string, required) — A free-text name to describe the HTTP endpoint.
- `url` (string, required) — The URL of the HTTP endpoint.
- `method` (string, optional) — The HTTP method of the endpoint. Default is GET.
- `headers` (dictionary, optional) — The HTTP headers to include in the request.
- `body` (string, optional) — The JSON-encoded HTTP body to include in the request.
