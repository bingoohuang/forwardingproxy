# Forwarding HTTP/S Proxy

A forwarding HTTP/S proxy. This server is useful when one wants to have
originating requests to a destination service from a set of well-known IPs.

## Usage

```sh
$ fproxy -h
Usage of fproxy:
  -addr                       string  Server address (default ":0")
  -auth                       string  Server authentication username:password
  -avoid                      string Site to be avoided
  -log                        string   Log level (default "info")
  -ca                         string   Filepath to certificate and private key, like -ca cert.pem,key.pem ()
  -le                                  Use letsencrypt for https
  -le.cache.dir               string   Cache directory for certificates
  -le.whitelist               string   Hostname to whitelist for letsencrypt (default "localhost")
  -server.idle.timeout        duration Server idle timeout (default 30s)
  -server.read.header.timeout duration Server read header timeout (default 30s)
  -server.read.timeout        duration Server read timeout (default 30s)
  -server.write.timeout       duration Server write timeout (default 30s)
  -client.read.timeout        duration Client read timeout (default 5s)
  -client.write.timeout       duration Client write timeout (default 5s)
  -dest.dial.timeout          duration Destination dial timeout (default 10s)
  -dest.read.timeout          duration Destination read timeout (default 5s)
  -dest.write.timeout         duration Destination write timeout (default 5s)
```

1. Proxy as HTTP server: fproxy
2. Proxy as HTTPS server: fproxy -ca cert.pem,key.pem
3. Enable Lets Encrypt `fproxy -le -le.whitelist proxy.somehostname.tld -le.cachedir /somewhere/.fproxycache`

The server can be configured to run on a specific interface and port (`-addr`),
be protected via `Proxy-Authorization` (`-auth`). Additionally, most
timeouts can be customized.

To enable verbose logging output, use `-verbose` flag.

## demo

```sh
bingoobjca@bogon a % fproxy -addr :7777 -log debug

{"level":"info","ts":"2022-07-13T14:53:50.457+0800","caller":"fproxy/main.go:196","msg":"http server starting","Listening":"[::]:7777"}
{"level":"info","ts":"2022-07-13T14:53:50.457+0800","caller":"fproxy/main.go:198","msg":"settings: export http_proxy=http://127.0.0.1:7777; export https_proxy=http://127.0.0.1:7777"}
{"level":"info","ts":"2022-07-13T14:53:57.237+0800","caller":"fproxy/proxy.go:37","msg":"Incoming request","host":"127.0.0.1:5003"}
{"level":"debug","ts":"2022-07-13T14:53:57.237+0800","caller":"fproxy/proxy.go:76","msg":"Connecting","host":"127.0.0.1:5003"}
{"level":"debug","ts":"2022-07-13T14:53:57.237+0800","caller":"fproxy/proxy.go:85","msg":"Connected","host":"127.0.0.1:5003"}
{"level":"debug","ts":"2022-07-13T14:53:57.237+0800","caller":"fproxy/proxy.go:87","msg":"Hijacking","host":"127.0.0.1:5003"}
{"level":"debug","ts":"2022-07-13T14:53:57.237+0800","caller":"fproxy/proxy.go:106","msg":"Hijacked connection","host":"127.0.0.1:5003"}
{"level":"debug","ts":"2022-07-13T14:53:57.293+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":">>> 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":267,"total":267}
{"level":"debug","ts":"2022-07-13T14:53:57.296+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":"<<< 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":1494,"total":1494}
{"level":"debug","ts":"2022-07-13T14:53:57.337+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":">>> 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":75,"total":342}
{"level":"debug","ts":"2022-07-13T14:53:57.364+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":">>> 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":6,"total":348}
{"level":"debug","ts":"2022-07-13T14:53:57.391+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":">>> 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":45,"total":393}
{"level":"debug","ts":"2022-07-13T14:53:57.392+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":"<<< 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":51,"total":1545}
{"level":"debug","ts":"2022-07-13T14:53:57.394+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":">>> 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":180,"total":573}
{"level":"debug","ts":"2022-07-13T14:53:57.396+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":"<<< 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":341,"total":1886}
{"level":"debug","ts":"2022-07-13T14:53:57.412+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":">>> 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":180,"total":753}
{"level":"debug","ts":"2022-07-13T14:53:57.449+0800","caller":"fproxy/debugconn.go:33","msg":"transferred","direction":"<<< 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","bytes":341,"total":2227}
{"level":"debug","ts":"2022-07-13T14:53:57.810+0800","caller":"fproxy/debugconn.go:40","msg":"read error","direction":">>> 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","total":753,"error":"EOF"}
{"level":"debug","ts":"2022-07-13T14:53:57.810+0800","caller":"fproxy/debugconn.go:40","msg":"read error","direction":"<<< 127.0.0.1:58344-127.0.0.1:7777-127.0.0.1:58345-127.0.0.1:5003","total":2227,"error":"read tcp 127.0.0.1:58345->127.0.0.1:5003: use of closed network connection"}
^C{"level":"info","ts":"2022-07-13T14:54:18.381+0800","caller":"fproxy/main.go:132","msg":"Server shutting down"}
{"level":"info","ts":"2022-07-13T14:54:18.382+0800","caller":"fproxy/main.go:157","msg":"Server stopped"}
```

