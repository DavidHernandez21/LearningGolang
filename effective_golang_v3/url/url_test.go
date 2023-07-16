package url

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	const rawurl = "http://host/some/fake/path"

	want := &URL{
		Scheme: "http",
		Host:   "foo.com",
		Path:   "go",
	}

	got, err := Parse(rawurl)
	if err != nil {
		t.Fatalf("Parse(%q) err = %q, want nil", rawurl, err)
	}
	if *got != *want {
		t.Fatalf("Parse(%q):\n\tgot: %s\n\twant: %s\n", rawurl, got.testString(), want.testString())
	}

}

func TestParseInvalidURL(t *testing.T) {
	tests := map[string]string{
		"missing scheme": "foo:com",
		"empty scheme":   "://foo.com",
	}

	for name, in := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse(in); err == nil {

				t.Errorf("Parse(%q)=nil; want an error", in)
			}
		})
	}
}

func TestURLPort(t *testing.T) {
	tests := map[string]struct {
		in       string
		hostname string
		port     string
	}{
		"with port":          {in: "foo.com:80", hostname: "foo.com", port: "80"},     // with port
		"with empty port":    {in: "foo.com:", hostname: "foo.com", port: ""},         // with empty port
		"without port":       {in: "foo.com", hostname: "foo.com", port: ""},          // without port
		"ip with port":       {in: "1.2.3.4:8080", hostname: "1.2.3.4", port: "8080"}, // ip with port
		"ip with empty port": {in: "1.2.3.4:", hostname: "1.2.3.4", port: ""},         // ip with empty port
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("Hostname/%s/%s", name, tt.in), func(t *testing.T) {
			u := &URL{Host: tt.in}
			if got, want := u.Hostname(), tt.hostname; got != want {
				t.Errorf("for host %q; got %q; want %q", tt.in, got, want)

			}
		})

		t.Run(fmt.Sprintf("Port/%s/%s", name, tt.in), func(t *testing.T) {
			u := &URL{Host: tt.in}
			if got, want := u.Port(), tt.port; got != want {
				t.Errorf("for host %q; got %q; want %q", tt.in, got, want)

			}
		})
	}

}

func TestURLString(t *testing.T) {

	tests := map[string]struct {
		url  *URL
		want string
	}{
		"nil URL":               {url: nil, want: ""},
		"empty URL":             {url: &URL{}, want: ""},
		"scheme":                {url: &URL{Scheme: "https"}, want: "https://"},
		"scheme and host":       {url: &URL{Scheme: "https", Host: "foo.com"}, want: "https://foo.com"},
		"schame, host and path": {url: &URL{Scheme: "https", Host: "foo.com", Path: "go"}, want: "https://foo.com/go"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.url.String(); got != tt.want {
				t.Errorf("url:%#v\ngot: %q\nwant: %q", tt.url, got, tt.want)
			}
		})
	}
}

func BenchmarkURLString(b *testing.B) {
	tests := []struct {
		name string // test name
		url  *URL
	}{
		{name: "nil URL", url: nil},
		{name: "empty URL", url: &URL{}},
		{name: "scheme", url: &URL{Scheme: "https"}},
		{name: "scheme and host", url: &URL{Scheme: "https", Host: "foo.com"}},
		{name: "schame, host and path", url: &URL{Scheme: "https", Host: "foo.com", Path: "go"}},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.url.String()
			}
		})
	}

}

// func TestURLPort(t *testing.T) {
// 	tests := []struct {
// 		name string // test name
// 		in   string // input host i.e. URL.Host
// 		port string
// 	}{
// 		{name: "with port", in: "foo.com:80", port: "80"},        // with port
// 		{name: "with empty port", in: "foo.com:", port: ""},      // with empty port
// 		{name: "without port", in: "foo.com", port: ""},          // without port
// 		{name: "ip with port", in: "1.2.3.4:8080", port: "8080"}, // ip with port
// 		{name: "ip with empty port", in: "1.2.3.4:", port: ""},   // ip with empty port
// 	}

// 	for _, tt := range tests {
// 		u := &URL{Host: tt.in}
// 		if got, want := u.Port(), tt.port; got != want {
// 			t.Errorf("%s: for host %q; got %q; want %q", tt.name, tt.in, got, want)

// 		}
// 	}

// }
