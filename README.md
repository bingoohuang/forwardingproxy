# Forwarding HTTP/S Proxy

A forwarding HTTP/S proxy. This server is useful when one wants to have
originating requests to a destination service from a set of well-known IPs.

## Usage

```
$ forwardingproxy -h
Usage of forwardingproxy:
  -addr string
        Server address
  -auth string
        Server authentication username:password
  -avoid string
        Site to be avoided
  -cert string
        Filepath to certificate
  -client.read.timeout duration
        Client read timeout (default 5s)
  -client.write.timeout duration
        Client write timeout (default 5s)
  -dest.dial.timeout duration
        Destination dial timeout (default 10s)
  -dest.read.timeout duration
        Destination read timeout (default 5s)
  -dest.write.timeout duration
        Destination write timeout (default 5s)
  -key string
        Filepath to private key
  -le
        Use letsencrypt for https
  -le.cache.dir string
        Cache directory for certificates (default "/tmp")
  -le.whitelist string
        Hostname to whitelist for letsencrypt
  -server.idle.timeout duration
        Server idle timeout (default 30s)
  -server.read.header.timeout duration
        Server read header timeout (default 30s)
  -server.read.timeout duration
        Server read timeout (default 30s)
  -server.write.timeout duration
        Server write timeout (default 30s)
  -verbose
        Set log level to DEBUG
```

To start the proxy as HTTP server, just run:

```
$ forwardingproxy
```

To start the proxy as HTTPS server, just provide server certificate and private
key files:

```
$ forwardingproxy -cert cert.pem -key key.pem
```

To create a self-signed certificate and private key for testing, run:

```
$ openssl req -newkey rsa:2048 -nodes -keyout key.pem -new -x509 -sha256 -days 3650 -out cert.pem

```

Or enable letsencrypt

```
$ forwardingproxy -le -le.whitelist proxy.somehostname.tld -le.cachedir /home/somewhere/.forwardingproxycache
```

The server can be configured to run on a specific interface and port (`-addr`),
be protected via `Proxy-Authorization` (`-auth`). Additionally, most
timeouts can be customized.

To enable verbose logging output, use `-verbose` flag.

## Implementation details

It is a simple HTTPS tunneling proxy that starts a Go HTTPS server at a given
port awaiting `CONNECT` requests, basically dropping everything else. To start
the HTTPS server one has to provide a server certificate and private key for the
TLS handshake phase.

Once a client requests a `CONNECT` it will create a TCP connection to the
provided destination host, and on successfully establishing this connection,
hijack the original client connection, and transparently and bidirectionally
copying incoming and outgoing TCP byte streams.

It has minimal logging using Uber's Zap logger.

## Features

This is NIH (_Not Invented Here_ syndrome), thus quality and feature set is not
en-par with hosted or off-the-shelf solutions.

Compared to especially hosted solutions, insight into the proxies operations
such as logging, monitoring, usage statistics need to be added if desired.
Additionally, one has to setup the binary as a reliable server and automate
deployments.

## Background

If one has a third-party requirement to have server requests originating from a
fixed IP address, there are mainly two options: (i) host code on a cloud
provider such as an EC2 instance and connect the instance to an EIP (Elastic
IP). (ii) But if code is hosted on a PaaS provider with no guarantee of a fixed
IP, such as Heroku, one would proxy requests through a proxy server and have
that proxy server attached to a fixed IP.

### Proxy

