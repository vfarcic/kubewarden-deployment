# Kubewarden Policies for dot-sql Crossplane Compositions

## Test

```bash
make test
```

## Publish

```bash
kwctl annotate policy.wasm \
    --metadata-path metadata.yml \
    --output-path annotated-policy.wasm

kwctl push annotated-policy.wasm \
    c8n.io/vfarcic/kubewarden-deployment:v0.0.2
```