# WAF Extension

This extension implements a Web Application Firewall using [OWASP Coraza](https://coraza.io/) and comes with rules from the [OWASP Core Rule Set (CRS)](https://coreruleset.org/) already embedded and ready to use.

## Configuration

The filter accepts a JSON configuration with the following fields:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `directives` | `[]string` | *(required)* | List of Coraza SecLang directives. Each entry is a single directive string. Supports special includes such as `Include @coraza.conf`, `Include @crs-setup.conf`, `Include @owasp_crs/*.conf`. |
| `mode` | `string` | `"REQUEST_ONLY"` | WAF inspection mode. `"REQUEST_ONLY"` processes only request phases (fast path, recommended for most deployments). `"FULL"` processes both request and response phases. `"RESPONSE_ONLY"` processes only response phases. |
| `header_mode` | `string` | `"FULL"` | Controls which request headers are forwarded to Coraza. `"FULL"` forwards all request headers (preserves existing behaviour). `"MINIMAL"` forwards only a security-relevant subset (see below) to reduce per-request allocation and rule-variable population cost. |

### `header_mode: MINIMAL` header allowlist

When `header_mode` is set to `"MINIMAL"`, only the following headers are forwarded to Coraza (in addition to `Host` which is always forwarded):

- `user-agent`
- `accept`
- `content-type`
- `content-length`
- `cookie`
- `authorization`
- `referer`
- `origin`
- `x-forwarded-for`
- `x-real-ip`

All other headers are ignored by the WAF. This reduces per-request allocations and speeds up CRS rule evaluation. Use `"FULL"` if you require rules that inspect arbitrary custom headers.

### GET/HEAD no-body fast path

Requests with method `GET` or `HEAD`, or with an explicit `Content-Length: 0` header, skip body buffering entirely. Phase-2 rules (request body phase) are run immediately within the request-headers callback with an empty body, and the filter returns `Continue` without waiting for a body. This avoids the `HeadersStatusStop` / buffering overhead for the common no-body case.

## Rule Files

The WAF resolves rule files from three layered filesystems (first match wins):

1. **Embedded rules** — tailored configurations shipped with this extension.
2. **Coraza CoreRuleSet package** — upstream OWASP CRS rules.
3. **Local filesystem** — support for user-provided overrides and custom rules.

## Upgrading CRS

The CRS rules are provided by the [Coraza CoreRuleSet package](https://github.com/corazawaf/coraza-coreruleset). To upgrade:

1. Bump the CRS dependency version in `go.mod`. You can find the latest CRS version on the [coraza-coreruleset releases page](https://github.com/corazawaf/coraza-coreruleset/releases):
```sh
go get github.com/corazawaf/coraza-coreruleset/v4@<new-version>
go mod tidy
```

1. Update the embedded CRS configuration file (`@crs-setup.conf`) if there are any relevant changes in the upstream `crs-setup.conf.example`. The embedded config should be updated while preserving any custom configurations specific to this repository.

## Upgrading Coraza

To upgrade to a new Coraza version:

1. Bump the Coraza dependency version in `go.mod`:
   ```sh
   go get github.com/corazawaf/coraza/v3@<new-version>
   go mod tidy
   ```

1. Update the embedded Coraza configuration file (`@coraza.conf`) if there are any relevant changes in the upstream `coraza.conf-recommended`. The embedded config should be updated while preserving any custom configurations specific to this repository.
