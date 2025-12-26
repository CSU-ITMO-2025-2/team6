package app

import (
	"context"
	"log"

	db "local-lib/database"
	"local-lib/database/pg"
	"main-service/internal/api/study"
	"main-service/internal/api/user"
	"main-service/internal/client"
	"main-service/internal/client/grpc/study_client"
	"main-service/internal/closer"
	"main-service/internal/config"
	"main-service/internal/repository"
	userRepository "main-service/internal/repository/user"
	"main-service/internal/service"
	studyService "main-service/internal/service/study"
	userService "main-service/internal/service/user"
	pb "study-service/pkg/study_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	pgConfig           config.PGConfig
	httpConfig         config.HTTPConfig
	studyServiceConfig config.GRPCConfig

	dbClient    db.Client
	studyClient client.StudyClient

	userRepository repository.UserRepository
	userService    service.UserService
	userImpl       *user.Implementation

	studyService service.StudyService
	studyImpl    *study.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) StudyServiceConfig() config.GRPCConfig {
	if s.studyServiceConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.studyServiceConfig = cfg
	}

	return s.studyServiceConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) StudyClient() client.StudyClient {
	if s.studyClient == nil {
		conn, err := grpc.NewClient(
			s.StudyServiceConfig().Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(),
		)
		if err != nil {
			log.Fatalf("error initializing grpc client: %v", err)
		}

		grpcClient := pb.NewStudyServiceClient(conn)

		s.studyClient = study_client.NewStudyClient(grpcClient)
	}

	return s.studyClient
}

func (s *serviceProvider) StudyService(ctx context.Context) service.StudyService {
	if s.studyService == nil {
		s.userService = studyService.NewService(
			s.StudyClient(),
		)
	}

	return s.studyService
}

func (s *serviceProvider) StudyImpl(ctx context.Context) *study.Implementation {
	if s.studyImpl == nil {
		s.studyImpl = study.NewImplementation(s.StudyService(ctx))
	}

	return s.studyImpl
}
