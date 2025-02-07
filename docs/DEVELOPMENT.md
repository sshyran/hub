# Development

- [Getting Started](#getting-started)
- [Checkout repository](#checkout-your-fork)
- [Requirements](#requirements)
- [API Service](#api-service)
  - [Running Database](#running-database)
  - [Running Database Migration](#database-migration)
  - [Adding GitHub OAuth Configuration](#adding-github-oAuth-configuration)
  - [Running API Service](#running-api-service)
  - [Running Tests](#running-tests)
  - [Creating JWT for testing](#creating-jwt-for-testing)
  - [API Documentation](#api-documentation)
- [UI](#ui)
  - [Running UI](#running-ui)
  - [Running UI Tests](#running-ui-tests)

## Getting Started

1. [Create a GitHub Account][join-github]
1. [Setup GitHub access via SSH][gh-ssh]

## Checkout your fork

To check out this repository:

1. Create your own [fork of this repository][fork-repo]
2. Clone it to your machine:

```shell
  git clone git@github.com:${YOUR_GITHUB_USERNAME}/hub.git
  cd hub

  git remote add upstream git@github.com:tektoncd/hub.git
  # prevent accidental push to upstream
  git remote set-url --push upstream no-push
  git fetch --all
```

Adding the upstream remote sets you up nicely for regularly [syncing your fork][sync-fork].

## Requirements

You must install these tools:

1. [`go`][install-go]: The language hub apis are built in.
1. [`git`][install-git]: For source control
1. [`node`][install-node]: To publish and install packages to and from the public npm registry

You may need to install more tools depending on the way you want to run the hub.

## API Service

### Running database

Two ways to run postgresql database:

- [Install postgresql][install-pg] on your local machine, or
- Run a postgresql container using [docker][install-docker] / [podman][install-podman]

If you have installed postgresql locally, you need to create a `hub` database.

**NOTE:** Use the same configuration mentioned in `.env.dev` or
update `.env.dev` with the configuration you used. The api service
and db migration uses the db configuration from `.env.dev`.

- If you want to run a postgres container, source the `.env.dev` so that
  `docker` can use the same database configuration as in `.env.dev` to create a container.

  Ensure you are in `hub/api` directory.

  ```bash
  source .env.dev

  docker run -d  --name hub \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_DB=$POSTGRES_DB \
    -p $POSTGRES_PORT:5432 \
    postgres
  ```

### Database Migration

Once the database is up and running, you can run migration to create tables.

Run the following command to run migration

```bash
go run ./cmd/db
```

Wait until the migration completes and logs to show

> DB initialisation successful !!

### Adding GitHub OAuth Configuration

Create a GitHub OAuth. You can find the steps to create it [here][gh-oauth] with `Authorization callback URL` as `http://localhost:4200` \
Create a Gitlab OAuth. You can find the steps to create it [here][gl-oauth] with `Authorization callback URL` as `http://localhost:4200/auth/gitlab/callback` \
Create a BitBucket OAuth. You can find the steps to create it [here][bb-oauth] with `Authorization callback URL` as `http://localhost:4200`



After creation, add the OAuth Client ID as \
OAuth Client ID `GH_CLIENT_ID` and Client Secret as `GH_CLIENT_SECRET` for Github \
OAuth Client ID `GL_CLIENT_ID` and Client Secret as `GL_CLIENT_SECRET` for Gitlab \
OAuth Client ID `BB_CLIENT_ID` and Client Secret as `BB_CLIENT_SECRET` for BitBucket \
in [.env.dev][env-dev].

For `JWT_SIGNING_KEY`, you can use any random word.

For `ACCESS_JWT_EXPIRES_IN` and `REFRESH_JWT_EXPIRES_IN` you can set the time as per your convenience. Example `15m`, `10y`.

### Running API Service

Once the database is setup and the migration has been run, you can run api service by

```bash
go run ./cmd/api
```

### Running tests

To run the tests, we need a test db.

- If you have installed postgresql, create a `hub_test` database.
- If you are running a container, create `hub_test` database in the same container.

```bash
  source .env.dev

  docker exec -it hub bash -c \
    "PGPASSWORD=$POSTGRES_PASSWORD \
     psql -h localhost -p 5432 -U postgres -c 'create database hub_test;'"
```

Once the `hub_test` database is created, you can run the test using following command:

```bash
  go test -p 1 -count=1 -v ./pkg/...
```

To re-generate the golden files use the below command

```bash
go test $(go list -f '{{ .ImportPath }} {{ .TestImports }}' ./... | grep gotest.tools/v3/golden | awk '{print $1}' | tr '\n' ' ') -test.update-golden=true
```

This will run `go test a/package -test.update-golden=true` on all packages that are importing `gotest.tools/v3/golden`

**NOTE:** `tests` use the database configurations from [test/config/env.test][env-test-file]

### Creating JWT for testing

To create a JWT, Open below URL in a browser.

```
https://github.com/login/oauth/authorize?client_id=<Add Client ID here>
```

Add your OAuth Client ID from [.env.dev][env-dev] in place of `<Add Client ID here>`.

It will redirect you to GitHub. Login using your GitHub Credentials and Once you authorize, it will redirect you to `localhost:8080`.

for ex. `http://localhost:8080/?code=32d4a0b4eb6e9fbea731`

Use the `code` from url in `/auth/login` API. It will add you as a user in db and return a JWT.

#### JWT with Additional Scopes

By default, the JWT has only defaut scopes. If you need additional scopes in your JWT then add your GitHub username with required scopes in your local [config.yaml][config-yaml]

And repeat the login process and you will have additonal scopes in JWT.

### API Documentation

API documentation is generated by goa in file [`gen/http/openapi.yaml`][swagger-doc].

Also, you can call the API `/swagger/swagger.json` to get the documentaion.

You can paste file content or API response in [swagger client][swagger].

## UI

### Running UI

Ensure you are in `hub/ui` directory

Run the following command to install the dependencies

```
npm install
```

To start the application run the following command

```
npm start
```

You will see Hub running at `http://localhost:3000`

### Running UI Tests

Run the following command to run all the tests

```
npm test
```

[join-github]: https://github.com/join
[gh-ssh]: https://help.github.com/articles/connecting-to-github-with-ssh/
[fork-repo]: https://help.github.com/articles/fork-a-repo/
[sync-fork]: https://help.github.com/articles/syncing-a-fork/
[install-go]: https://golang.org/doc/install
[install-goa]: https://github.com/goadesign/goa
[install-git]: https://help.github.com/articles/set-up-git/
[install-pg]: https://www.postgresql.org/docs/12/tutorial-install.html
[install-docker]: https://docs.docker.com/engine/install/
[install-podman]: https://podman.io/getting-started/installation.html
[env-dev]: https://github.com/tektoncd/hub/blob/master/api/.env.dev
[env-test-file]: https://github.com/tektoncd/hub/blob/master/api/test/config/env.test
[gh-oauth]: https://docs.github.com/en/developers/apps/creating-an-oauth-app
[gl-oauth]: https://docs.gitlab.com/ee/integration/oauth_provider.html#user-owned-applications
[bb-oauth]: https://support.atlassian.com/bitbucket-cloud/docs/use-oauth-on-bitbucket-cloud
[config-yaml]: https://github.com/tektoncd/hub/blob/master/config.yaml
[swagger]: https://editor.swagger.io
[swagger-doc]: https://github.com/tektoncd/hub/blob/master/api/gen/http/openapi.yaml
[install-node]: https://nodejs.org/en/download
