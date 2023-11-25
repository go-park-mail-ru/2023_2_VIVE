package domain

import pb "HnH/services/csat/csatPB"

type QuestionStatistics struct {
	AvgStars            float32
	StarsNumList        []pb.StarsNum
	QuestionCommentList []string
}
