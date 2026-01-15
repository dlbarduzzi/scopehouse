# ScopeHouse

A centralized control plane for managing and synchronizing Prometheus alerts across multiple clusters.

## Getting started

### Local development

Export required environment variables.

```sh
# Make sure these values match your database details.
export SH_DATABASE_URL='postgresql://user:pass@127.0.0.1:5432/scopehouse?sslmode=disable'
```

Running database as a docker container.

```sh
docker compose -f docker/compose.local.db.yml up -d
```

## License

[MIT](./LICENSE)
