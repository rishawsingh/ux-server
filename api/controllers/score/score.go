package score

import (
	"net/http"
	"sort"
	"time"

	"github.com/remotestate/golang/services"

	"github.com/remotestate/golang/utils"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/remotestate/golang/internal"
)

const maxProductsForSubcategory = 4
const maxScoredProducts = 50

type Controller struct {
	logger         *internal.Logger
	productService services.Product
	surveyService  services.Survey
	inMemScore     *utils.InMemScore
	tx             *internal.Transactor
	tracer         *internal.Tracer
}

func NewController(logger *internal.Logger,
	productService services.Product,
	surveyService services.Survey,
	inMemScore *utils.InMemScore,
	tx *internal.Transactor,
	tracer *internal.Tracer) *Controller {
	return &Controller{
		logger:         logger,
		productService: productService,
		surveyService:  surveyService,
		inMemScore:     inMemScore,
		tx:             tx,
		tracer:         tracer,
	}
}

func (c *Controller) Calculate(ctx *gin.Context) {
	surveyInviteID := ctx.Param("surveyID")
	if _, parseErr := uuid.Parse(surveyInviteID); parseErr != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, parseErr)
		return
	}
	start := time.Now()
	defer func(start time.Time, c *Controller, ctx *gin.Context) {
		diff := time.Since(start).Milliseconds()
		c.logger.WithCtx(ctx.Request.Context()).Infof("Calculate took %d millisecond", diff)
	}(start, c, ctx)
	/*
		 	a sample of tx
			c.tx.Wrap(ctx.Request.Context(), func(tx *sqlx.Tx) error {
				_, err := c.productService.GetAllProductWithAttributes(ctx.Request.Context(), tx)
				return err
			})
	*/
	if c.inMemScore.Len() == 0 {
		productWithAttributes, scoreErr := c.productService.GetAllProductWithAttributes(ctx.Request.Context(), nil)
		if scoreErr != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, errors.Errorf("%s, stack: %+v",
				"failed to get all product with attributes", scoreErr))
			return
		}
		c.inMemScore.Set(productWithAttributes)
	}
	genderAttribute, genderErr := c.surveyService.FindGenderAttribute(ctx.Request.Context(), nil, surveyInviteID)
	if genderErr != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, errors.Errorf("%s, stack: %+v",
			"failed to get gender attribute for survey", genderErr))
		return
	}
	surveyAttributeWeight, surveyAttributeErr := c.surveyService.FindAttributeWeightForUser(ctx.Request.Context(), nil, surveyInviteID)
	if surveyAttributeErr != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, errors.Errorf("%s, stack: %+v",
			"failed to get all product with attributes", surveyAttributeErr))
		return
	}
	finalWeight := calculateProductScore(surveyAttributeWeight, c.inMemScore.Get(), genderAttribute)
	products := c.inMemScore.GetSegregatedData()
	type responseData struct {
		SubCategoryID string `json:"subcategoryId"`
		BrandID       string `json:"brandId"`
		ProductID     string `json:"productId"`
		Weight        int    `json:"weight"`
	}
	response := make([]responseData, 0)
	for subCategoryID, brands := range products {
		subcategoryProducts := make([]responseData, 0)
		for brandID, brandProducts := range brands {
			topScoreForThisBrand, index := 0, 0
			for i := range brandProducts {
				if weight, ok := finalWeight[brandProducts[i].ID]; ok {
					if topScoreForThisBrand < weight {
						topScoreForThisBrand = weight
						index = i
					}
				}
			}
			if topScoreForThisBrand > 0 {
				subcategoryProducts = append(subcategoryProducts, responseData{
					SubCategoryID: subCategoryID,
					BrandID:       brandID,
					ProductID:     brandProducts[index].ID,
					Weight:        topScoreForThisBrand,
				})
			}
		}
		// sort
		sort.SliceStable(subcategoryProducts, func(i, j int) bool {
			return subcategoryProducts[i].Weight > subcategoryProducts[j].Weight
		})
		if len(subcategoryProducts) > maxProductsForSubcategory {
			subcategoryProducts = subcategoryProducts[0:maxProductsForSubcategory]
		}
		if len(subcategoryProducts) > 0 {
			response = append(response, subcategoryProducts...)
		}
	}
	sort.SliceStable(response, func(i, j int) bool {
		return response[i].Weight > response[j].Weight
	})
	if len(response) > maxScoredProducts {
		response = response[0:maxScoredProducts]
	}
	ctx.JSON(http.StatusOK, response)
}
