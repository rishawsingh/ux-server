package score

import (
	"math"

	"github.com/samber/lo"

	"github.com/remotestate/golang/models"
)

func calculateProductScore(attributeScore []models.SurveyAnswerAttributeWeight,
	productList []models.ProductAttributeList,
	genderAttribute []string) map[string]int {
	shouldCheckGenderAttribute := len(genderAttribute) > 0
	attributeScoreMap := make(map[string]int)
	for i := range attributeScore {
		attributeScoreMap[attributeScore[i].AttributeID] = attributeScore[i].Weight
	}
	productScore := make(map[string]int)
	for i := range productList {
		var numerator float64
		var denominator float64
		percentage := 100.0
		if shouldCheckGenderAttribute {
			ok := lo.Some[string](productList[i].AttributeIDs, genderAttribute)
			if !ok {
				// this product does not have any attributes for gender
				continue
			}
		}
		for j := range productList[i].AttributeIDs {
			val, ok := attributeScoreMap[productList[i].AttributeIDs[j]]
			if ok {
				numerator += float64(val)
				denominator += float64(val)
			} else {
				denominator++
			}
		}
		weight := int(math.Round((numerator / denominator) * percentage))
		if weight > 0 {
			productScore[productList[i].Product.ID] = weight
		}
	}
	return productScore
}
