package gateway

import (
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransConfig struct {
	ServerKey string
	Env       midtrans.EnvironmentType
}

type MidtransGateway struct {
	snapClient snap.Client
	coreClient coreapi.Client
}

func NewMidtransGateway(cfg MidtransConfig) domain.PaymentGateway {
	var s snap.Client
	s.New(cfg.ServerKey, cfg.Env)

	var c coreapi.Client
	c.New(cfg.ServerKey, cfg.Env)

	return &MidtransGateway{
		snapClient: s,
		coreClient: c,
	}
}

func (g *MidtransGateway) CreatePayment(t *domain.CreatePaymentRequest) (*domain.PaymentResponse, error) {
	// prepare snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  t.OrderID,
			GrossAmt: t.Amount,
		},

		CustomerDetail: &midtrans.CustomerDetails{
			FName: t.Customer.Name,
			Email: t.Customer.Email,
		},

		EnabledPayments: []snap.SnapPaymentType{
			snap.SnapPaymentType(t.PaymentMethod),
		},

		Items: func() *[]midtrans.ItemDetails {
			var items []midtrans.ItemDetails
			for _, item := range t.Items {
				items = append(items, midtrans.ItemDetails{
					Name:  item.Name,
					Qty:   item.Quantity,
					Price: item.Price,
				})
			}
			return &items
		}(),

		Expiry: &snap.ExpiryDetails{
			Unit:     "minute",
			Duration: t.ExpiryMinutes,
		},
	}

	// call midtrans snap API
	snapResp, err := g.snapClient.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return &domain.PaymentResponse{
		Token:      snapResp.Token,
		PaymentURL: snapResp.RedirectURL,
	}, nil
}

func (g *MidtransGateway) CheckStatus(orderID string) (string, error) {
	// call midtrans core API to check transaction status
	res, err := g.coreClient.CheckTransaction(orderID)
	if err != nil {
		return "", err
	}

	// map midtrans status and fraud status to internal transaction status
	finalStatus := pkg.MapMidtransStatus(res.TransactionStatus, res.FraudStatus)

	return finalStatus, nil
}
