package delivery

import (
	pb "HnH/services/csat/csatPB"
	"HnH/services/csat/internal/repository/psql"
	"HnH/services/csat/internal/usecase"
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
