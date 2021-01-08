# casedhttp

casedhttp is a package used to add auditing capabilities to the `net/http` package.

By using casedhttp's provided `casedhttp.ContextMiddleware` around your `http.HandlerFunc` each audit event published using `cased.PublishWithContext` will include the following properties:

| Key                 | Example                              | Sensitive Value |
| ------------------- | ------------------------------------ | --------------- |
| location            | 1.1.1.1                              | ip-address      |
| request_url         | /login                               | -               |
| request_http_method | POST                                 | -               |
| request_user_agent  | Mozilla/5.0                          | -               |
| request_id          | 1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed | -               |

## Usage

See [example](/example/http/main.go) for an example implementation.
