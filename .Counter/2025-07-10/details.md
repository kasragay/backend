# Details

Date : 2025-07-10 04:58:49

Directory /home/yuhaha/kasragay/backend

Total : 58 files,  7931 codes, 170 comments, 762 blanks, all 8863 lines

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [.github/workflows/ci-cd.yml](/.github/workflows/ci-cd.yml) | YAML | 43 | 0 | 7 | 50 |
| [Dockerfiles/gateway](/Dockerfiles/gateway) | (Unsupported) | 0 | 0 | 0 | 0 |
| [Dockerfiles/user](/Dockerfiles/user) | (Unsupported) | 0 | 0 | 0 | 0 |
| [Makefile](/Makefile) | Makefile | 148 | 0 | 27 | 175 |
| [README.md](/README.md) | Markdown | 125 | 0 | 43 | 168 |
| [assets/favicon.ico](/assets/favicon.ico) | (Unsupported) | 0 | 0 | 0 | 0 |
| [assets/favicon16x16.ico](/assets/favicon16x16.ico) | (Unsupported) | 0 | 0 | 0 | 0 |
| [assets/favicon32x32.ico](/assets/favicon32x32.ico) | (Unsupported) | 0 | 0 | 0 | 0 |
| [assets/logo150x150.png](/assets/logo150x150.png) | (Unsupported) | 0 | 0 | 0 | 0 |
| [assets/logo250x250.png](/assets/logo250x250.png) | (Unsupported) | 0 | 0 | 0 | 0 |
| [assets/logo400x400.png](/assets/logo400x400.png) | (Unsupported) | 0 | 0 | 0 | 0 |
| [cmd/gateway/main.go](/cmd/gateway/main.go) | Go | 59 | 0 | 15 | 74 |
| [cmd/settings/main.go](/cmd/settings/main.go) | Go | 127 | 0 | 17 | 144 |
| [cmd/user/main.go](/cmd/user/main.go) | Go | 49 | 0 | 14 | 63 |
| [docker-compose.yml](/docker-compose.yml) | YAML | 194 | 6 | 20 | 220 |
| [docs/swagger.yaml](/docs/swagger.yaml) | YAML | 2,087 | 2 | 8 | 2,097 |
| [go.mod](/go.mod) | Go Module File | 95 | 0 | 4 | 99 |
| [go.sum](/go.sum) | Go Checksum File | 286 | 0 | 1 | 287 |
| [internal/clients/consts.go](/internal/clients/consts.go) | Go | 5 | 0 | 2 | 7 |
| [internal/clients/shared.go](/internal/clients/shared.go) | Go | 15 | 0 | 4 | 19 |
| [internal/ports/abc.repositories.go](/internal/ports/abc.repositories.go) | Go | 47 | 0 | 12 | 59 |
| [internal/ports/abc.server.go](/internal/ports/abc.server.go) | Go | 35 | 1 | 9 | 45 |
| [internal/ports/abc.services.go](/internal/ports/abc.services.go) | Go | 50 | 1 | 12 | 63 |
| [internal/ports/consts.go](/internal/ports/consts.go) | Go | 4 | 0 | 2 | 6 |
| [internal/ports/io.auth.go](/internal/ports/io.auth.go) | Go | 126 | 0 | 28 | 154 |
| [internal/ports/io.user.go](/internal/ports/io.user.go) | Go | 30 | 0 | 7 | 37 |
| [internal/ports/models.user.go](/internal/ports/models.user.go) | Go | 201 | 0 | 31 | 232 |
| [internal/ports/validations.go](/internal/ports/validations.go) | Go | 272 | 1 | 32 | 305 |
| [internal/repository/cache.go](/internal/repository/cache.go) | Go | 107 | 0 | 16 | 123 |
| [internal/repository/consts.go](/internal/repository/consts.go) | Go | 4 | 0 | 2 | 6 |
| [internal/repository/mongo.go](/internal/repository/mongo.go) | Go | 46 | 0 | 9 | 55 |
| [internal/repository/relational.go](/internal/repository/relational.go) | Go | 415 | 0 | 28 | 443 |
| [internal/repository/s3.go](/internal/repository/s3.go) | Go | 115 | 0 | 11 | 126 |
| [internal/server/consts.go](/internal/server/consts.go) | Go | 4 | 0 | 2 | 6 |
| [internal/server/gateway/consts.go](/internal/server/gateway/consts.go) | Go | 106 | 0 | 6 | 112 |
| [internal/server/gateway/middlewares.go](/internal/server/gateway/middlewares.go) | Go | 88 | 0 | 6 | 94 |
| [internal/server/gateway/routes.go](/internal/server/gateway/routes.go) | Go | 477 | 0 | 29 | 506 |
| [internal/server/gateway/server.go](/internal/server/gateway/server.go) | Go | 28 | 0 | 5 | 33 |
| [internal/server/gateway/swagger.go](/internal/server/gateway/swagger.go) | Go | 212 | 130 | 83 | 425 |
| [internal/server/middlewares.go](/internal/server/middlewares.go) | Go | 98 | 0 | 6 | 104 |
| [internal/server/routes.go](/internal/server/routes.go) | Go | 24 | 0 | 5 | 29 |
| [internal/server/server.go](/internal/server/server.go) | Go | 153 | 0 | 22 | 175 |
| [internal/server/user/consts.go](/internal/server/user/consts.go) | Go | 4 | 0 | 2 | 6 |
| [internal/server/user/routes.go](/internal/server/user/routes.go) | Go | 93 | 0 | 12 | 105 |
| [internal/server/user/server.go](/internal/server/user/server.go) | Go | 24 | 0 | 5 | 29 |
| [internal/services/auth.go](/internal/services/auth.go) | Go | 564 | 1 | 47 | 612 |
| [internal/services/consts.go](/internal/services/consts.go) | Go | 4 | 0 | 2 | 6 |
| [internal/services/mailcom.go](/internal/services/mailcom.go) | Go | 171 | 6 | 19 | 196 |
| [internal/services/ratelimiter.go](/internal/services/ratelimiter.go) | Go | 262 | 0 | 28 | 290 |
| [internal/services/telecom.go](/internal/services/telecom.go) | Go | 120 | 1 | 18 | 139 |
| [internal/services/user.go](/internal/services/user.go) | Go | 92 | 0 | 10 | 102 |
| [internal/utils/consts.go](/internal/utils/consts.go) | Go | 2 | 0 | 2 | 4 |
| [internal/utils/error.go](/internal/utils/error.go) | Go | 255 | 16 | 30 | 301 |
| [internal/utils/funcs.go](/internal/utils/funcs.go) | Go | 190 | 1 | 26 | 217 |
| [internal/utils/logger.go](/internal/utils/logger.go) | Go | 155 | 0 | 22 | 177 |
| [internal/utils/set.go](/internal/utils/set.go) | Go | 32 | 0 | 7 | 39 |
| [kasragay-backend.code-workspace](/kasragay-backend.code-workspace) | JSON with Comments | 15 | 0 | 1 | 16 |
| [templates/otp\_email.html](/templates/otp_email.html) | HTML | 73 | 4 | 6 | 83 |

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)