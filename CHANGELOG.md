## [2.1.1](https://github.com/founda/aws-rds-authenticator/compare/v2.1.0...v2.1.1) (2024-03-20)


### Dependencies and Other Build Updates

* **deps:** bump github.com/aws/aws-sdk-go-v2 from 1.18.1 to 1.21.1 ([2a3e866](https://github.com/founda/aws-rds-authenticator/commit/2a3e866e647e1d45840829156689677b79163bba))

## [2.1.0](https://github.com/founda/aws-rds-authenticator/compare/v2.0.4...v2.1.0) (2023-07-25)


### Features

* add the ability to skip token generation ([d742e18](https://github.com/founda/aws-rds-authenticator/commit/d742e18b8306884c0b9b4ca6ce55492e490e8ee9))

## [2.0.4](https://github.com/founda/aws-rds-authenticator/compare/v2.0.3...v2.0.4) (2023-07-05)


### Bug Fixes

* use PR number as tag for Dependabot-trigger branch builds ([c298f6c](https://github.com/founda/aws-rds-authenticator/commit/c298f6c9bac2cfd74464c7c3136b5992ed5d619a))

## [2.0.3](https://github.com/founda/aws-rds-authenticator/compare/v2.0.2...v2.0.3) (2023-07-04)


### Dependencies and Other Build Updates

* **deps:** bump github.com/aws/aws-sdk-go-v2/config ([97e8510](https://github.com/founda/aws-rds-authenticator/commit/97e8510a9b67061520d95c62436af7f5f592fe10))
* **deps:** bump github.com/aws/aws-sdk-go-v2/feature/rds/auth ([9a8e8f3](https://github.com/founda/aws-rds-authenticator/commit/9a8e8f3284d50810e7d38634ea4b9dba85f39f34))

## [2.0.2](https://github.com/founda/aws-rds-authenticator/compare/v2.0.1...v2.0.2) (2023-07-04)


### Dependencies and Other Build Updates

* **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.4 ([9f69a29](https://github.com/founda/aws-rds-authenticator/commit/9f69a29a6dbdc4ce3835955fde03ab190ef63644))

## [2.0.1](https://github.com/founda/aws-rds-authenticator/compare/v2.0.0...v2.0.1) (2023-05-31)


### Bug Fixes

* update module path to match major release ([ce9c8d2](https://github.com/founda/aws-rds-authenticator/commit/ce9c8d2e2c53e1dd39bb5eee17ed11a34e6add45))

## [2.0.0](https://github.com/founda/aws-rds-authenticator/compare/v1.1.2...v2.0.0) (2023-05-26)


### ⚠ BREAKING CHANGES

* Instead of a string of key/value pairs we now return a connection string.

### Features

* **pgclient:** add support for verify-ca ([93a44b4](https://github.com/founda/aws-rds-authenticator/commit/93a44b439328dd8735dd8051c0ab8d4862a55b86))


### Bug Fixes

* connect to the `postgres` database by default ([4c7657d](https://github.com/founda/aws-rds-authenticator/commit/4c7657dd7046932abb84058f0dcc3c9c70685a57))
* prevent args being passed as single string ([1ef35d8](https://github.com/founda/aws-rds-authenticator/commit/1ef35d84ed45d1a998b404b65e27c47f323c512b))
* remove committed binary ([d5f5654](https://github.com/founda/aws-rds-authenticator/commit/d5f565498eac6f1790a8cc44c8b59b5ce637bfd2))
* return a connection string instead of the key/value pairs ([6ccf0c8](https://github.com/founda/aws-rds-authenticator/commit/6ccf0c8fbf4564f302e81dc2a1a0a1e3ec8bb671))


### Dependencies and Other Build Updates

* **readme:** update README ([f4a1f36](https://github.com/founda/aws-rds-authenticator/commit/f4a1f3699aca21f6200f15ac5c76bdd480a7362d))

## [1.1.2](https://github.com/founda/aws-rds-authenticator/compare/v1.1.1...v1.1.2) (2023-05-02)


### Bug Fixes

* **pgclient:** add missing gexec ([b1d9bb9](https://github.com/founda/aws-rds-authenticator/commit/b1d9bb9945496f3de7b666b1de0e86af7bf720df))

## [1.1.1](https://github.com/founda/aws-rds-authenticator/compare/v1.1.0...v1.1.1) (2023-04-24)


### Bug Fixes

* include Dependabot dependency updates in semantic-release ([42cffec](https://github.com/founda/aws-rds-authenticator/commit/42cffec26876b3603fdb12bbaeb9479dcd35adeb)), closes [/mattbun.io/posts/semantic-release-dependabot/#approach-2](https://github.com/founda//mattbun.io/posts/semantic-release-dependabot//issues/approach-2)
* return dummy version in Dependabot PRs ([0b26170](https://github.com/founda/aws-rds-authenticator/commit/0b26170da6e30be21c3de8fea462fd82c95025dd))

# [1.1.0](https://github.com/founda/aws-rds-authenticator/compare/v1.0.0...v1.1.0) (2023-04-22)


### Features

* enable Dependabot ([d40ddd0](https://github.com/founda/aws-rds-authenticator/commit/d40ddd00bc7885a2c9ba43fbce953b274473711a)), closes [/github.com/dependabot/dependabot-core/issues/3253#issuecomment-852541544](https://github.com//github.com/dependabot/dependabot-core/issues/3253/issues/issuecomment-852541544)

# 1.0.0 (2023-04-22)


### Bug Fixes

* database is not required for authenticator ([50a5caa](https://github.com/founda/aws-rds-authenticator/commit/50a5caade949d0fecc2e2471c3802bb3a5802188))


### Features

* add Dockerfiles for MySQL and Postgres clients ([7691131](https://github.com/founda/aws-rds-authenticator/commit/7691131a3bcd218d2b1f554b44f0d47ce1693152))
* add semantic-release ([13ca0b6](https://github.com/founda/aws-rds-authenticator/commit/13ca0b6e432c7346f991667593cc2ce356e63277))
* **ci:** add CI ([beb89b9](https://github.com/founda/aws-rds-authenticator/commit/beb89b97a41dfb5b27a9374367b58faf23da79ab)), closes [/github.com/orgs/community/discussions/25725#discussioncomment-3248924](https://github.com//github.com/orgs/community/discussions/25725/issues/discussioncomment-3248924)
* **ci:** add client images to CI ([277892f](https://github.com/founda/aws-rds-authenticator/commit/277892f43d395bc85183525fdc95cc554244ce74))
* configurable ssl ([4e7af92](https://github.com/founda/aws-rds-authenticator/commit/4e7af92c1bdeb2937aa5dc3be3ce841a6cb5020b))
