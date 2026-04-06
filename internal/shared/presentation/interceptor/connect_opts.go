// Package interceptor
package interceptor

import "connectrpc.com/connect"

func CommonHandlerOptions(interceptors ...connect.Interceptor) []connect.HandlerOption {
	return []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
	}
}
