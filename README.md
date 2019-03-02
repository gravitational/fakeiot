# Fakeiot

Fake IOT test cluster is a client side implementation
for Gravitational Full Stack Engineer coding challenge.

## Protocol

**Request Format**

Every time a user logs in into imaginary Fake IOT device,
it connects to the metrics server and issues the following
JSON POST request:

```
HTTP POST /metrics

{
  "account_id": "781df840-09da-42f4-ba29-996d2ff76a73",
  "user_id": "bf506b23-8c4e-4c8e-af95-e331dba766ab",
  "timestamp": "2019-03-03T18:02:30.424878129Z"
}
```

* `account_id` is UUID of the user account.
* `user_id` is UUID of the user
* `timestamp` is an [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) timestamp
of the login in UTC timezone.

**Authentication**

Client uses [Bearer Auth](https://tools.ietf.org/html/rfc6750#page-5) headers
to authenticate the request.

**Transport Security**

The `fakeiot` client will only connect to the target
server if it trusts the server's x509 certificate to establish a TLS connection
and use valid HTTPS.

You could use [Letsencrypt](https://letsencrypt.org/) for your test server,
or generate your own certificate and tell `fakeiot` to trust the cert's
certificate authority using `--ca-cert` command line flag.

## Tool Usage

**Building**

`fakeiot` could be installed using go:

```bash
$ go install github.com/gravitational/fakeiot
```

**Running Tests**

Once you create your first server, we recommend to run a set
of compliance tests on it:

```bash
$ fakeiot --token=shmoken --url="https://127.0.0.1:8443" --ca-cert=./fixtures/ca-cert.pem test
2019-03-03T10:18:56-08:00 INFO [RUNNER]    Starting compliance tests. runner/runner.go:110
2019-03-03T10:18:56-08:00 INFO [RUNNER]    [PASS] Sending OK request. runner/runner.go:137
2019-03-03T10:18:56-08:00 INFO [RUNNER]    [PASS] Sending Bogus request. runner/runner.go:164
2019-03-03T10:18:56-08:00 INFO [RUNNER]    [PASS] Sending Bogus request. runner/runner.go:164
2019-03-03T10:18:56-08:00 INFO             Fake IOT program run successfully. fakeiot/main.go:45
```

**Running A Demo**

To generate some traffic, you can run a simulation:

```bash
# 3 users log in sequentially over 10 seconds with 1 second frequency
$ fakeiot --token=shmoken --url="https://127.0.0.1:8443" --ca-cert=./fixtures/ca-cert.pem run --period=10s --freq=1s --users=3
```

Try generating more traffic to see how your server performs

```bash
# 3 users log in sequentially over 10 seconds with 1 second frequency
$ fakeiot --token=shmoken --url="https://127.0.0.1:8443" --ca-cert=./fixtures/ca-cert.pem run --period=10s --freq=1s --users=100
```

## Feedback 

If you have any feedback, please create an issue in this repository.

Good luck!
