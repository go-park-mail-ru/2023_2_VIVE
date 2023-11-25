package delivery

import (
	pb "HnH/services/csat/csatPB"
	"HnH/services/csat/internal/repository/psql"
	"HnH/services/csat/internal/usecase"
	"context"
	"database/sql"
)

type CsatServer struct {
	pb.UnimplementedCsatServer
	csatUscase usecase.ICsatUsecase
}

func NewServer(db *sql.DB) (*CsatServer, error) {
	csatRepo := psql.NewPsqlCsatRepository(db)
	csatUsecase := usecase.NewCsatUscase(csatRepo)

	return &CsatServer{
		csatUscase: csatUsecase,
	}, nil
}

func (s *CsatServer) GetQuestions(ctx context.Context, userID *pb.UserID) (*pb.QuestionList, error) {
	return s.csatUscase.GetQuestions(ctx, userID)
}

func (s *CsatServer) GetStatistic(ctx context.Context, empty *pb.Empty) (*pb.Statistics, error) {
	return s.csatUscase.GetStatistics(ctx)
}

func (s *CsatServer) RegisterAnswer(ctx context.Context, answer *pb.Answer) (*pb.Empty, error) {
	return s.csatUscase.RegisterAnswer(ctx, answer)
}
