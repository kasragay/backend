# Kasragay API

[![License: CC BY-NC-SA 4.0](https://img.shields.io/badge/License-CC%20BY--NC--SA%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-sa/4.0/)


## Commands

---

```bash
# lint
make lint

# run application
APP="all" LONG_VERSION="v1.0.0" make run

# docker build & push
APP="all" LONG_VERSION="v1.0.0" make docker-bp

# docker up (all apps)
make docker-up

# docker down (all apps)
make docker-down

# docker down and volume delete (all apps)
make docker-downv

# docker down
APP="gateway" LONG_VERSION="v1.0.0" make docker-down-app

# docker down and volume delete
APP="gateway" LONG_VERSION="v1.0.0" make docker-downv-app

# create superuser 
make createsuperuser

# delete superuser
make deletesuperuser

# default of:
#   - APP is all
#   - LONG_VERSION is v1.0.0
M="Commit message" APP="gateway" LONG_VERSION="v1.0.0" make compush 
M="Commit message" APP="gateway" LONG_VERSION="v1.0.0" make commit

make push
```

### Environment Variables

---

```env
LONG_VERSION=v1.0.0
VERSION=v1

DEBUG=<bool>
DOMAIN=kasragay.com
PORT=<port>


POSTGRES_DB_HOST=<url>
POSTGRES_DB_PORT=<port>
POSTGRES_DB_USERNAME=<string>
POSTGRES_DB_PASSWORD=<string>
POSTGRES_DB_DATABASE=<string>

MONGO_HOST=<string>
MONGO_PORT=<port>
    
DRAGONFLYDB_HOST=<url>
DRAGONFLYDB_PORT=<port>
DRAGONFLYDB_PASSWORD=<string>

MINIO_USERNAME=<string>
MINIO_PASSWORD=<string>
MINIO_HOST=<url>
MINIO_PORT=<port>
MINIO_USE_SSL=<bool>
MINIO_AVATARS_BUCKET=<string>

TWILIO_ACCOUNT_SID=<string>
TWILIO_AUTH_TOKEN=<string>
NOREPLY_PHONE=0742692xxxx

SMTP_HOST=smtpout.secureserver.net
SMTP_PORT=587
NOREPLY_EMAIL=no-reply@kasragay.com
NOREPLY_EMAIL_PASSWORD=<string>
SUPPORT_EMAIL=support@kasragay.com

# Gateway related
GATEWAY_PORT=<port>

TLS_ON=true
TLS_CRT_FILE=/etc/ssl/crt.pem
TLS_KEY_FILE=/etc/ssl/key.pem

ADDITIONAL_ALLOWED_HOSTS=<comma-sep-string>
ADDITIONAL_ALLOWED_ORIGINS=<comma-sep-string>

JWT_SECRET_KEY=<string>
JWT_ACCESS_EXP=1200
JWT_REFRESH_EXP=14400

PASSWORD_HASH_COST=<int>
PASSWORD_HASH_SALT=<string>
USER_BACK_MAXIMUM_REFERENCE=<int>
#

# User related
USER_PORT=<port>
#

# docker-compose related
MINIO_CONSOLE_PORT=<port>
#
```

### Count Progress

---

# Summary

Date : 2025-07-10 04:58:49

Total : 58 files,  7931 codes, 170 comments, 762 blanks, all 8863 lines

Summary / [Details](.Counter/2025-07-10/details.md)

## Languages
| language | files | code | comment | blank | total |
| :--- | ---: | ---: | ---: | ---: | ---: |
| Go | 41 | 4,865 | 158 | 645 | 5,668 |
| YAML | 3 | 2,324 | 8 | 35 | 2,367 |
| Go Checksum File | 1 | 286 | 0 | 1 | 287 |
| Makefile | 1 | 148 | 0 | 27 | 175 |
| Markdown | 1 | 125 | 0 | 43 | 168 |
| Go Module File | 1 | 95 | 0 | 4 | 99 |
| HTML | 1 | 73 | 4 | 6 | 83 |
| JSON with Comments | 1 | 15 | 0 | 1 | 16 |
| (Unsupported) | 8 | 0 | 0 | 0 | 0 |

## Directories
| path | files | code | comment | blank | total |
| :--- | ---: | ---: | ---: | ---: | ---: |
| . | 58 | 7,931 | 170 | 762 | 8,863 |
| . (Files) | 6 | 863 | 6 | 96 | 965 |
| .github | 1 | 43 | 0 | 7 | 50 |
| .github/workflows | 1 | 43 | 0 | 7 | 50 |
| Dockerfiles | 2 | 0 | 0 | 0 | 0 |
| assets | 6 | 0 | 0 | 0 | 0 |
| cmd | 3 | 235 | 0 | 46 | 281 |
| cmd/gateway | 1 | 59 | 0 | 15 | 74 |
| cmd/settings | 1 | 127 | 0 | 17 | 144 |
| cmd/user | 1 | 49 | 0 | 14 | 63 |
| docs | 1 | 2,087 | 2 | 8 | 2,097 |
| internal | 38 | 4,630 | 158 | 599 | 5,387 |
| internal/clients | 2 | 20 | 0 | 6 | 26 |
| internal/ports | 8 | 765 | 3 | 133 | 901 |
| internal/repository | 5 | 687 | 0 | 66 | 753 |
| internal/server | 12 | 1,311 | 130 | 183 | 1,624 |
| internal/server (Files) | 4 | 279 | 0 | 35 | 314 |
| internal/server/gateway | 5 | 911 | 130 | 129 | 1,170 |
| internal/server/user | 3 | 121 | 0 | 19 | 140 |
| internal/services | 6 | 1,213 | 8 | 124 | 1,345 |
| internal/utils | 5 | 634 | 17 | 87 | 738 |
| templates | 1 | 73 | 4 | 6 | 83 |

Summary / [Details](.Counter/2025-07-10/details.md)
