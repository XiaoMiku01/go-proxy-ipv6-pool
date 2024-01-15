package main

import (
	"io"
	"log"
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
)

var httpProxy = goproxy.NewProxyHttpServer()

func init() {
	httpProxy.Verbose = true

	httpProxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			// 为 IPv6 地址添加方括号
			outgoingIP, err := generateRandomIPv6(cidr)
			if err != nil {
				log.Printf("Generate random IPv6 error: %v", err)
				return req, nil
			}
			outgoingIP = "[" + outgoingIP + "]"
			// 使用指定的出口 IP 地址创建连接
			localAddr, err := net.ResolveTCPAddr("tcp", outgoingIP+":0")
			if err != nil {
				log.Printf("[http] Resolve local address error: %v", err)
				return req, nil
			}
			dialer := net.Dialer{
				LocalAddr: localAddr,
			}

			// 通过代理服务器建立到目标服务器的连接
			// 发送 http 请求
			// 使用自定义拨号器设置 HTTP 客户端
			// 创建新的 HTTP 请求

			newReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
			if err != nil {
				log.Printf("[http] New request error: %v", err)
				return req, nil
			}
			newReq.Header = req.Header

			// 设置自定义拨号器的 HTTP 客户端
			client := &http.Client{
				Transport: &http.Transport{
					DialContext: dialer.DialContext,
				},
			}

			// 发送 HTTP 请求
			resp, err := client.Do(newReq)
			if err != nil {
				log.Printf("[http] Send request error: %v", err)
				return req, nil
			}
			return req, resp
		},
	)

	httpProxy.OnRequest().HijackConnect(
		func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
			// 通过代理服务器建立到目标服务器的连接
			outgoingIP, err := generateRandomIPv6(cidr)
			if err != nil {
				log.Printf("Generate random IPv6 error: %v", err)
				return
			}
			outgoingIP = "[" + outgoingIP + "]"
			// 使用指定的出口 IP 地址创建连接
			localAddr, err := net.ResolveTCPAddr("tcp", outgoingIP+":0")
			if err != nil {
				log.Printf("[http] Resolve local address error: %v", err)
				return
			}
			dialer := net.Dialer{
				LocalAddr: localAddr,
			}

			// 通过代理服务器建立到目标服务器的连接
			server, err := dialer.Dial("tcp", req.URL.Host)
			if err != nil {
				log.Printf("[http] Dial to %s error: %v", req.URL.Host, err)
				client.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
				client.Close()
				return
			}

			// 响应客户端连接已建立
			client.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
			// 从客户端复制数据到目标服务器
			go func() {
				defer server.Close()
				defer client.Close()
				io.Copy(server, client)
			}()

			// 从目标服务器复制数据到客户端
			go func() {
				defer server.Close()
				defer client.Close()
				io.Copy(client, server)
			}()

		},
	)
}
