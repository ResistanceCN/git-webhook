# git-webhook

A generic webhook service.

It will execute a sh command when receive a HTTP request.

## Usage

```
GET /{name}
GET /{name}?password={password}
POST /{name}
POST /{name}?password={password}
```

Only ASCII characters and numbers could be contained in `name`.