```sh
$ fproxy
{"level":"info","ts":"2022-07-06T09:35:04.196+0800","caller":"fproxy/main.go:191","msg":"Server starting","Listening":"[::]:54741"}
{"level":"info","ts":"2022-07-06T09:36:18.941+0800","caller":"fproxy/proxy.go:37","msg":"Incoming request","host":"127.0.0.1:5003"}
{"level":"info","ts":"2022-07-06T09:37:24.379+0800","caller":"fproxy/proxy.go:37","msg":"Incoming request","host":"127.0.0.1:5004"}
```

```sh
$ gurl http://127.0.0.1:5003/v -proxy http://127.0.0.1:54741
2022/07/06 09:36:18.940579 main.go:131: Proxy URL: http://127.0.0.1:54741/
Conn-Session: 127.0.0.1:54927->127.0.0.1:54741 (reused: false, wasIdle: false, idle: 0s)
GET /v HTTP/1.1
Host: 127.0.0.1:5003
Accept: application/json
Accept-Encoding: gzip, deflate
Content-Type: application/json
Gurl-Date: Wed, 06 Jul 2022 01:36:18 GMT
User-Agent: gurl/1.0.0


HTTP/1.1 200 OK
Vary: Accept-Encoding
X-Forwarded-For: 127.0.0.1
Content-Encoding: gzip
Content-Length: 139
Content-Type: application/json; charset=utf-8
Date: Wed, 06 Jul 2022 01:36:18 GMT

{
  "build": "2022-06-28T15:38:41+0800",
  "git": "master-b524d91@2022-06-28T15:38:05+08:00",
  "go": "go1.18.3_darwin/amd64",
  "version": "1.0.0"
}

   DNS Lookup   TCP Connection   Request Transfer   Server Processing   Response Transfer
[       0 ms  |          0 ms  |            0 ms  |            0 ms  |             0 ms  ]
              |                |                  |                  |                   |
  namelookup: 0 ms             |                  |                  |                   |
                      connect: 0 ms               |                  |                   |
                                   wrote request: 0 ms               |                   |
                                                      starttransfer: 1 ms                |
                                                                                  total: 1 ms
2022/07/06 09:36:18.942304 main.go:176: current request cost: 1.64843ms
2022/07/06 09:36:18.942315 main.go:69: complete, total cost: 1.812369ms

$ gurl https://127.0.0.1:5004/v -proxy http://127.0.0.1:54741 -i
2022/07/06 09:37:24.379036 main.go:131: Proxy URL: http://127.0.0.1:54741/
Conn-Session: 127.0.0.1:55092->127.0.0.1:54741 (reused: false, wasIdle: false, idle: 0s)
GET /v HTTP/1.1
Host: 127.0.0.1:5004
Accept: application/json
Accept-Encoding: gzip, deflate
Content-Type: application/json
Gurl-Date: Wed, 06 Jul 2022 01:37:24 GMT
User-Agent: gurl/1.0.0


HTTP/1.1 200 OK
Date: Wed, 06 Jul 2022 01:37:24 GMT
Content-Length: 139
Content-Encoding: gzip
Content-Type: application/json; charset=utf-8
Vary: Accept-Encoding

{
  "build": "2022-06-28T15:38:41+0800",
  "git": "master-b524d91@2022-06-28T15:38:05+08:00",
  "go": "go1.18.3_darwin/amd64",
  "version": "1.0.0"
}

  DNS Lookup   TCP Connection   TLS Handshake   Request Transfer   Server Processing   Response Transfer
[      0 ms  |          0 ms  |         2 ms  |           0 ms  |             0 ms  |            0 ms  ]
             |                |               |                 |                   |                  |
 namelookup: 0 ms             |               |                 |                   |                  |
                     connect: 0 ms            |                 |                   |                  |
                                 pretransfer: 2 ms              |                   |                  |
                                                 wrote request: 2 ms                |                  |
                                                                     starttransfer: 3 ms               |
                                                                                                total: 3 ms
2022/07/06 09:37:24.382855 main.go:176: current request cost: 3.742847ms
2022/07/06 09:37:24.382869 main.go:69: complete, total cost: 3.918938ms
```

