# Golang temporal service template

## Overview

## Contributing quickstart

### Global dependencies

* Go `1.21`
* Sdkman to manage Java versions

```sh
# On macs install coreutils
brew install coreutils
brew install graphviz

# On macs go can be installed through homebrew
brew install go@1.21

# If using sdkman to manage Java versions
curl -s "https://get.sdkman.io" | bash
sdk env install
sdk env
```

### Local config files

```sh
# Copy the example .env.local file (used to override environment variables for local development)
cp docs/examples/.env.local .

# Copy the recommended vscode settings to your workspace config
cp .vscode/settings.recommended.json .vscode/settings.json
```

If you are using `direnv` (recommended), copy the example `.envrc` file.

```sh
cp docs/examples/.envrc .
direnv allow
direnv reload
```

### Run development scripts

```sh
# Make sure the correct version of language tooling is active before running any commands
sdk env use;

# show available makefile targets
make help;

# Install project dependencies (installs dependencies and tools)
make deps.install;

# Fetch devstack containers
make devstack.pull;

# Run code verification (static analysis, linting etc)
make verify;

# Verify code using static analysis tools and automatically apply fixes when possible
make verify.fix;

# Run all code generation
make codegen;

# Run unit tests
make test.unit;

# Start local devstack dependencies (Postgres and PgAdmin)
make devstack.start;

# Run tests (some tests rely on having the DB running)
make test.integration.blackbox;
```

At this point you should have a development version of this API project running 🎉

You can stop or recreate the devstack using the following commands

```sh
# Stop shutdown the containers running as part of the devstack
make devstack.stop

# Delete/reset the devstack, removes all the containers, volumes etc of the compose stack
make devstack.clean
```

## Next steps

* TBD
