package url

import (
	"errors"
	"fmt"
	"strings"
)

// URL represents a parsed URL (technically, a URI reference).
type URL struct {
	// https://foo.com
	Scheme string // https
	Host   string // foo.com
	Path   string // /go
}

// Port returns u.Host's port, without the leading colon.
// If u.Host doesn't contain a port, Port returns an empty string.
func (u *URL) Port() string {
	_, port, _ := split(u.Host, ":", 0)
	return port

}

// Hostname returns u.Host without the port.
func (u *URL) Hostname() string {
	host, _, ok := split(u.Host, ":", 0)
	if !ok {
		return u.Host
	}
	return host
}

func (u *URL) String() string {
	if u == nil {
		return ""
	}
	var sbuilder strings.Builder
	sbuilder.Grow(len(u.Scheme) + len(u.Host) + len(u.Path) + 4)
	if sc := u.Scheme; sc != "" {
		sbuilder.WriteString(sc)
		sbuilder.WriteString("://")
	}
	if h := u.Host; h != "" {
		sbuilder.WriteString(h)
	}
	if p := u.Path; p != "" {
		sbuilder.WriteString("/")
		sbuilder.WriteString(p)
	}

	return sbuilder.String()
}

func (u *URL) testString() string {
	return fmt.Sprintf("scheme=%q, host=%q, path=%q,", u.Scheme, u.Host, u.Path)
}

// Parse parses rawurl into a URL structure.
func Parse(rawurl string) (*URL, error) {
	scheme, rest, ok := parseSchema(rawurl)
	if !ok {
		return nil, errors.New("missing protocol scheme")
	}

	host, path := parseHostPath(rest)
	return &URL{Scheme: scheme, Host: host, Path: path}, nil
}

func parseSchema(rawurl string) (scheme, rest string, ok bool) {
	return split(rawurl, "://", 1)
}

func parseHostPath(hostpath string) (host, path string) {
	host, path, ok := split(hostpath, "/", 0)
	if !ok {
		host = hostpath
	}
	return host, path

}

// split s by sep
//
// split returns empty string if sep is not found in s at index n
func split(s, sep string, n int) (s1, s2 string, ok bool) {
	i := strings.Index(s, sep)
	if i < n {
		return "", "", false
	}
	if len(s) > i {
		s1 = s[:i]
	}

	// s2 = s[i+len(sep):]  boundcheck even when len(s) > i + len(sep)
	return s1, s[i+len(sep):], true
}