Can use httplive to create localhost.pem and localhost.key for demo

```sh
$ httplive -k -l
2022-07-06 09:57:51.294 [INFO ] 6837 --- [1    ] [-]  : log file created: /Users/bingoobjca/logs/b/httplive.log
2022-07-06 09:57:52.400 [INFO ] 6837 --- [1    ] [-]  : Created a new local CA ✅
2022-07-06 09:57:52.574 [INFO ] 6837 --- [1    ] [-]  : Created a new certificate valid for the name - "localhost" 📜
2022-07-06 09:57:52.574 [INFO ] 6837 --- [1    ] [-]  : The certificate is at ".cert/localhost.pem" and the key at ".cert/localhost.key" ✅
2022-07-06 09:57:52.574 [INFO ] 6837 --- [1    ] [-]  : The certificate will expire at 2024-07-06 01:57:52 🗓
```

```sh
$ fproxy -ca localhost.pem,localhost.key
{"level":"info","ts":"2022-07-06T09:38:33.450+0800","caller":"fproxy/main.go:175","msg":"Server starting","Listening":"[::]:55254"}
{"level":"info","ts":"2022-07-06T09:38:58.283+0800","caller":"fproxy/proxy.go:37","msg":"Incoming request","host":"127.0.0.1:5003"}
{"level":"info","ts":"2022-07-06T09:39:45.639+0800","caller":"fproxy/proxy.go:37","msg":"Incoming request","host":"127.0.0.1:5004"}
```

```sh
$ gurl http://127.0.0.1:5003/v -proxy https://127.0.0.1:55254 -i
2022/07/06 09:38:58.280752 main.go:131: Proxy URL: https://127.0.0.1:55254/
Conn-Session: 127.0.0.1:55314->127.0.0.1:55254 (reused: false, wasIdle: false, idle: 0s)
GET /v HTTP/1.1
Host: 127.0.0.1:5003
Accept: application/json
Accept-Encoding: gzip, deflate
Content-Type: application/json
Gurl-Date: Wed, 06 Jul 2022 01:38:58 GMT
User-Agent: gurl/1.0.0


HTTP/1.1 200 OK
Date: Wed, 06 Jul 2022 01:38:58 GMT
Vary: Accept-Encoding
X-Forwarded-For: 127.0.0.1
Content-Encoding: gzip
Content-Length: 139
Content-Type: application/json; charset=utf-8

{
  "build": "2022-06-28T15:38:41+0800",
  "git": "master-b524d91@2022-06-28T15:38:05+08:00",
  "go": "go1.18.3_darwin/amd64",
  "version": "1.0.0"
}

   DNS Lookup   TCP Connection   Request Transfer   Server Processing   Response Transfer
[       0 ms  |          2 ms  |            0 ms  |            1 ms  |             0 ms  ]
              |                |                  |                  |                   |
  namelookup: 0 ms             |                  |                  |                   |
                      connect: 2 ms               |                  |                   |
                                   wrote request: 2 ms               |                   |
                                                      starttransfer: 3 ms                |
                                                                                  total: 3 ms
2022/07/06 09:38:58.284910 main.go:176: current request cost: 4.083942ms
2022/07/06 09:38:58.284922 main.go:69: complete, total cost: 4.253366ms

$ gurl https://127.0.0.1:5004/v -proxy https://127.0.0.1:55254 -i
2022/07/06 09:39:45.636832 main.go:131: Proxy URL: https://127.0.0.1:55254/
Conn-Session: 127.0.0.1:55427->127.0.0.1:55254 (reused: false, wasIdle: false, idle: 0s)
GET /v HTTP/1.1
Host: 127.0.0.1:5004
Accept: application/json
Accept-Encoding: gzip, deflate
Content-Type: application/json
Gurl-Date: Wed, 06 Jul 2022 01:39:45 GMT
User-Agent: gurl/1.0.0


HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Type: application/json; charset=utf-8
Vary: Accept-Encoding
Date: Wed, 06 Jul 2022 01:39:45 GMT
Content-Length: 139

{
  "build": "2022-06-28T15:38:41+0800",
  "git": "master-b524d91@2022-06-28T15:38:05+08:00",
  "go": "go1.18.3_darwin/amd64",
  "version": "1.0.0"
}

  DNS Lookup   TCP Connection   TLS Handshake   Request Transfer   Server Processing   Response Transfer
[      0 ms  |          0 ms  |         2 ms  |           0 ms  |             0 ms  |            0 ms  ]
             |                |               |                 |                   |                  |
 namelookup: 0 ms             |               |                 |                   |                  |
                     connect: 0 ms            |                 |                   |                  |
                                 pretransfer: 5 ms              |                   |                  |
                                                 wrote request: 5 ms                |                  |
                                                                     starttransfer: 5 ms               |
                                                                                                total: 6 ms
2022/07/06 09:39:45.643151 main.go:176: current request cost: 6.237827ms
2022/07/06 09:39:45.643161 main.go:69: complete, total cost: 6.415185ms
```

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

