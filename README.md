Nomtail
=======

Application that aggregates (follow) logs from multiple Nomad allocations
into one output stream. This is similar to running `nomad logs -f <alloc-id>`
but for multiple allocations.

Installation
------------

Download the *nomtail* file from the [https://github.com/dsrkoc/nomtail/releases](releases)
section for your platform and run it.

Usage
-----

For running `nomad alloc logs` one needs allocation id (or at least identifier's unique prefix).
Identifiers change after redeployment so it's not very useful to remember them. To get the
fresh set of identifiers one would use `nomad job status` with job identifier unique prefix.

Nomtail requires just the identifier's prefix:

```
$ nomtail -job-prefix=my-service
```

It reads nomad's address from the `NOMAD_ADDR` environment variable. This can be overridden:

```
$ nomtail -job-prefix=my-service -nomad='http://localhost:4646'
```

Supply `-help` for additional options:

```
$ nomtail -help
```


Notes
-----

This project is inspired by [https://github.com/johanhaleby/kubetail](kubetail), which does
the same thing, but for Kubernetes pods. Note that kubetail is written is *bash*, whereas this
project uses *go* (the reason being your author never used go before and this seemed like a
nice project to give it a go (bad pun actually intended, sorry)).
