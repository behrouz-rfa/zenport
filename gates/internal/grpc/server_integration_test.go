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
	"zenport/gates/gatespb"
	"zenport/gates/internal/application"
	"zenport/gates/internal/domain"
)

type serverSuite struct {
	mocks struct {
		time *domain.MockTimeRepository
		ntp  *domain.MockNtpRepository
	}
	server *grpc.Server
	client gatespb.GatesServiceClient
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
		time *domain.MockTimeRepository
		ntp  *domain.MockNtpRepository
	}{
		time: domain.NewMockTimeRepository(s.T()),
		ntp:  domain.NewMockNtpRepository(s.T()),
	}

	// create app
	app := application.NewApplication(s.mocks.time, s.mocks.ntp)

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
	s.client = gatespb.NewGatesServiceClient(conn)
}
func (s *serverSuite) TearDownTest() {
	s.server.GracefulStop()
}

func (s *serverSuite) TestAskTime_WithQuestion() {

	ctx := context.Background()
	s.mocks.ntp.On("FetchTime", mock.Anything, mock.AnythingOfType("string")).Return(string("What time is it?"), nil)
	_, err := s.client.GetTime(ctx, &gatespb.GetTimeRequest{Ask: "What time is it?"})

	s.Assert().NoError(err)
}