To proxy HTTPS requests, one broadly has two options in software: Use the [HTTP
Tunnel](https://en.wikipedia.org/wiki/HTTP_tunnel) feature via the
[CONNECT](https://www.ietf.org/rfc/rfc2817.txt) method, also called a
[Forwarding Proxy](https://en.wikipedia.org/wiki/Proxy_server), or a [Reverse
Proxy](https://en.wikipedia.org/wiki/Reverse_proxy). There are hardware
solutions on OSI layer 3 instead of layer 7, namely a NAT proxy, but this is not
discussed here as it more convenient nowadays to not require access to physical
hardware or want to invest into a NAT proxy e.g. on AWS.

#### Forwarding Proxy

A forwarding proxy can come in two flavours:

One in which the proxy terminates an incoming client request, evaluates it, and
forwards the request to a destination. This works for HTTP as well as for HTTPS.
A subtle but important side effect of using a forwarding proxy for HTTPS is that
it would terminate the request, thus being able to inspect the request's content
(and modify it).

The other in which the proxy uses tunneling via the `CONNECT` method. By this, a
proxy accepts an initial CONNECT request entailing the entire URL as `HOST`
value, rather than just the host address. The proxy then opens up a TCP
connection to the destination and transparently forwards the raw communication
from the client to the destination. This comes with the subtle difference that
only the initial `CONNECT` request from the client to the proxy is terminated
and can be analyzed, however, any further communication is not terminated nor
intercepted thus SSL communication can't be read by the proxy.

One additional subtle thing to mention is that forwarding proxies using
tunneling rely on clients to understand and comply to the HTTP tunneling RFC
2817 and thus have to be explicitly configured to use HTTP/S proxying, usually
picking up the proxy url from environment variables such as `HTTP_PROXY` and
`HTTPS_PROXY` respectively. For Go, see
[net/http/Transport](https://golang.org/pkg/net/http/#Transport) and

#### Reverse Proxy

A reverse proxy accepts incoming requests from clients and routes them to a
specific destination based on the request. Discriminators for destination could
be the host, path, query parameters, any header even the body. A reverse proxy,
in any case, intercepts and terminates the HTTP and HTTPS connections and
creates new requests to the destinations.

## Alternatives

One can consider alternatives in protocol and product, as well as _Make vs Buy_.

Alternatives to a forwarding proxy would be a reverse proxy or NAT proxy.
Reverse proxies as outlined above would require more routing logic. A NAT proxy
would require more network or hardware configuration even if it's abstracted by
IaaS providers.

Alternative products would be [Squid](http://squid-cache.org),
[NGINX](https://www.nginx.com), [HAProxy](http://haproxy.org),
[Varnish](http://varnish-cache.org), [TinyProxy](https://tinyproxy.github.io),
etc. NGINX is primarily meant to be used as a reverse proxy and can be difficult
to set up acting as a forwarding proxy. HAProxy is similar to NGINX meant to be
used as reverse proxy as well as HTTP cache. Varnish is primarily meant to be
used as HTTP cache. Squid and TinyProxy are closest to be working as forwarding
proxies, however it can be difficult to set them up on a new Amazon Linux 2 AMI
EC2 instance, and TinyProxy is not maintained anymore since several years.

## License

MIT License.

## HTTP(S) Proxy in Golang in less than 100 lines of code

[blog](https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c)

### HTTP

To support HTTP we’ll use built-in HTTP server and client. The role of proxy is to handle HTTP request, pass such
request to destination server and send response back to the client.

![img.png](_assets/img.png)

### HTTP CONNECT tunneling

Suppose client wants to use either HTTPS or WebSockets in order to talk to server. Client is aware of using proxy.
Simple HTTP request / response flow cannot be used since client needs to e.g. establish secure connection with server (
HTTPS) or wants to use other protocol over TCP connection (WebSockets). Technique which works is to use
HTTP [CONNECT](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/CONNECT)
method. It tells the proxy server to establish TCP connection with destination server and when done to proxy the TCP
stream to and from the client. This way proxy server won’t terminate SSL but will simply pass data between client and
destination server so these two parties can establish secure connection.

![img.png](_assets/img02.png)

Presented code is not a production-grade solution. It lacks e.g.
handling [hop-by-hop headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#hbh), setting up timeouts while
copying data between two connections or the ones exposed by net/http — more on this in “[The complete guide to Go
net/http timeouts](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)”.
