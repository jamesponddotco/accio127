# Accio127

[![builds.sr.ht status](https://builds.sr.ht/~jamesponddotco/accio127.svg)](https://builds.sr.ht/~jamesponddotco/accio127?)

**Accio127** is a super simple, privacy-respecting, and ad-free API you
can use to get your public IP address. It's pretty fast and reliable.

Use the official instance at
[api.accio127.com](https://api.accio127.com/v1/ip) or host the service
yourself. It should be well-suited for both small and large deployments.

## Usage

* [Getting started](doc/getting-started.md)
* [OpenAPI Specification](doc/openapi.json)

## Installation

### From source

First install the dependencies:

- Go 1.20 or above.
- make.
- sqlite.
- [scdoc](https://git.sr.ht/~sircmpwn/scdoc).

Then compile and install:

```bash
make
sudo make install
```

## Contributing

Anyone can help make `accio127` better. Send patches on the [mailing
list](https://lists.sr.ht/~jamesponddotco/accio127-devel) and report
bugs on the [issue
tracker](https://todo.sr.ht/~jamesponddotco/accio127).

You must sign-off your work using `git commit --signoff`. Follow the
[Linux kernel developer's certificate of
origin](https://www.kernel.org/doc/html/latest/process/submitting-patches.html#sign-your-work-the-developer-s-certificate-of-origin)
for more details.

All contributions are made under [the EUPL license](LICENSE.md).

## Resources

The following resources are available:

- [Support and general discussions](https://lists.sr.ht/~jamesponddotco/accio127-discuss).
- [Patches and development related questions](https://lists.sr.ht/~jamesponddotco/accio127-devel).
- [Instructions on how to prepare patches](https://git-send-email.io/).
- [Feature requests and bug reports](https://todo.sr.ht/~jamesponddotco/accio127).

---

Released under the [EUPL License](LICENSE.md).
