# REST Auth part
## Installation
1. Clone the repository
2. Run `docker build -t "build_image" -f .\deploy\Dockerfile .` in console
3. Run `docker-compose build --no-cache`
4. Run `docker-compose up -d` in console.
## Description
This is a simple REST API for user authentication. It has 2 endpoints:
1. GET `/auth/{guid}/` - for user authentication
2. POST `/refresh/` - for token refresh
body: `{"refreshT": "your_refresh_token", "accessT": "your_access_token"}`

