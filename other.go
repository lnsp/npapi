package npapi

// EarningsItem stores information about the projected earnings for a specified interval.
type EarningsItem struct {
	// Earned coins in ETH
	Coins float64
	// Earned bitcoins
	Bitcoins float64
	// Earned USD
	Dollars float64
	// Earned CNY
	Yuan float64
	// Earned EUR
	Euros float64
	// Earned RUR
	Rubles float64
}

// EarningsReport stores approximated earnings for the intervals minute, hour, day, week and month.
type EarningsReport struct {
	PerMinute, PerHour, PerDay, PerWeek, PerMonth EarningsItem
}

// PriceReport stores the price information.
type PriceReport struct {
	// ETH price in USD
	USDollar float64
	// ETH price in EUR
	Euro float64
	// ETH price in RUR
	Rubles float64
	// ETH price in CNY
	Yuan float64
	// ETH price in BTC
	Bitcoins float64
}

// ApproximatedEarnings calculates the approximated earnings projected by the hashrate.
func ApproximatedEarnings(hashrate float64) (EarningsReport, error) {
	jsonReport := map[string]struct {
		Coins    float64 `json:"coins"`
		Bitcoins float64 `json:"bitcoins"`
		Dollars  float64 `json:"dollars"`
		Yuan     float64 `json:"yuan"`
		Euros    float64 `json:"euros"`
		Rubles   float64 `json:"rubles"`
	}{}
	if err := fetch(&jsonReport, approximatedEarningsEndpoint, hashrate); err != nil {
		return EarningsReport{}, err
	}
	return EarningsReport{
		PerMinute: EarningsItem(jsonReport["minute"]),
		PerHour:   EarningsItem(jsonReport["hour"]),
		PerDay:    EarningsItem(jsonReport["day"]),
		PerWeek:   EarningsItem(jsonReport["week"]),
		PerMonth:  EarningsItem(jsonReport["month"]),
	}, nil
}

// Prices fetches a price report from the server, storing the current exchange rates for ETH.
func Prices() (PriceReport, error) {
	jsonPrices := struct {
		USDollar float64 `json:"price_usd"`
		Euro     float64 `json:"price_eur"`
		Rubles   float64 `json:"price_rur"`
		Yuan     float64 `json:"price_cny"`
		Bitcoins float64 `json:"price_btc"`
	}{}
	if err := fetch(&jsonPrices, pricesEndpoint); err != nil {
		return PriceReport{}, err
	}
	return PriceReport(jsonPrices), nil
}
