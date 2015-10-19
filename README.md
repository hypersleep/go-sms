# go-sms
SMS notifications using Billingrad API

This application using Consul KV storage for serving configuration.

By default it running in development environment and loads keys from Consul with `go-sms/development` prefix.

You can use any different environment using `_ENV` environment variable. (ex. `_ENV=staging` = `go-sms/staging` prefix in consul)

For production environment ensure this consul keys exist:

`go-sms/production/auth_token` - auth token for go-sms HTTP API

`go-sms/production/close_api_key` - Billingrad's API

`go-sms/production/delivery_id` - Billingrad's API

`go-sms/production/open_api_key` - Billingrad's API

`go-sms/production/port` - port for go-sms HTTP API

`go-sms/production/receivers` - SMS receivers (splitted by `,` comma)

