package protocol

import (
	"crypto/tls"
	"io"
	"net/url"
	"strings"
)

func FetchUrl(url_addr string, cs *CertStore) (string, error) {
	parsed_url, err := url.Parse(url_addr)
	if err != nil {
		return "", err
	}

	addr := strings.Builder{}
	addr.WriteString(parsed_url.Hostname())
	if p := parsed_url.Port(); p != "" {
		addr.WriteString(":")
		addr.WriteString(p)
	} else {
		addr.WriteString(":1965")
	}

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", addr.String(), config)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	valid, err := cs.CheckCertificate(parsed_url.Hostname(), cert)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", ErrInvalidCert
	}

	_, err = conn.Write([]byte(url_addr))
	if err != nil {
		return "", err
	}
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}
	_, err = io.Copy(&sb, conn)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}
