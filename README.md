# touyakun

## Overview

LINE DC BOT AWARDS 2024 提出作品
投薬くん(とうやくん)

## Requirement

### OS

- Mac OS Sonoma 14.4(動作確認済み)

### Library

- Go
  - Gin
- Python
  - FastAPI
- NGINX  
- Docker
- docker-compose

## Installation(local)

1. Clone this repository

```
git clone git@github.com:git@github.com:GoRuGoo/fib_api.git
```
2. Build

```
docker compose up -d
```
3. Adding a self-signed certificate
```
cd reverseproxy
```
```
mkdir key
```
```
cd key
```
```
brew install mkcert
```
```
mkcert -install
```
```
mkcert touyakun.com
```
```
mkcert ai.touyakun.com
```

4. Add to /etc/hosts(Not necessary if you use cURL's "--resolve" option when sending a request)
```
127.0.1 touyakun.com
```


## Usage(local)

1. Build & start container

```
docker compose up -d
```

2. Access
```
curl -X GET -H "Content-Type application/json" "https://touyakun.com"
```




## Author

- [Yuta Ito](https://github.com/GoRuGoo)
