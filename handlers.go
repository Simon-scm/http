package http

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func writeResponse(w *bufio.Writer, status int, statusText string, contentType string, body []byte) error {
	if statusText == "" {
		statusText = "OK"
	}

	if contentType == "" {
		contentType = "text/plain; charset=utf-8"
	}

	// 1) Status-line
	fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", status, statusText)

	// 2) Header
	fmt.Fprintf(w, "Date: %s\r\n", time.Now().UTC().Format(time.RFC1123))
	fmt.Fprintf(w, "Server: local-go-mini/0.1\r\n")
	fmt.Fprintf(w, "Content-Type: %s\r\n", contentType)
	fmt.Fprintf(w, "Content-Length: %d\r\n", len(body))

	// 3) Empty line
	fmt.Fprintf(w, "\r\n")

	// 4) Body
	if len(body) > 0 {
		if _, err := w.Write(body); err != nil {
			return err
		}
	}

	return w.Flush()
}

func handleConn(c net.Conn) {
	defer c.Close()

	// Safety timeout
	_ = c.SetDeadline(time.Now().Add(30 * time.Second))

	br := bufio.NewReader(c)
	bw := bufio.NewWriter(c)

	req, err := parseRequest(br)
	if err != nil {
		log.Printf("parse error: %v", err)
		_ = writeResponse(bw, 400, "Bad Request", "text/plain; charset=utf-8", []byte("bad reqeust\r\n"))
		return
	}

	// Log request line, headers length, body length
	log.Printf("%s %s %s, Headers=%d, Body=%d", req.Method, req.Target, req.Proto, len(req.Headers), len(req.Body))

	// Routing - currently only 200 respones on /
	_ = writeResponse(bw, 200, "OK", "text/plain; charset=utf-8", []byte("ok\r\n"))
}
