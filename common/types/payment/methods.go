package types

import (
	"handworks/common/grpc/payment"
)

func (q DbQuote) ToProto() *payment.QuoteResponse {
	breakdowns := make([]*payment.AddOnBreakdown, 0, len(q.Addons))
	for _, a := range q.Addons {
		breakdowns = append(breakdowns, a.ToAddOnBreakdownProto())
	}

	return &payment.QuoteResponse{
		AddonBreakdown:   breakdowns,
		AddonTotal:       q.AddonTotal,
		MainServiceName:  q.MainService,
		MainServiceTotal: q.Subtotal,
		TotalPrice:       q.TotalPrice,
		QuoteId:          q.ID,
	}
}

func (a DbQuoteAddon) ToAddOnBreakdownProto() *payment.AddOnBreakdown {
	return &payment.AddOnBreakdown{
		AddonId:   a.ID,
		AddonName: a.ServiceType,
		Price:     a.AddonPrice,
	}
}
