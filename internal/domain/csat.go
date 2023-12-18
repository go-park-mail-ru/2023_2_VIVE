package domain

//easyjson:json
type StarsNum struct {
	StarsNum int32 `json:"starsNum,omitempty"`
	Count    int64 `json:"count,omitempty"`
}

//easyjson:json
type StarsNumSlice []StarsNum

//easyjson:json
type QuestionComment struct {
	Comment string `json:"comment,omitempty"`
}

//easyjson:json
type QuestionCommentSlice []QuestionComment

//easyjson:json
type QuestionStatistics struct {
	AvgStars            float32              `json:"avgStars,omitempty"`
	StarsNumList        StarsNumSlice        `json:"starsNumList,omitempty"`
	QuestionCommentList QuestionCommentSlice `json:"questionCommentList,omitempty"`
	QuestionText        string               `json:"question_text,omitempty"`
}

//easyjson:json
type QuestionStatisticsSlice []QuestionStatistics

//easyjson:json
type Statistics struct {
	StatisticsList QuestionStatisticsSlice `json:"statisticsList,omitempty"`
}

//easyjson:json
type Question struct {
	Question   string `json:"question,omitempty"`
	Name       string `json:"name,omitempty"`
	QuestionId int64  `json:"question_id,omitempty"`
}

//easyjson:json
type QuestionSlice []Question

//easyjson:json
type QuestionList struct {
	Questions QuestionSlice `json:"questions,omitempty"`
}

//easyjson:json
type Answer struct {
	Starts     int32  `json:"starts,omitempty"`
	Comment    string `json:"comment,omitempty"`
	QuestionId int64  `json:"question_id,omitempty"`
}