## HTTP 代理

代理是网络中的一项重要的功能，其功能就是代理网络用户去取得网络信息。形象的说：它是网络信息的中转站，对于客户端来说，代理扮演的是服务器的角色，接收请求报文，返回响应报文；对于web服务器来说，代理扮演的是客户端的角色，发送请求报文，接收响应报文。

代理具有多种类型，如果是根据网络用户划分的话，可以划分为正向代理和反向代理：

1. 正向代理：将客户端作为网络用户。客户端访问服务端时，先访问代理服务器，随后代理服务器再访问服务端。此过程需客户端进行代理配置，对服务端透明。
2. 反向代理：将服务端作为网络用户。访问过程与正向代理相同，不过此过程对客户端透明，需服务端进行代理配置（也可不配置）。

针对正向代理和反向代理，分别有不同的代理协议，即代理服务器和网络用户之间通信所使用的协议：

正向代理： http https socks4 socks5 vpn：就功能而言，vpn也可以被认为是代理
反向代理： tcp udp http https

HTTP 代理是正向代理中较为简单的代理方式，它使用 HTTP 协议作为客户端和代理服务器的传输协议。

HTTP 代理可以承载 HTTP 协议，HTTPS 协议，FTP 协议等等。对于不同的协议，客户端和代理服务器间的数据格式略有不同。

### HTTP Forwarding Proxy

HTTP Proxy 如同一个中间人

客户端浏览器向它发起 TCP 连接，并把实际需要访问的 web service 放在 GET 和 HOST HTTP header 中。
它收到请求后，解析 HTTP header，并向真正的 web service 发起请求。
稍后它再将 web service 返回的内容返回至客户端。
在上述这个过程中，HTTP Proxy 同时维持了2个 TCP 连接。因为它非常清楚地知道客户端和服务端之间的对话内容，所以带来了安全隐患。
另外有些无良的 Proxy 还会自己往服务端所返回的 HTML 页面里面添加烦人的广告。

```sh
+--------+             +-------+             +-------------+
| Client |             | Proxy |             | Destination |
+--------+             +-------+             +-------------+
    |   1. HTTP Request    |                       |
    | -------------------> |   2. HTTP Request     |
    |                      | --------------------> |
    |                      |   3. HTTP Response    |
    |   4. HTTP Response   | <-------------------- |
    | <------------------- |                       x
    x                      :
```

To support HTTP we’ll use built-in HTTP server and client. The role of proxy is to handle HTTP request, pass such
request to destination server and send response back to the client.

![img.png](_assets/2022-07-02-09-09-16.png)

### HTTP CONNECT tunneling

Suppose client wants to use either HTTPS or WebSockets in order to talk to server. Client is aware of using proxy.
Simple HTTP request / response flow cannot be used since client needs to e.g. establish secure connection with server (
HTTPS) or wants to use other protocol over TCP connection (WebSockets). Technique which works is to use
HTTP [CONNECT](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/CONNECT)
method. It tells the proxy server to establish TCP connection with destination server and when done to proxy the TCP
stream to and from the client. This way proxy server won’t terminate SSL but will simply pass data between client and
destination server so these two parties can establish secure connection.

