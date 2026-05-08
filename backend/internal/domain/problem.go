package domain

type Problem struct {
	ID                 int32  `yaml:"id"`
	Title              string `yaml:"title"`
	SetupSQL           string `yaml:"setup_sql"`
	ExpectedResultJSON string `yaml:"expected_result_json"`
	AnswerSQL          string `yaml:"answer_sql"`
	IsOrderMatters     bool   `yaml:"is_order_matters"`
	Description        string `yaml:"description"`
}
