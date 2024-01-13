# Manually build containers

```sh
# build container and tag it as "worker"
podman build \
  -t worker \
  --build-arg entrypoint=modules/worker/cmd/worker \
  .

# create container from the local image called "worker"
podman run --net=host -t worker
```