```sh
> CONNECT example.host.com:443 HTTP/1.1
> Host: example.host.com:443
> Proxy-Authorization: Basic base64-encoded-proxy-credentials
> Proxy-Connection: Keep-Alive
< HTTP/1.1 200 OK
> GET /foo/bar?baz#qux HTTP/1.1
> Host: example.host.com
> Authorization: Basic base64-encoded-destination-credentials
< HTTP/1.1 200 OK
< Connection: close
+--------+              +-------+             +-------------+
| Client |              | Proxy |             | Destination |
+--------+              +-------+             +-------------+
    |   1. HTTP/S CONNECT    |                       |
    | ---------------------> |   2. TLS handshake    |
    |                        | --------------------> |
    |                        |   3. Established      |
    |   4. HTTP/S 200 OK     | <-------------------- |
    | <--------------------- |                       |
    |   5. HTTPS Request     |                       |
    | ---------------------> |   6. TCP stream       |
    |                        | --------------------> |
    |                        |   7. TCP stream       |
    |   8. HTTPS Response    | <-------------------- |
    | <--------------------- |                       x
    x                        :
```

时代在变，技术在变，需求也在变。如今放眼望去，大部分网站都在使用 HTTPS 。这个时候再使用 HTTP Proxy 的话，就会出现一个问题：只有 Proxy 能看到 web service 的证书，但客户端却只能接收到 Proxy 的证书。很明显，Proxy 证书里面的名字和客户端所想要访问的地址完全不同，这个时候浏览器会毫不留情地给出一个警告，TLS handshake 也没法通过。

为了解决这个问题，HTTPS Proxy 应运而生。图 2 展示了它的工作流程。

1. ① 客户端浏览器首先向 HTTPS Proxy 发起 TCP 连接，但在 HTTP 请求里，用的是 CONNECT 方法。
2. ② HTTPS Proxy 收到这个连接后，利用 CONNECT 请求里面的数据，向 google.com 发起 TCP 连接。
3. ③ 从此以后 HTTPS Proxy 只负责 relay 客户端和 google.com 之间的数据，不会做任何解析，也没有办法解析，因为它根本没有参与 TLS handshake 这个过程。它所 relay 的数据包括：TLS handshake 数据包(如证书、加密算法、加密秘钥等)、加密过之后的HTTP payload。

在步骤 ③ 这里，HTTPS Proxy 通过对客户端及服务端会话数据的透明转发，实现了 tunnel 的效果。

![HTTPS Proxy](_assets/2022-07-13-13-54-04.png)

图 2：通过 HTTPS Proxy ，客户端和服务器直接创建了一个 tunnel

需要强调的是上述过程中，浏览器从始至终都没有和 google.com 发起直接的TCP连接，通信双方所有的数据都是经过 HTTPS Proxy中转的。

由于 HTTPS 下客户端和服务端的通信除了开头的协商以外都是密文，中间的代理服务器不再承担修改 HTTPS 报文再转发的功能，而是一开始就和客户端协商好服务端的地址，随后的 TCP 密文直接转发即可。

![通过 HTTPS Proxy ，客户端和服务器直接创建了一个 tunnel](_assets/2022-07-02-09-10-15.png)

Presented code is not a production-grade solution. It lacks e.g.
handling [hop-by-hop headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#hbh), setting up timeouts while
copying data between two connections or the ones exposed by net/http — more on this in “[The complete guide to Go
net/http timeouts](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)”.

### Go net/http timesouts

Server Timeouts

![Server Timeouts](_assets/2022-07-06-10-17-09.png)

There are three main timeouts exposed in http.Server: ReadTimeout, WriteTimeout and IdleTimeout. You set them by
explicitly using a Server:

```go
srv := &http.Server{
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 10 * time.Second,
	IdleTimeout:  120 * time.Second,
	TLSConfig:    tlsConfig,
	Handler:      serveMux,
}
log.Println(srv.ListenAndServeTLS("", ""))
```

Client Timeouts

![Client Timeouts](_assets/2022-07-06-10-19-08.png)

```go
c := &http.Client{
Transport: &http.Transport{
	Dial: (&net.Dialer{
		imeout:   30 * time.Second,
		eepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout:   10 * time.Second,
	ResponseHeaderTimeout: 10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
}
```

## Resources

1. [golang实现一个简单的http代理](https://blog.51cto.com/u_15127596/3297318)
2. [sipt/shuttle](https://github.com/sipt/shuttle) A web proxy in Golang with amazing features.
3. HTTP proxy written in Go. [COW](https://github.com/cyfdecyf/cow) can automatically identify blocked sites and use
   parent proxies to access.
4. HTTP(S) Proxy in Golang in less than 100 lines of code [blog](https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c)
5. [手边的Tunnel知多少](https://www.51cto.com/article/708337.html)
6. [HTTP 代理原理及实现](https://imququ.com/post/web-proxy.html)
