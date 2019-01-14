package http

import (
	"context"

	"bitbucket.org/ibizsoftware/berith-chain/swarm/api"
	"bitbucket.org/ibizsoftware/berith-chain/swarm/sctx"
)

type uriKey struct{}

func GetRUID(ctx context.Context) string {
	v, ok := ctx.Value(sctx.HTTPRequestIDKey{}).(string)
	if ok {
		return v
	}
	return "xxxxxxxx"
}

func SetRUID(ctx context.Context, ruid string) context.Context {
	return context.WithValue(ctx, sctx.HTTPRequestIDKey{}, ruid)
}

func GetURI(ctx context.Context) *api.URI {
	v, ok := ctx.Value(uriKey{}).(*api.URI)
	if ok {
		return v
	}
	return nil
}

func SetURI(ctx context.Context, uri *api.URI) context.Context {
	return context.WithValue(ctx, uriKey{}, uri)
}
