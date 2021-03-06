Nomtail
=======

Application that aggregates (follow) logs from multiple Nomad allocations
into one output stream. This is similar to running `nomad alloc logs -f <allocation>`
but for multiple allocations.

Installation
------------

Download the *nomtail* file from the [releases](https://github.com/dsrkoc/nomtail/releases)
section for your platform and run it.

Usage
-----

For running `nomad alloc logs` one needs allocation identifier (or at least its unique prefix).
 To get the fresh set of identifiers one would use `nomad job status` with job identifier's unique prefix.

Nomtail requires just the job's identifier (or its prefix):

```
$ nomtail my-service
```

It reads nomad's address from the `NOMAD_ADDR` environment variable. This can be overridden:

```
$ nomtail -address='http://localhost:4646' my-service
```

In cases where several allocations frequently output log messages and those messages start with a
well formatted timestamp (like ISO 8601), messages can be printed in a more chronological order,
rather than in the order they were received from allocations. Nomtail can be instructed to
lexicographically sort log messages it keeps in its buffer before outputing them to the stdout:

```
$ nomtail -sort-buffer my-service
```

Supply `-help` for additional options:

```
$ nomtail -help
```


Notes
-----

This project is inspired by [kubetail](https://github.com/johanhaleby/kubetail), which does
the same thing, but for Kubernetes pods. Note that kubetail is written is *bash*, whereas this
project uses *go* (the reason being your author never used go before and this seemed like a
nice project to give it a go (bad pun actually intended, sorry)).
