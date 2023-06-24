# Getting started

## Server side

Start by [building and installing
`accio127`](https://git.sr.ht/~jamesponddotco/accio127#installation),
and then copy [the example `config.json`
file](https://git.sr.ht/~jamesponddotco/accio127/tree/trunk/item/config/config.example.json)
to whatever location you prefer in your filesystem,
`/etc/accio127/config.json` being the recommended location.

[Grab a TLS certificate for your API](https://certbot.eff.org/) and edit
the `config.json` with your preferred settings and the location of your
TLS certificate. Here's an example configuration file:

```json
{
  "address": "api.accio127.com:1997",
  "pid": "/var/run/accio127.pid",
  "proxy": "45.77.163.57",
  "dsn": "file:/var/lib/accio127/sqlite.db?cache=shared&mode=rwc&_pragma_cache_size=-20000&_journal_mode=WAL&_synchronous=NORMAL",
  "certFile": "/etc/nginx/ssl/api.accio127.com_ecc/fullchain.cer",
  "certKey": "/etc/nginx/ssl/api.accio127.com_ecc/api.accio127.com.key",
  "minTLSVersion": "TLS13",
  "privacyPolicy": "https://accio127.com/privacy"
}
```

Now, to start `accio127`, run this command:

```console
accio127ctl start --config /path/to/your/config.json
```

For production you'll probably want to have a `systemd` service to run
that command for you. Here's a simple example of one.

```console
[Unit]
Description=Privacy-focused public IP address API
Documentation=https://sr.ht/~jamesponddotco/accio127/
ConditionFileIsExecutable=/usr/bin/accio127ctl
ConditionFileNotEmpty=/etc/accio127/config.json
After=network.target nss-lookup.target

[Service]
Type=simple
UMask=117
ExecStart=/usr/bin/accio127ctl start --config /etc/accio127/config.json
ExecStop=/usr/bin/accio127ctl stop --config /etc/accio127/config.json
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
```

For production you'll want to improve this file with sandbox and
security features, but that's beyond the scope of this file.

You'll also need to have a server such as NGINX in front of the service,
as it was written to sit behind one. Here's an example `location` for
NGINX.

```nginx
location / {
  proxy_set_header Host $http_host;
  proxy_http_version 1.1;

  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_ssl_server_name on;
  proxy_pass https://api.accio127.com:1997;
}
```

Again, for production you'll want to improve this `location` and have a
proper NGINX configuration file in place with rate limiting and other
security features, since the service itself don't implement these.

## Client side

You can access the service by either accessing the endpoints directly
with your browser, or by using something like `curl`.

**https://api.accio127.com/v1/ip** — Grab your public IP address.
```console
curl -s https://api.accio127.com/v1/ip
```

**https://api.accio127.com/v1/ip/hashed** — Grab a SHA256 hash of your
public IP address.
```console
curl -s https://api.accio127.com/v1/ip/hashed
```

**https://api.accio127.com/v1/ip/anonymized** — Grab your IP address
with the last two octets removed.
```console
curl -s https://api.accio127.com/v1/ip/anonymized
```

**https://api.accio127.com/v1/metrics** — See how many times the service
has been accessed.
```console
curl -s https://api.accio127.com/v1/metrics
```

**https://api.accio127.com/v1/health** — Check the health of the
service.
```console
curl -s https://api.accio127.com/v1/health
```

**https://api.accio127.com/v1/ping** — Check if the service is live or
not. Useful for monitoring.
```console
curl -s https://api.accio127.com/v1/ping
```
