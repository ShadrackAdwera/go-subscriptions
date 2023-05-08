package db

import (
	"context"
	"fmt"
	"os"

	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"
)

type CreatePackageTxInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	//AfterCreate func()error
}

type CreatePackageTxResult struct {
	Message string `json:"message"`
}

func (s *Store) CreatePackageTx(ctx context.Context, args CreatePackageTxInput) (CreatePackageTxResult, error) {
	sKey := os.Getenv("STRIPE_KEY")

	if sKey == "" {
		return CreatePackageTxResult{}, fmt.Errorf("no stripe key found")
	}

	err := s.execTx(ctx, func(q *Queries) error {
		//create a product in stripe
		stripe.Key = sKey
		productParams := &stripe.ProductParams{
			Name:        stripe.String(args.Name),
			Description: stripe.String(args.Description),
		}
		pd, err := product.New(productParams)

		if err != nil {
			return err
		}
		// create a price in stripe
		priceParams := &stripe.PriceParams{
			Currency: stripe.String(string(stripe.CurrencyUSD)),
			Product:  stripe.String(pd.ID),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String("month"),
			},
			UnitAmount: stripe.Int64(args.Price),
		}
		prodPrice, err := price.New(priceParams)

		if err != nil {
			return err
		}

		// create package locally using price ID
		_, err = q.CreatePackage(ctx, CreatePackageParams{
			ID:            pd.ID,
			Name:          pd.Name,
			Description:   pd.Description,
			Price:         prodPrice.UnitAmount,
			StripePriceID: prodPrice.ID,
		})

		if err != nil {
			return err
		}
		return nil
	})

	return CreatePackageTxResult{
		Message: "Subscription Package Successfully Created",
	}, err
}
