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
		AddonBreakdown: breakdowns,
		AddonTotal:     q.AddonTotal,
		Subtotal:       q.Subtotal,
		TotalPrice:     q.TotalPrice,
	}
}

func (a DbQuoteAddon) ToAddOnBreakdownProto() *payment.AddOnBreakdown {
	return &payment.AddOnBreakdown{
		AddonId:   a.AddonID,
		AddonName: a.AddonName,
		Price:     a.AddonPrice,
	}
}
