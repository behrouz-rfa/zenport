package grpc

import (
	"context"
	"google.golang.org/grpc"
	"zenport/gates/internal/domain"
	"zenport/ntps/ntpspb"
)

type NtpRepository struct {
	client ntpspb.TimeServiceClient
}

var _ domain.NtpRepository = (*NtpRepository)(nil)

func NewNtpRepository(conn *grpc.ClientConn) NtpRepository {
	return NtpRepository{client: ntpspb.NewTimeServiceClient(conn)}
}

func (r NtpRepository) FetchTime(ctx context.Context, request string) (string, error) {
	t, err := r.client.GetTime(ctx, &ntpspb.GetTimeRequest{Time: request})
	if err != nil {
		return "", err
	}
	return t.Time, nil
}
