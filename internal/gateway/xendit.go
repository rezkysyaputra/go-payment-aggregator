package gateway

import (
	"context"
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg"

	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
)

type XenditConfig struct {
	ApiKey string
}

type XenditGateway struct {
	xenditClient *xendit.APIClient
}

func NewXenditGateway(cfg XenditConfig) domain.PaymentGateway {
	c := xendit.NewClient(cfg.ApiKey)

	return &XenditGateway{
		xenditClient: c,
	}

}

func mapPaymentMethodToXendit(method string) []string {
	mapping := map[string][]string{
		"credit_card": {"CREDIT_CARD"},
		"bank_transfer": {
			"BNI",
			"BCA",
			"MANDIRI",
			"PERMATA",
			"BRI",
		},
		"e_wallet": {"OVO", "DANA", "SHOPEEPAY", "LINKAJA"},
		"qris":     {"QRIS"},
	}

	if methods, ok := mapping[method]; ok {
		return methods
	}
	return []string{}
}

func (x *XenditGateway) CreatePayment(req *domain.CreatePaymentRequest) (*domain.PaymentResponse, error) {

	invoiceDurationSeconds := float32(req.ExpiryMinutes * 60)

	reqInvoice := invoice.CreateInvoiceRequest{
		PayerEmail: &req.Customer.Email,
		ExternalId: req.OrderID,
		Amount:     float64(req.Amount),
		Currency:   &req.Currency,
		Customer: &invoice.CustomerObject{
			GivenNames: *invoice.NewNullableString(&req.Customer.Name),
			Email:      *invoice.NewNullableString(&req.Customer.Email),
		},
		Items: func() []invoice.InvoiceItem {
			var items []invoice.InvoiceItem
			for _, item := range req.Items {
				items = append(items, invoice.InvoiceItem{
					Name:     item.Name,
					Quantity: float32(item.Quantity),
					Price:    float32(item.Price),
				})
			}
			return items
		}(),

		PaymentMethods: mapPaymentMethodToXendit(req.PaymentMethod),

		InvoiceDuration: &invoiceDurationSeconds,
	}

	ctx := context.Background()
	inv, _, err := x.xenditClient.InvoiceApi.CreateInvoice(ctx).CreateInvoiceRequest(reqInvoice).Execute()
	if err != nil {
		return nil, err
	}

	return &domain.PaymentResponse{
		PaymentURL: inv.InvoiceUrl,
		Token:      *inv.Id,
	}, nil
}

func (x *XenditGateway) CheckStatus(orderID string) (string, error) {
	ctx := context.Background()
	inv, _, err := x.xenditClient.InvoiceApi.GetInvoiceById(ctx, orderID).Execute()
	if err != nil {
		return "", err
	}

	status := pkg.MapXenditStatus(inv.Status.String())

	return status, nil
}
