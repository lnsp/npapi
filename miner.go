package npapi

import (
	"strconv"
	"time"
)

const (
	OneHour         Interval = "h1"
	ThreeHours               = "h3"
	SixHours                 = "h6"
	TwelveHours              = "h12"
	TwentyfourHours          = "h24"
)

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	secs, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*t = Time(time.Unix(secs, 0))
	return nil
}

type Interval string

// Worker is a nanopool.org worker. It represents one mining machine.
type Worker struct {
	ID               string
	Hashrate         float64
	LastShare        Time
	AverageHashrates map[Interval]float64
}

// User is a nanopool.org user identified by his address. A user can have multiple workers.
type User struct {
	Address            string
	Balance            float64
	UnconfirmedBalance float64
	Hashrate           float64
	AverageHashrates   map[Interval]float64
	Worker             []Worker
}

// ChartItem stores hashrate metrics of a specific point in time.
type ChartItem struct {
	Date     Time
	Shares   uint
	Hashrate float64
}

// HistoryItem stores hashrate history metrics.
type HistoryItem struct {
	Date     Time
	Hashrate float64
}

// UserInfo retrieves a complete set of user information including workers and hashrate statistics.
func UserInfo(addr string) (*User, error) {
	var user struct {
		Balance            string              `json:"balance"`
		UnconfirmedBalance string              `json:"unconfirmed_balance"`
		Hashrate           string              `json:"hashrate"`
		AverageHashrates   map[Interval]string `json:"avghashrate"`
		Worker             []struct {
			ID                 string `json:"id"`
			Hashrate           string `json:"hashrate"`
			LastShare          Time   `json:"lastShare"`
			AvgOneHour         string `json:"avg_h1"`
			AvgThreeHours      string `json:"avg_h3"`
			AvgSixHours        string `json:"avg_h6"`
			AvgTwelveHours     string `json:"avg_h12"`
			AvgTwentyfourHours string `json:"avg_h24"`
		} `json:"worker"`
	}
	if err := fetch(&user, userEndpoint, addr); err != nil {
		return nil, err
	}
	worker := make([]Worker, len(user.Worker))
	for i, w := range user.Worker {
		worker[i] = Worker{
			ID:        w.ID,
			Hashrate:  mustf(w.Hashrate),
			LastShare: w.LastShare,
			AverageHashrates: map[Interval]float64{
				OneHour:         mustf(w.AvgOneHour),
				ThreeHours:      mustf(w.AvgThreeHours),
				SixHours:        mustf(w.AvgSixHours),
				TwelveHours:     mustf(w.AvgTwelveHours),
				TwentyfourHours: mustf(w.AvgTwentyfourHours),
			},
		}
	}
	return &User{
		Address:            addr,
		Balance:            mustf(user.Balance),
		UnconfirmedBalance: mustf(user.UnconfirmedBalance),
		Hashrate:           mustf(user.Hashrate),
		AverageHashrates: map[Interval]float64{
			OneHour:         mustf(user.AverageHashrates[OneHour]),
			ThreeHours:      mustf(user.AverageHashrates[ThreeHours]),
			SixHours:        mustf(user.AverageHashrates[SixHours]),
			TwelveHours:     mustf(user.AverageHashrates[TwelveHours]),
			TwentyfourHours: mustf(user.AverageHashrates[TwentyfourHours]),
		},
		Worker: worker,
	}, nil
}

type jsonBalance struct {
	Status  bool    `json:"status"`
	Balance float64 `json:"data"`
}

// Balance retrieves the accounts balance.
func Balance(addr string) (float64, error) {
	var balance float64
	if err := fetch(&balance, accountBalanceEndpoint, addr); err != nil {
		return balance, err
	}
	return balance, nil
}

// AverageHashrateIn retrieves the average hashrate in the last x hours.
func AverageHashrateIn(addr string, hours uint) (float64, error) {
	var hashrate float64
	if err := fetch(&hashrate, averageHashrateLimitedEndpoint, addr, hours); err != nil {
		return hashrate, err
	}
	return hashrate, nil
}

// AverageHashrate retrieves the average hashrate in the last one to twentyfour hours.
func AverageHashrate(addr string) (map[Interval]float64, error) {
	avgs := map[Interval]float64{}
	if err := fetch(&avgs, averageHashrateEndpoint, addr); err != nil {
		return nil, err
	}
	return avgs, nil
}

// HashrateChart retrieves the hashrate chart data.
func HashrateChart(addr string) ([]ChartItem, error) {
	jsonItems := []struct {
		Date     Time    `json:"date"`
		Shares   uint    `json:"shares"`
		Hashrate float64 `json:"hashrate"`
	}{}
	if err := fetch(&jsonItems, hashrateChartEndpoint, addr); err != nil {
		return nil, err
	}
	items := make([]ChartItem, len(jsonItems))
	for i := range jsonItems {
		items[i] = ChartItem(jsonItems[i])
	}
	return items, nil
}

// Exists checks if the account exists.
func Exists(addr string) error {
	var data string
	if err := fetch(&data, accountExistEndpoint, addr); err != nil {
		return err
	}
	return nil
}

// CurrentHashrate retrieves the current calculated hashrate.
func CurrentHashrate(addr string) (float64, error) {
	var hashrate float64
	if err := fetch(&hashrate, currentHashrateEndpoint, addr); err != nil {
		return hashrate, err
	}
	return hashrate, nil
}

// HashrateHistory fetches the latest hashrate history.
func HashrateHistory(addr string) ([]HistoryItem, error) {
	jsonHistory := []struct {
		Date     Time    `json:"date"`
		Hashrate float64 `json:"hashrate"`
	}{}
	if err := fetch(&jsonHistory, historyEndpoint, addr); err != nil {
		return nil, err
	}
	history := make([]HistoryItem, len(jsonHistory))
	for i := range jsonHistory {
		history[i] = HistoryItem(jsonHistory[i])
	}
	return history, nil
}
