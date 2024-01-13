package connection

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrConnectionAddressEmpty = errors.New("connection address is empty")
	ErrSystemCertPoolEmpty    = errors.New("system cert pool is empty")
)

func GRPCDialContext(ctx context.Context, address string, insecureConnection bool, extraOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if address == "" {
		return nil, fmt.Errorf("%w", ErrConnectionAddressEmpty)
	}

	// prepend default connection options so that can be overwritten
	var opts []grpc.DialOption
	if insecureConnection {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		rootCAs, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		} else if rootCAs == nil {
			return nil, fmt.Errorf("%w", ErrSystemCertPoolEmpty)
		}
		cfg := &tls.Config{RootCAs: rootCAs}
		creds := credentials.NewTLS(cfg)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	opts = append(opts, extraOpts...)

	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
