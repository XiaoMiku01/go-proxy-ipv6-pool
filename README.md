# Go Proxy IPV6 Pool

Random ipv6 egress proxy server (support http/socks5) 
The Go language implementation of [zu1k/http-proxy-ipv6-pool](https://github.com/zu1k/http-proxy-ipv6-pool)

## Usage

```bash
    go run . --port <port> --cidr < your ipv6 cidr >  # e.g. 2001:399:8205:ae00::/64
```

### Use as a proxy server

```bash
    curl -x http://xxx:52122 http://6.ipw.cn/ # 2001:399:8205:ae00:456a:ab12 (random ipv6 address)
```

```bash
    curl -x socks5://xxx:52123 http://6.ipw.cn/ # 2001:399:8205:ae00:456a:ab12 (random ipv6 address)
```

## License

MIT License (see [LICENSE](LICENSE))
