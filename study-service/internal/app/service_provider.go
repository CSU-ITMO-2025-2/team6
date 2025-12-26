package app

import (
	"context"
	"log"

	db "local-lib/database"
	"local-lib/database/pg"
	"local-lib/queue"
	"study-service/internal/api/study"
	"study-service/internal/closer"
	"study-service/internal/config"
	"study-service/internal/repository"
	storage "study-service/internal/repository/minio"
	studyRepository "study-service/internal/repository/study"
	"study-service/internal/service"
	studyService "study-service/internal/service/study"
	studyAnalyzeWorkerService "study-service/internal/service/study_analyze_worker"
	"study-service/internal/workers"

	"go.uber.org/zap"
)

type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	minioConfig config.MinIOConfig

	dbClient db.Client
	queue    queue.Study

	studyRepository repository.StudyRepository
	studyService    service.StudyService
	studyImpl       *study.Implementation

	storage repository.Storage

	studyAnalyzeWorkerService service.StudyAnalyzeWorkerService
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

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) MinioConfig() config.MinIOConfig {
	if s.minioConfig == nil {
		cfg, err := config.NewMinIOConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.minioConfig = cfg
	}

	return s.minioConfig
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

func (s *serviceProvider) Queue() queue.Study {
	if s.queue == nil {
		nConf := queue.InitFromEnv()
		n := queue.Connect(nConf)
		s.queue = queue.WrapNatsCore(n)
		s.queue.AddMLStream()
	}

	return s.queue
}

func (s *serviceProvider) StartWorkers(ctx context.Context) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	workers.StartWorkers(ctx, sugar, s.Queue(), s.StudyAnalyzeWorkerService(ctx))
}

func (s *serviceProvider) StudyRepository(ctx context.Context) repository.StudyRepository {
	if s.studyRepository == nil {
		s.studyRepository = studyRepository.NewRepository(s.DBClient(ctx))
	}

	return s.studyRepository
}

func (s *serviceProvider) Storage() repository.Storage {
	if s.storage == nil {
		s.storage = storage.NewStorage(
			s.MinioConfig().Host(),
			s.MinioConfig().User(),
			s.MinioConfig().Password(),
			s.MinioConfig().UseSSL())
	}

	return s.storage
}

func (s *serviceProvider) StudyService(ctx context.Context) service.StudyService {
	if s.studyService == nil {
		s.studyService = studyService.NewService(
			s.StudyRepository(ctx),
			s.Storage(),
			s.Queue(),
		)
	}

	return s.studyService
}

func (s *serviceProvider) StudyAnalyzeWorkerService(ctx context.Context) service.StudyAnalyzeWorkerService {
	if s.studyAnalyzeWorkerService == nil {
		s.studyAnalyzeWorkerService = studyAnalyzeWorkerService.NewService(
			s.StudyRepository(ctx),
			s.Storage(),
			s.Queue(),
		)
	}

	return s.studyAnalyzeWorkerService
}

func (s *serviceProvider) StudyImpl(ctx context.Context) *study.Implementation {
	if s.studyImpl == nil {
		s.studyImpl = study.NewImplementation(s.StudyService(ctx))
	}

	return s.studyImpl
}
