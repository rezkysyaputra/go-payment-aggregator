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

func (g *MidtransGateway) CreatePayment(req *domain.CreatePaymentRequest) (*domain.PaymentResponse, error) {
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.Amount,
		},

		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.Customer.Name,
			Email: req.Customer.Email,
		},

		EnabledPayments: []snap.SnapPaymentType{
			snap.SnapPaymentType(req.PaymentMethod),
		},

		Items: func() *[]midtrans.ItemDetails {
			var items []midtrans.ItemDetails
			for _, item := range req.Items {
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
			Duration: req.ExpiryMinutes,
		},
	}

	snapResp, err := g.snapClient.CreateTransaction(snapReq)
	if err != nil {
		return nil, err
	}

	return &domain.PaymentResponse{
		Token:      snapResp.Token,
		PaymentURL: snapResp.RedirectURL,
	}, nil
}

func (g *MidtransGateway) CheckStatus(orderID string) (string, error) {
	res, err := g.coreClient.CheckTransaction(orderID)
	if err != nil {
		return "", err
	}

	finalStatus := pkg.MapMidtransStatus(res.TransactionStatus, res.FraudStatus)

	return finalStatus, nil
}
