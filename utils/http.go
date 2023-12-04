package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context/ctxhttp"
)

type HttpOption struct {
	Method      string
	Host        string
	Url         *url.URL
	Header      map[string]string
	RequestBody interface{}
	Response    interface{}
}

func (ho *HttpOption) Send(ctx context.Context) error {
	if ho.Url == nil {
		return fmt.Errorf("no url specificed")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ho.RequestBody); err != nil {
		return err
	}
	req, err := http.NewRequest(ho.Method, ho.Url.String(), &buf)
	if err != nil {
		return err
	}
	if ho.Host != "" {
		req.Host = ho.Host
	}
	if ho.Header != nil {
		for k, v := range ho.Header {
			req.Header.Set(k, v)
		}
	}

	resp, err := ctxhttp.Do(ctx, &http.Client{}, req)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if ho.Response != nil {
		if err := json.Unmarshal(data, &ho.Response); err != nil {
			return err
		}
	}
	_ = resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}
	return nil
}

func (ho *HttpOption) SendWithCode(ctx context.Context) (int, error) {
	if ho.Url == nil {
		return 0, fmt.Errorf("no url specificed")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ho.RequestBody); err != nil {
		return 0, err
	}
	req, err := http.NewRequest(ho.Method, ho.Url.String(), &buf)
	if err != nil {
		return 0, err
	}

	if ho.Host != "" {
		req.Host = ho.Host
	}
	if ho.Header != nil {
		for k, v := range ho.Header {
			req.Header.Set(k, v)
		}
	}

	resp, err := ctxhttp.Do(ctx, &http.Client{}, req)
	if err != nil {
		return 0, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	if ho.Response != nil {
		if err := json.Unmarshal(data, &ho.Response); err != nil {
			return 0, err
		}
	}

	_ = resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return resp.StatusCode, fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}
	return resp.StatusCode, nil
}

func IsHttps(uri string) bool {
	return strings.Index(strings.ToLower(uri), "https://") == 0
}

func IsHttp(uri string) bool {
	return strings.Index(strings.ToLower(uri), "http://") == 0
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
