# dyngo
Dynamic DNS on your own INWX Domain, written in Go.
It will determine the external IP of the system it is running on and updates a given domain record.


## Quickstart
To get up and running ASAP, follow these steps:
1. `cp dyngo.env.example dyngo.env`
2. Enter your credentials and target domain in `dnygo.env`
3. `docker-compose build`
4. `docker-compose up -d`


## Running
The recommended way of running `dyngo` is via **docker**.
`dyngo` is configured completely through environment variables.
Please ensure the following variables are set before executing the binary

Variable | Example | Description
---------|----------|-------------
INWX_USERNAME | abcd1234 | The username you use for logging in at INWX
INWX_PASSWORD | your-pass | The password you use for logging in at INWX
INWX_DOMAIN_RECORD | dyn.yourdomain.com | The *full domain name* you want to update
INWX_SLEEP_MINUTES | 10 | Time dyngo waits between domain ip checks


## Docker
You can build a current snapshot of `dyngo` by running a `docker build` in the repo root:
`docker build --tag dyngo .`
Afterwards, you can start it by running
`docker run --rm --env-file=/path/to/dyngo.env -t dyngo`

It is recommended that you provide the needed environment variables in a [separate env file](https://docs.docker.com/engine/reference/commandline/run/#set-environment-variables--e---env---env-file) for readability.
