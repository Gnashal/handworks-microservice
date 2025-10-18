package helpers

import (
	model "handworks-gateway/graph/generated/models"
	"handworks/common/grpc/payment"
)

func MapQuote(res *payment.QuoteResponse) *model.Quote {
	if res == nil {
		return nil
	}

	addons := make([]*model.AddOnBreakdown, len(res.AddonBreakdown))
	for i, a := range res.AddonBreakdown {
		addons[i] = &model.AddOnBreakdown{
			AddonID:   a.AddonId,
			AddonName: a.AddonName,
			Price:     float64(a.Price),
		}
	}

	return &model.Quote{
		QuoteID:          res.QuoteId,
		MainServiceName:  res.MainServiceName,
		MainServiceTotal: float64(res.MainServiceTotal),
		AddonBreakdown:   addons,
		AddonTotal:       float64(res.AddonTotal),
		TotalPrice:       float64(res.TotalPrice),
	}
}
