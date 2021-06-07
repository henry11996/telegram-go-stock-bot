package main

import (
	"fmt"
	"math/big"

	"github.com/RainrainWu/fugle-realtime-go/client"
	"github.com/shopspring/decimal"
)

func convertInfo(data client.FugleAPIData) string {
	status := ""

	if data.Meta.Issuspended {
		status += "暫停買賣 "
	}
	if data.Meta.Canshortmargin && data.Meta.Canshortlend {
		status += "暫停買賣 "
	}
	if data.Meta.Canshortmargin && data.Meta.Canshortlend {
		status += "可融資券 "
	} else if data.Meta.Canshortmargin {
		status += "禁融券 "
	} else if data.Meta.Canshortlend {
		status += "禁融資 "
	} else {
		status += "禁融資券 "
	}

	if data.Meta.Candaybuysell && data.Meta.Candaysellbuy {
		status += "買賣現沖 "
	} else if data.Meta.Candaybuysell {
		status += "現沖買 "
	} else if data.Meta.Candaysellbuy {
		status += "現沖賣 "
	} else {
		status += "禁現沖 "
	}

	return fmt.Sprintf("[%s\\(%s\\)](https://tw.stock.yahoo.com/q/bc?s=%s)\n"+
		"產業：%s\n"+
		"狀態：%s\n"+
		"現價：%s\n",
		data.Meta.Namezhtw, data.Info.SymbolID, data.Info.SymbolID,
		data.Meta.Industryzhtw,
		status,
		data.Meta.Pricereference,
	)
}

func convertQuote(data client.FugleAPIData) string {
	var status string
	if data.Quote.Istrial {
		status = "試搓中"
	} else if data.Quote.Iscurbingrise {
		status = "緩漲試搓"
	} else if data.Quote.Iscurbingfall {
		status = "緩跌試搓"
	} else if data.Quote.Isclosed {
		//已收盤
		status = ""
	} else if data.Quote.Ishalting {
		status = "暫停交易"
	} else {
		//正常交易
		status = ""
	}

	var currentPirce decimal.Decimal
	zero := decimal.NewFromInt(0)
	if data.Quote.Trade.Price.Equal(zero) {
		currentPirce = data.Quote.Trial.Price
	} else {
		currentPirce = data.Quote.Trade.Price
	}

	var percent, minus *big.Float
	hunded := decimal.NewFromInt(100)
	percent = currentPirce.Sub(data.Meta.Pricereference).Div(data.Meta.Pricereference).Mul(hunded).BigFloat()
	minus = currentPirce.Sub(data.Meta.Pricereference).BigFloat()
	var bestPrices string

	if len(data.Quote.Order.Bestbids) > 0 || len(data.Quote.Order.Bestasks) > 0 {
		for i := 0; i < 5; i++ {
			bidPrice := ""
			askPrice := ""
			bidUnit := ""
			askUnit := ""
			if len(data.Quote.Order.Bestbids) > i {
				bestbids := data.Quote.Order.Bestbids[len(data.Quote.Order.Bestbids)-1-i]
				bidPrice = bestbids.Price.StringFixed(2)
				if bidPrice == "0.00" {
					bidPrice = "市價"
				}
				bidUnit = bestbids.Unit.String()
			}
			for j := 0; j < 5; j++ {
				if len(data.Quote.Order.Bestasks) > j {
					bestasks := data.Quote.Order.Bestasks[j]
					askPrice = bestasks.Price.StringFixed(2)
					if askPrice == "0.00" {
						askPrice = "市價"
					}
					askUnit = bestasks.Unit.String()
				}
				if i == j {
					bestPrices += fmt.Sprintf("%6s %5s \\| %6s %5s\n", bidPrice, bidUnit, askPrice, askUnit)
				}
			}
		}
	} else {
		bestPrices = ""
	}

	return fmt.Sprintf("``` %9s - %s  %s \n"+
		"高 %4v \\| 低 %4v \\| 總 %5v\n"+
		"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
		"            %v         \n"+
		"    買   %2.2f %2.2f%%   賣\n"+
		"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
		"%s```", data.Meta.Namezhtw, data.Info.SymbolID, status,
		data.Quote.PriceHigh.Price, data.Quote.PriceLow.Price, data.Quote.Total.Unit,
		currentPirce.BigFloat(), minus, percent,
		bestPrices,
	)
}

func convertLegalPerson(legalPerson LegalPerson) string {
	return fmt.Sprintf("%s \\- %s\n"+
		"日期：%s\n"+
		"買超股數：%s\n"+
		"賣超股數：%s\n"+
		"買賣超股數：\\%s",
		legalPerson.StockName, legalPerson.StockId,
		legalPerson.Date,
		legalPerson.Buy,
		legalPerson.Sell,
		legalPerson.Total,
	)
}
