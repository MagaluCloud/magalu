package openapi

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

type UrlHost string
type PathHost string
type SchemaHost string

type UrlHostKey struct{}
type PathHostKey struct{}
type SchemaHostKey struct{}

func GetUrlHostFlag(ctx context.Context) string {
	host, found := ctx.Value(UrlHostKey{}).(UrlHost)
	if !found {
		return ""
	}
	return string(host)
}

func GetPathHostFlag(ctx context.Context) string {
	path, found := ctx.Value(PathHostKey{}).(PathHost)
	if !found {
		return ""
	}
	return string(path)
}

func GetSchemaHostFlag(ctx context.Context) string {
	schema, found := ctx.Value(SchemaHostKey{}).(SchemaHost)
	if !found {
		return ""
	}
	return string(schema)
}

func WithUrlHostFlag(ctx context.Context, host string) context.Context {
	_, host, _, err := parseURL(host)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, UrlHostKey{}, UrlHost(host))
}

func WithPathHostFlag(ctx context.Context, path string) context.Context {
	_, _, path, err := parseURL(path)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, PathHostKey{}, PathHost(path))
}

func WithSchemaHostFlag(ctx context.Context, schema string) context.Context {
	schema, _, _, err := parseURL(schema)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, SchemaHostKey{}, SchemaHost(schema))
}

func parseURL(input string) (schema string, host string, path string, err error) {
	if input == "" {
		return "", "", "", errors.New("empty URL")
	}
	if input == "http" || input == "https" {
		return input, "", "", nil
	}

	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "http://" + input
	}

	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", "", "", err
	}

	// Extrai as partes requisitadas
	schema = parsedURL.Scheme
	host = parsedURL.Host
	path = parsedURL.Path

	return schema, host, path, nil
}
