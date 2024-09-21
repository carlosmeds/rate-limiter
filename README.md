# Rate Limiter Application

This application is a rate limiting service that uses Redis to store and manage request limits. The application is built in Go and can be run using Docker Compose.

The rate limit is applied when the configured number of requests is reached, returning status code 429. There is a default limit that indicates how many requests an IP can make. However, if a valid API token is provided, the limit is based on the token's configuration. If the token is invalid, status code 401 is returned.

## Prerequisites

- Docker
- Docker Compose

## Configuration

### Environment Variables

The application uses environment variables to configure Redis and other options. An example file is provided. You should copy this file and adjust the values as necessary.

Rename the file [`.env.example`](./.env.example) to `.env`:

```bash
cp .env.example .env
```
`BLOCKED_TIME`

Sets the block time (in seconds) after a client exceeds the request limit

`DEFAULT_LIMIT`

Default number of requests allowed per IP before being blocked

`API_KEYS`

Specifies rate limits for certain API keys. It is a key-value config. For example, in *your_api_key_value:2,another_api_key:5*, your_api_key_value allows 2 requests, and another_api_key allows 5.

`WEB_SERVER_PORT`

Port where the web server will run, set to 8080 in this case.

## How to Run the Application

1. **Clone o repositório:**

```sh
cd your-repository
git clone https://github.com/carlosmeds/rate-limiter.git
```

2. **Start the containers:**

```sh
docker-compose up --build
```

### Making a request

Run the command in the [`ip.http`](./api/ip.http) file.