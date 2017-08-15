package npapi

import (
	"strconv"
	"time"
)

// These constants are used for identifying hashrate intervals.
const (
	// One hour interval
	OneHour Interval = "h1"
	// Three hour interval
	ThreeHours = "h3"
	// Six hour interval
	SixHours = "h6"
	// Twelve hour interval
	TwelveHours = "h12"
	// Twentyfour hour interval
	TwentyfourHours = "h24"
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

// Payment is a nanopool.org payment.
type Payment struct {
	Date      Time
	TxHash    string
	Amount    float64
	Confirmed bool
}

// Worker is a nanopool.org worker. It represents one mining machine.
type Worker struct {
	ID               string
	Hashrate         float64
	LastShare        Time
	Rating           uint
	AverageHashrates map[Interval]float64
}

type HashrateItem struct {
	ID       string
	Hashrate float64
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

// ShareItem stores share history metrics.
type ShareItem struct {
	Date   Time
	Shares uint
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

// HashrateAndBalance retrieves the current hashrate and balance.
func HashrateAndBalance(addr string) (float64, float64, error) {
	data := struct {
		Hashrate float64 `json:"hashrate"`
		Balance  float64 `json:"balance"`
	}{}
	if err := fetch(&data, balanceHashrateEndpoint, addr); err != nil {
		return data.Hashrate, data.Balance, err
	}
	return data.Hashrate, data.Balance, nil
}

// ReportedHashrate retrieves the last reported hashrate.
func ReportedHashrate(addr string) (float64, error) {
	var hashrate float64
	if err := fetch(&hashrate, reportedHashrateEndpoint, addr); err != nil {
		return hashrate, err
	}
	return hashrate, nil
}

// Workers retrieves a list of workers bound to this account.
func Workers(addr string) ([]Worker, error) {
	jsonWorkers := []struct {
		ID        string  `json:"id"`
		Hashrate  float64 `json:"hashrate"`
		LastShare Time    `json:"lastShare"`
		Rating    uint    `json:"rating"`
	}{}
	if err := fetch(&jsonWorkers, workersEndpoint, addr); err != nil {
		return nil, err
	}
	workers := make([]Worker, len(jsonWorkers))
	for i, w := range jsonWorkers {
		workers[i] = Worker{
			ID:               w.ID,
			Hashrate:         w.Hashrate,
			LastShare:        w.LastShare,
			Rating:           w.Rating,
			AverageHashrates: map[Interval]float64{},
		}
	}
	return workers, nil
}

// Payments retrieves a list of occured payments from nanopool to the user.
func Payments(addr string) ([]Payment, error) {
	jsonPayments := []struct {
		Date      Time    `json:"date"`
		TxHash    string  `json:"txhash"`
		Amount    float64 `json:"amount"`
		Confirmed bool    `json:"confirmed"`
	}{}
	if err := fetch(&jsonPayments, paymentsEndpoint, addr); err != nil {
		return nil, err
	}
	payments := make([]Payment, len(jsonPayments))
	for i, p := range jsonPayments {
		payments[i] = Payment(p)
	}
	return payments, nil
}

// ShareHistory retrieves a history of share rate metrics.
func ShareHistory(addr string) ([]ShareItem, error) {
	jsonHistory := []struct {
		Date   Time `json:"date"`
		Shares uint `json:"shares"`
	}{}
	if err := fetch(&jsonHistory, sharerateHistoryEndpoint, addr); err != nil {
		return nil, err
	}
	history := make([]ShareItem, len(jsonHistory))
	for i, s := range jsonHistory {
		history[i] = ShareItem(s)
	}
	return history, nil
}

// WorkersAverageHashrateIn retrieves a list of workers, each associated with its hashrate in the given interval.
func WorkersAverageHashrateIn(addr string, interval uint) ([]HashrateItem, error) {
	jsonWorkers := []struct {
		ID       string  `json:"worker"`
		Hashrate float64 `json:"hashrate"`
	}{}
	if err := fetch(&jsonWorkers, workersAverageHashrateLimitedEndpoint, addr, interval); err != nil {
		return nil, err
	}
	workers := make([]HashrateItem, len(jsonWorkers))
	for i, w := range jsonWorkers {
		workers[i] = HashrateItem(w)
	}
	return workers, nil
}

// WorkerAverageHashrate retrieves a list of workers, each associated with its hashrates.
func WorkersAverageHashrate(addr string) (map[Interval][]HashrateItem, error) {
	jsonIntervals := map[Interval][]struct {
		ID       string  `json:"worker"`
		Hashrate float64 `json:"hashrate"`
	}{}
	if err := fetch(&jsonIntervals, workersAverageHashrateEndpoint, addr); err != nil {
		return nil, err
	}
	intervals := make(map[Interval][]HashrateItem)
	for key, jsonWorkers := range jsonIntervals {
		workers := make([]HashrateItem, len(jsonWorkers))
		for i, w := range jsonWorkers {
			workers[i] = HashrateItem(w)
		}
		intervals[key] = workers
	}
	return intervals, nil
}

// WorkerReportedHashrate retrieves the last reported hashrate associated with each worker.
func WorkerReportedHashrate(addr string) ([]HashrateItem, error) {
	jsonWorkers := []struct {
		ID       string  `json:"worker"`
		Hashrate float64 `json:"hashrate"`
	}{}
	if err := fetch(&jsonWorkers, workersReportedHashrateEndpoint, addr); err != nil {
		return nil, err
	}
	workers := make([]HashrateItem, len(jsonWorkers))
	for i, w := range jsonWorkers {
		workers[i] = HashrateItem(w)
	}
	return workers, nil
}
