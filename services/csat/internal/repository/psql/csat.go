package psql

import (
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	pb "HnH/services/csat/csatPB"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ICsatRepository interface {
	RegisterRequestTime(ctx context.Context, userID int64) error
	GetLastUpdate(ctx context.Context, userID int64) (time.Time, error)
	GetQuestions(ctx context.Context) ([]*pb.Question, error)
	RegisterAnswer(ctx context.Context, answer *pb.Answer) error
	GetStatistics(ctx context.Context) (*pb.Statistics, error)
}

type psqlCsatRepository struct {
	DB *sql.DB
}

func NewPsqlCsatRepository(db *sql.DB) ICsatRepository {
	return &psqlCsatRepository{
		DB: db,
	}
}

func (repo *psqlCsatRepository) RegisterRequestTime(ctx context.Context, userID int64) error {
	var exists bool

	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("writing last request time by 'user_id' in postgres")

	err := repo.DB.QueryRow(`SELECT EXISTS (SELECT user_id FROM csat_data.user_info WHERE user_id = $1)`, userID).Scan(&exists)
	if err != nil {
		contextLogger.WithField("err", err.Error()).Info("error while checking existence of last request time")
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	if exists {
		err = repo.DB.QueryRow(`UPDATE csat_data.user_info SET "last_request_at" = now() WHERE user_id = $1`, userID).Err()
		if err != nil {
			contextLogger.WithField("err", err.Error()).Info("error while updating last request time")
			return serverErrors.INTERNAL_SERVER_ERROR
		}
	} else {
		err = repo.DB.QueryRow(`INSERT INTO csat_data.user_info ("user_id") VALUES ($1)`, userID).Err()
		if err != nil {
			contextLogger.WithField("err", err.Error()).Info("error while inserting new 'user_id' in csat_data.user_info")
			return serverErrors.INTERNAL_SERVER_ERROR
		}
	}

	return nil
}

func (repo *psqlCsatRepository) GetLastUpdate(ctx context.Context, userID int64) (time.Time, error) {
	var reqTime time.Time

	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting last request time by 'user_id' in postgres")

	err := repo.DB.QueryRow(`SELECT last_request_at FROM csat_data.user_info WHERE user_id = $1`, userID).Scan(&reqTime)
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, serverErrors.ErrNoLastUpdate
	} else if err != nil {
		contextLogger.WithField("err", err.Error()).Info("error while getting last request time")
		return time.Time{}, serverErrors.INTERNAL_SERVER_ERROR
	}

	return reqTime, nil
}

func (repo *psqlCsatRepository) GetQuestions(ctx context.Context) ([]*pb.Question, error) {
	fmt.Printf("start GetQuestions")
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting csat questions from postgres")

	rows, err := repo.DB.Query(`SELECT id, "name", "text" FROM csat_data.question`)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, serverErrors.ErrQuestionsNotFound
	}
	if err != nil {
		contextLogger.WithField("err", err.Error()).Info("error while getting csat questions")
		return nil, serverErrors.INTERNAL_SERVER_ERROR
	}
	defer rows.Close()

	questionsToReturn := make([]*pb.Question, 0)

	for rows.Next() {
		question := &pb.Question{}

		err := rows.Scan(&question.QuestionId, &question.Name, &question.Question)
		if err != nil {
			contextLogger.WithField("err", err.Error()).Info("error while scanning questions from DB")
			return nil, serverErrors.INTERNAL_SERVER_ERROR
		}
		questionsToReturn = append(questionsToReturn, question)
	}

	// contextLogger.Debugf("%d", len(questionsToReturn))
	fmt.Printf("%d", len(questionsToReturn))

	if len(questionsToReturn) == 0 {
		return nil, serverErrors.ErrQuestionsNotFound
	}

	return questionsToReturn, nil
}

func (repo *psqlCsatRepository) RegisterAnswer(ctx context.Context, answer *pb.Answer) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("writing csat answer to postgres")

	var mess *string
	if answer.Comment != "" {
		mess = new(string)
		*mess = answer.Comment
	} else {
		mess = nil
	}

	err := repo.DB.QueryRow(`INSERT INTO csat_data.answer ("stars", "message", "question_id") VALUES ($1, $2, $3)`,
		answer.Starts, mess, answer.QuestionId).Err()
	if err != nil {
		contextLogger.WithField("err", err.Error()).Info("error while inserting new answer into DB")
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	return nil
}

func (repo *psqlCsatRepository) GetStatistics(ctx context.Context) (*pb.Statistics, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("collecting csat statistics from postgres")

	QuestRows, err := repo.DB.Query(`SELECT id, "text" FROM csat_data.question`)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, serverErrors.ErrQuestionsNotFound
	}
	if err != nil {
		contextLogger.WithField("err", err.Error()).Info("error while getting csat questions")
		return nil, serverErrors.INTERNAL_SERVER_ERROR
	}
	defer QuestRows.Close()

	questionStats := make([]*pb.QuestionStatistics, 0)

	for QuestRows.Next() {
		var quest_id int
		questionStat := &pb.QuestionStatistics{}

		err := QuestRows.Scan(&quest_id, &questionStat.QuestionText)
		if err != nil {
			contextLogger.WithField("err", err.Error()).Info("error while scanning questions from DB")
			return nil, serverErrors.INTERNAL_SERVER_ERROR
		}

		AnsRows, err := repo.DB.Query(`SELECT stars, message FROM csat_data.answer WHERE question_id = $1`, quest_id)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serverErrors.ErrAnswerNotFound
		}
		if err != nil {
			contextLogger.WithField("err", err.Error()).Info("error while getting answers from DB")
			return nil, serverErrors.INTERNAL_SERVER_ERROR
		}
		defer AnsRows.Close()

		messages := make([]*pb.QuestionComment, 0)
		countOfStars := map[int32]int64{}

		for AnsRows.Next() {
			var stars int32
			var message *string

			err = AnsRows.Scan(&stars, &message)
			if err != nil {
				contextLogger.WithField("err", err.Error()).Info("error while scanning answers from DB")
				return nil, serverErrors.INTERNAL_SERVER_ERROR
			}

			countOfStars[stars]++

			if message != nil {
				mess := &pb.QuestionComment{
					Comment: *message,
				}
				messages = append(messages, mess)
			}
		}

		var starsSum int64
		var starsCount int64
		marks := make([]*pb.StarsNum, 0, 5)
		for stars, count := range countOfStars {
			starsSum += int64(stars) * count
			starsCount += count

			starsNum := &pb.StarsNum{
				StarsNum: stars,
				Count:    count,
			}

			marks = append(marks, starsNum)
		}

		questionStat.QuestionCommentList = messages
		questionStat.StarsNumList = marks
		questionStat.AvgStars = float32(starsSum) / float32(starsCount)

		questionStats = append(questionStats, questionStat)
	}

	if len(questionStats) == 0 {
		return nil, serverErrors.ErrEntityNotFound
	}

	result := &pb.Statistics{
		StatisticsList: questionStats,
	}

	return result, nil
}
