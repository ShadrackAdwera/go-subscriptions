package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/subscription"
)

type SubscribePackageTxInput struct {
	CustomerID            string
	SubscriptionPackageID string
}

type SubscribePackageTxOutput struct {
	Message string
}

func (s *Store) SubscribePackageTx(ctx context.Context, args SubscribePackageTxInput) (SubscribePackageTxOutput, error) {
	sKey := os.Getenv("STRIPE_KEY")

	if sKey == "" {
		return SubscribePackageTxOutput{}, fmt.Errorf("no stripe key found")
	}

	err := s.execTx(ctx, func(q *Queries) error {
		stripe.Key = sKey
		// check if customer exists locally
		localCust, err := q.GetSubscriptionUserByStripeId(ctx, args.CustomerID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("user does not exist")
			}
			return err
		}
		// check if package exists locally
		localPackage, err := q.GetPackage(ctx, args.SubscriptionPackageID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("subscription package does not exist")
			}
			return err
		}
		// create subscription in stripe
		sParams := &stripe.SubscriptionParams{
			Customer: stripe.String(localCust.StripeID),
			Items: []*stripe.SubscriptionItemsParams{
				{
					Price: stripe.String(localPackage.StripePriceID),
				},
			},
		}

		s, err := subscription.New(sParams)

		if err != nil {
			return err
		}

		currentStart := time.Unix(s.CurrentPeriodStart, 0).Format(time.RFC3339)
		stTm, _ := time.Parse(time.RFC3339, currentStart)
		currentEnd := time.Unix(s.CurrentPeriodEnd, 0).Format(time.RFC3339)
		endTm, _ := time.Parse(time.RFC3339, currentEnd)

		// create subscription locally - tho sijui kama tunahitaji hii ata - we can fetch directly from stripe
		_, err = q.CreateUserPackage(ctx, CreateUserPackageParams{
			ID:        s.ID,
			UserID:    localCust.ID,
			PackageID: localCust.ID,
			Status:    SubscriptionStatusActive,
			StartDate: stTm,
			EndDate:   endTm,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return SubscribePackageTxOutput{
		Message: "subscription successfully created",
	}, err
}
