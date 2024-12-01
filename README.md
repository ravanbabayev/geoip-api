# GeoIP API

GeoIP API is a service developed to retrieve geolocation information based on IP addresses. This API uses the MaxMind GeoLite2 database to provide information such as country, city, location coordinates of IP addresses.

## Features

- Provides IP address based geolocation information.
- MaxMind provides data using the GeoLite2-City database.
- Easy installation with Docker support.
- The API reads port and configuration settings from the `.env` file.

## Requirements

- Docker (Optional, project can be run with Docker container)
- Go language (used in the development of the API)

## Installation

#### 1. Installation Using Docker

You can follow the steps below to run the project via Docker:

Clone the project:
   ```bash
   git clone https://github.com/ravanbabayev/geoip-api.git
   cd geoip-api
   docker build -t geoip-api .
   docker run -d --name geoip-api -p 8080:8080 --env-file .env geoip-api
  ```
