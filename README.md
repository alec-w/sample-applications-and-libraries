# Sample Applications and Libraries

Collection of sample applications and libraries that are useful across a variety of projects to serve as initial bases.

Start the development environment with:
```bash
# Run on host
make dev-env/up
```

Get a shell in the dev container with:
```bash
# Run on host
make dev-env/shell
```

Connect to the postgres database (from the dev container) with:
```bash
# Run in dev container
make database/connect
```

Stop the development environment with:
```bash
# Run on host
make dev-env/down
```

## Applications

### REST API

Setup the REST API's database schema with:
```bash
make rest-api/database/up
```

Populate the REST API's database with:
```bash
make rest-api/database/populate
```

Truncate th REST API's database with:
```bash
make rest-api/database/truncate
```

Remove the REST API's database schema with:
```bash
make rest-api/database/down
```
