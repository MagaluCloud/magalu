package schema

import (
	"net/url"
	"path"
	"strings"

	"github.com/invopop/jsonschema"
)

type URI string

func (u URI) JSONSchemaExtend(s *jsonschema.Schema) {
	s.Format = "uri"
}

func (u URI) String() string {
	return string(u)
}

func (u URI) JoinPath(parts ...string) URI {
	if parsed, err := url.Parse(string(u)); err == nil {
		return URI(parsed.JoinPath(parts...).String())
	}

	toBeJoined := make([]string, 0, 1+len(parts))
	toBeJoined = append(toBeJoined, string(u))
	toBeJoined = append(toBeJoined, parts...)
	return URI(path.Join(toBeJoined...))
}

func (u URI) Hostname() string {
	if parsed, err := url.Parse(string(u)); err == nil {
		if hostname := parsed.Hostname(); hostname != "" {
			return hostname
		}
		if parsed.Path != "" {
			return strings.Split(parsed.Path, "/")[0]
		}
	}
	return u.String()
}

func (u URI) Path() string {
	if parsed, err := url.Parse(string(u)); err == nil {
		var path string
		if parsed.Hostname() != "" {
			path = parsed.Path
		} else if pathEntries := strings.Split(parsed.Path, "/"); len(pathEntries) > 1 {
			path = strings.Join(pathEntries[1:], "/")
		} else {
			return ""
		}
		return strings.Trim(path, "/")
	}
	return u.String()
}

func (u URI) AsFilePath() FilePath {
	return FilePath(u.Path())
}

func (u URI) AsDirPath() DirPath {
	return DirPath(u.Path())
}

func (u URI) Scheme() string {
	if parsed, err := url.Parse(string(u)); err == nil && parsed.Scheme != "" {
		return parsed.Scheme
	}
	return ""
}

func (u URI) IsRoot() bool {
	path := u.Path()
	return path == "" || path == "/"
}

func (u URI) Filename() string {
	filename := path.Base(u.String())
	if filename == "/" || filename == "." {
		return ""
	}
	return filename
}

type FilePath string

func (p FilePath) JSONSchemaExtend(s *jsonschema.Schema) {
	s.ContentMediaType = "application/octet-stream"
}

func (p FilePath) AsURI() URI {
	u := url.URL{}
	u.Scheme = "path"
	u.Path = string(p)
	return URI(u.String())
}

func (p FilePath) Join(parts ...string) FilePath {
	toBeJoined := make([]string, 0, 1+len(parts))
	toBeJoined = append(toBeJoined, string(p))
	toBeJoined = append(toBeJoined, parts...)
	return FilePath(path.Join(toBeJoined...))
}

func (p FilePath) String() string {
	return string(p)
}

type DirPath string

func (p DirPath) JSONSchemaExtend(s *jsonschema.Schema) {
	s.ContentMediaType = "inode/directory"
}

func (p DirPath) AsURI() URI {
	u := url.URL{}
	u.Scheme = "path"
	u.Path = string(p)
	return URI(u.String())
}

func (p DirPath) Join(parts ...string) DirPath {
	toBeJoined := make([]string, 0, 1+len(parts))
	toBeJoined = append(toBeJoined, string(p))
	toBeJoined = append(toBeJoined, parts...)
	return DirPath(path.Join(toBeJoined...))
}

func (p DirPath) String() string {
	return string(p)
}
