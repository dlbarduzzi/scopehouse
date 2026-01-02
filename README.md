# ScopeHouse

<p>
  <a
    href="https://github.com/dlbarduzzi/scopehouse/actions/workflows/build.yml"
    target="_blank"
    rel="noopener">
    <img
      src="https://github.com/dlbarduzzi/scopehouse/actions/workflows/build.yml/badge.svg"
      alt="build"
    />
  </a>
</p>

A centralized control plane for managing and synchronizing Prometheus alerts across multiple clusters.

## Docker

### Local Development

You can run the application locally using Docker.

#### Run the application

From the project root, execute:

```sh
docker compose -f docker-compose.yaml up
```

If you make changes to the source code or Dockerfile, rebuild the image by adding the --build flag:

```sh
docker compose -f docker-compose.yaml up --build
```

## Acknowledgements

This project is heavily inspired by the open-source project
[PocketBase](https://github.com/pocketbase/pocketbase).

Several design patterns and core features are adapted from that project,
with modifications and extensions to meet the goals of this application.

## License

[MIT](./LICENSE)
