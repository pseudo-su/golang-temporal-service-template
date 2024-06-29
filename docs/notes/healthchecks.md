# Healthchecks

```sh
# frontdoor gRPC readiness check
grpcurl -plaintext localhost:9090 grpc.health.v1.Health/Check
# worker gRPC readiness check
grpcurl -plaintext localhost:9091 grpc.health.v1.Health/Check
```

```sh
# Deephealth checks
curl http://localhost:9090/v1/health/deep | jq
curl --header "Content-Type: application/json" --data '{}' \
  http://localhost:9090/deephealth.v1.DeepHealth/Check
grpcurl -plaintext localhost:9090 deephealth.v1.DeepHealth/Check
```
