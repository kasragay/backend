<div style="display: flex; align-items: flex-end; margin-bottom: 20px">
  <img src="https://api.kasragay.com/v1/assets/logo150x150.png" alt="Kasragay Logo" style="height: 50px; margin-right: 20px;">
  <h1 style="position: relative; top: 16px;">Kasragay Frontend</h1>
</div>

---

## This foundation is inspired by Kasra's gayness

[![Website](https://img.shields.io/badge/Website-kasragay.com-blue.svg)](https://kasragay.com/) [![API](https://img.shields.io/badge/API-api.kasragay.com-green.svg)](https://api.kasragay.com)

[![Telegram](https://img.shields.io/badge/Telegram-kasra__gay-0088cc.svg)](https://t.me/kasra_gay) [![Discord](https://img.shields.io/badge/Discord-PghhrARr-5865F2.svg)](https://discord.gg/PghhrARr)

[![License: CC BY-NC-SA 4.0](https://img.shields.io/badge/License-CC%20BY--NC--SA%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-sa/4.0/)

---

## Docs
[![V1 Swagger Doc](https://img.shields.io/badge/V1-Doc-green)](https://api.kasragay.com/v1)

---

## How to Use?
**Just do** `ssh kasragay.com` **to try it out.**


---

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
NOREPLY_PHONE=+44742692xxxx

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

Date : 2025-07-10 07:12:47

Total : 58 files,  7937 codes, 171 comments, 765 blanks, all 8873 lines

Summary / [Details](.Counter/2025-07-10/details.md)

## Languages
| language | files | code | comment | blank | total |
| :--- | ---: | ---: | ---: | ---: | ---: |
| Go | 41 | 4,861 | 159 | 643 | 5,663 |
| YAML | 3 | 2,324 | 8 | 35 | 2,367 |
| Go Checksum File | 1 | 286 | 0 | 1 | 287 |
| Makefile | 1 | 148 | 0 | 27 | 175 |
| Markdown | 1 | 135 | 0 | 48 | 183 |
| Go Module File | 1 | 95 | 0 | 4 | 99 |
| HTML | 1 | 73 | 4 | 6 | 83 |
| JSON with Comments | 1 | 15 | 0 | 1 | 16 |
| (Unsupported) | 8 | 0 | 0 | 0 | 0 |

## Directories
| path | files | code | comment | blank | total |
| :--- | ---: | ---: | ---: | ---: | ---: |
| . | 58 | 7,937 | 171 | 765 | 8,873 |
| . (Files) | 6 | 873 | 6 | 101 | 980 |
| .github | 1 | 43 | 0 | 7 | 50 |
| .github/workflows | 1 | 43 | 0 | 7 | 50 |
| Dockerfiles | 2 | 0 | 0 | 0 | 0 |
| assets | 6 | 0 | 0 | 0 | 0 |
| cmd | 3 | 235 | 0 | 46 | 281 |
| cmd/gateway | 1 | 59 | 0 | 15 | 74 |
| cmd/settings | 1 | 127 | 0 | 17 | 144 |
| cmd/user | 1 | 49 | 0 | 14 | 63 |
| docs | 1 | 2,087 | 2 | 8 | 2,097 |
| internal | 38 | 4,626 | 159 | 597 | 5,382 |
| internal/clients | 2 | 20 | 0 | 6 | 26 |
| internal/ports | 8 | 762 | 3 | 133 | 898 |
| internal/repository | 5 | 687 | 0 | 66 | 753 |
| internal/server | 12 | 1,311 | 130 | 183 | 1,624 |
| internal/server (Files) | 4 | 279 | 0 | 35 | 314 |
| internal/server/gateway | 5 | 911 | 130 | 129 | 1,170 |
| internal/server/user | 3 | 121 | 0 | 19 | 140 |
| internal/services | 6 | 1,212 | 9 | 122 | 1,343 |
| internal/utils | 5 | 634 | 17 | 87 | 738 |
| templates | 1 | 73 | 4 | 6 | 83 |

Summary / [Details](.Counter/2025-07-10/details.md)
