package url

import (
	"fmt"
	"net/url"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/option"
)

func ParseIncomingRequest(in types.IncomingRequest) (*url.URL, error) {
	scheme, ok := option.Handle(in.Scheme())
	if !ok {
		return nil, fmt.Errorf("scheme is required")
	}

	authority, ok := option.Handle(in.Authority())
	if !ok {
		return nil, fmt.Errorf("authority is required")
	}

	path := option.UnwrapOr(in.PathWithQuery(), "/")

	rawURL := fmt.Sprintf("%s://%s%s",
		scheme.String(),
		authority,
		path,
	)

	return url.ParseRequestURI(rawURL)
}

func SetOutgoingRequestURL(out types.OutgoingRequest, u *url.URL) error {
	out.SetScheme(cm.Some(toScheme(u.Scheme)))

	out.SetAuthority(cm.Some(u.Host))
	pathWithQuery := u.Path
	if u.RawQuery != "" {
		pathWithQuery += "?" + u.RawQuery
	}
	out.SetPathWithQuery(cm.Some(pathWithQuery))
	return nil
}

func toScheme(scheme string) types.Scheme {
	switch scheme {
	case "http":
		return types.SchemeHTTP()
	case "https":
		return types.SchemeHTTPS()
	default:
		return types.SchemeOther(scheme)
	}
}
