package service

import (
	"context"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/rexcfnghk/pricing-store/model"
	"github.com/rexcfnghk/pricing-store/repository/currencypair"
	"github.com/rexcfnghk/pricing-store/repository/customer"
	"github.com/rexcfnghk/pricing-store/repository/providercurrencyconfig"
	"github.com/rexcfnghk/pricing-store/repository/quote"
	"github.com/samber/lo"
)

type BestPriceService struct {
	CurrencyPairRepo           *currencypair.RedisRepo
	ProviderCurrencyConfigRepo *providercurrencyconfig.RedisRepo
	CustomerRepo               *customer.RedisRepo
	QuoteRepo                  *quote.RedisRepo
}

func (s *BestPriceService) GetBestPrice(ctx context.Context, currencyPair *model.CurrencyPair, customerId int) (model.BestPrice, error) {
	_, err := s.CustomerRepo.GetById(ctx, customerId)
	if err != nil {
		return model.BestPrice{}, fmt.Errorf("unable to retrieve customer: %w", err)
	}

	currencyPairId, err := s.CurrencyPairRepo.GetByCurrencyPairId(ctx, currencyPair.Base, currencyPair.Quote)
	if err != nil {
		return model.BestPrice{}, fmt.Errorf("unable to retrieve currency pair: %w", err)
	}

	quotes, err := s.QuoteRepo.GetAllByCurrencyPairId(ctx, currencyPairId)
	if err != nil {
		return model.BestPrice{}, fmt.Errorf("unable to retrieve quotes: %w", err)
	}

	providerCurrencyConfigs, err := s.getProviderCurrencyConfigs(ctx, quotes, currencyPairId)
	if err != nil {
		return model.BestPrice{}, fmt.Errorf("unable to retrieve provider currency configs: %w", err)
	}

	fmt.Println(providerCurrencyConfigs)

	filteredQuotes := s.getFilteredQuotes(quotes, providerCurrencyConfigs)

	fmt.Println(filteredQuotes)

	maxBidPriceQuote := lo.MaxBy[model.MarketQuote](filteredQuotes, func(item model.MarketQuote, max model.MarketQuote) bool {
		return item.BidPrice.Compare(max.BidPrice) > 0
	})

	minAskPriceQuote := lo.MinBy[model.MarketQuote](filteredQuotes, func(item model.MarketQuote, min model.MarketQuote) bool {
		return item.AskPrice.Compare(min.AskPrice) < 0
	})

	return model.BestPrice{
		BidPrice:                maxBidPriceQuote.BidPrice,
		BidQuantity:             maxBidPriceQuote.BidQuantity,
		AskPrice:                minAskPriceQuote.AskPrice,
		AskQuantity:             minAskPriceQuote.AskQuantity,
		BestBidMarketProviderId: maxBidPriceQuote.MarketProviderId,
		BestAskMarketProviderId: minAskPriceQuote.MarketProviderId,
	}, nil

	// Get customer rating factor from customer ID
	// Get currency mapping ID from query["base"] and query["quote"]
	// Get all quotes with "quotes:{currencymappingid}"
	// Get all currency configs with quotes.DistinctBy(q => q.MarketProviderId)
	// Filter quotes to only show active based on currency configs
	// BEST PRICE = max bid price and min ask price
}

func (s *BestPriceService) getFilteredQuotes(quotes []model.MarketQuote, providerCurrencyConfigs map[int]bool) []model.MarketQuote {
	var filteredQuotes []model.MarketQuote
	linq.From(quotes).Where(func(q interface{}) bool {
		return providerCurrencyConfigs[q.(model.MarketQuote).MarketProviderId]
	}).ToSlice(&filteredQuotes)

	return filteredQuotes
}

func (s *BestPriceService) getProviderCurrencyConfigs(ctx context.Context, quotes []model.MarketQuote, currencyPairId currencypair.Id) (map[int]bool, error) {
	var uniqueProviderIds []int
	linq.From(quotes).Select(func(q interface{}) interface{} {
		return q.(model.MarketQuote).MarketProviderId
	}).Distinct().ToSlice(&uniqueProviderIds)

	fmt.Println(uniqueProviderIds)

	providerCurrencyConfigs := make(map[int]bool)
	for _, uniqueProviderId := range uniqueProviderIds {
		providerCurrencyConfig, err := s.ProviderCurrencyConfigRepo.GetById(ctx, uniqueProviderId, currencyPairId)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve provider currency config: %w", err)
		}

		providerCurrencyConfigs[uniqueProviderId] = providerCurrencyConfig.IsEnabled
	}

	return providerCurrencyConfigs, nil
}
