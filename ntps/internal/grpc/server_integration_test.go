//go:build integration

package grpc

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"zenport/ntps/internal/application"
	"zenport/ntps/internal/domain"
	"zenport/ntps/ntpspb"
)

type serverSuite struct {
	mocks struct {
		ntp  *domain.MockNtpRepository
		time *domain.MockTimeRepository
	}
	server *grpc.Server
	client ntpspb.TimeServiceClient
	suite.Suite
}

func TestServer(t *testing.T) {
	suite.Run(t, &serverSuite{})
}

func (s *serverSuite) SetupSuite()    {}
func (s *serverSuite) TearDownSuite() {}

func (s *serverSuite) SetupTest() {
	const grpcTestPort = ":10912"

	var err error
	// create server
	s.server = grpc.NewServer()
	var listener net.Listener
	listener, err = net.Listen("tcp", grpcTestPort)
	if err != nil {
		s.T().Fatal(err)
	}

	// create mocks
	s.mocks = struct {
		ntp  *domain.MockNtpRepository
		time *domain.MockTimeRepository
	}{
		time: domain.NewMockTimeRepository(s.T()),
		ntp:  domain.NewMockNtpRepository(s.T()),
	}

	// create app
	app := application.New(s.mocks.time)

	// register app with server
	if err = RegisterServer(app, s.server); err != nil {
		s.T().Fatal(err)
	}
	go func(listener net.Listener) {
		err := s.server.Serve(listener)
		if err != nil {
			s.T().Fatal(err)
		}
	}(listener)

	// create client
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(grpcTestPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.T().Fatal(err)
	}
	s.client = ntpspb.NewTimeServiceClient(conn)
}
func (s *serverSuite) TearDownTest() {
	s.server.GracefulStop()
}

func (s *serverSuite) TestAskTime_Send() {
	ctx := context.Background()
	s.mocks.time.On("Save", mock.Anything, mock.AnythingOfType("*domain.Time")).Return(nil)

	_, err := s.client.GetTime(ctx, &ntpspb.GetTimeRequest{Time: "What time is it?"})

	s.Assert().NoError(err)
}
