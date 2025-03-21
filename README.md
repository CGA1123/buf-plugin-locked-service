# `buf-plugin-locked-service`

Implements a "breaking" change detection rule via [`bufplugin-go`] which
restricts addition to an existing RPC service definition.

Useful if you want to restrict growing a large service interface.

Can be configured via:

```yaml
# buf.yaml
version: v2
lint:
  use:
    - LOCKED_SERVICE
plugins:
  - plugin: buf-plugin-locked-service
    options:
      locked_services:
        - your.package.v1.YourLockedService
```

[`bufplugin-go`]: https://github.com/bufbuild/bufplugin-go
