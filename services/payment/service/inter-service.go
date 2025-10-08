package service

import "context"

func (p *PaymentService) HandleSubscriptions(ctx context.Context) error {
	if err := p.SubscribeBookingRequests(); err != nil {
		p.L.Error("%v\n", err)
		return err
	}
	<-ctx.Done()
	return ctx.Err()

}

func (p *PaymentService) SubscribeBookingRequests() error {
	return nil
}
