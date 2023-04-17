package models

type SurveyAnswerAttributeWeight struct {
	Weight        int    `json:"weight" db:"weight"`
	AttributeName string `json:"attributeName" db:"name"`
	AttributeID   string `json:"attributeId" db:"id"`
}
