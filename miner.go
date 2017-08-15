package npapi

import (
	"strconv"
	"time"
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

// Payment is a nanopool.org payment.
type Payment struct {
	// Payment date
	Date Time
	// Payment transaction hash
	TxHash string
	// Payment amount
	Amount float64
	// Payment status
	Confirmed bool
}

// Worker is a nanopool.org worker. It represents one mining machine.
type Worker struct {
	// Worker ID
	ID string
	// Worker Hashrate [MH/s]
	Hashrate float64
	// Last Share date of Worker
	LastShare Time
	// Worker Rating
	Rating uint
	// Average hashrates
	AverageHashrates HashrateReport
}

// HashrateItem stores an association between a worker and an (averaged) hashrate.
type HashrateItem struct {
	// Worker ID
	ID string
	// Worker Hashrate [MH/s]
	Hashrate float64
}

// HashrateReport storing the (average) hashrates in the last one, six, three, twelve and twentyfour hours.
type HashrateReport struct {
	LastHour, LastThreeHours, LastSixHours, LastTwelveHours, LastDay float64
}

// toHashrateReport parses a hashrate map to a well defined report.
func toHashrateReport(data map[string]float64) HashrateReport {
	return HashrateReport{
		LastHour:        data["h1"],
		LastThreeHours:  data["h3"],
		LastSixHours:    data["h6"],
		LastTwelveHours: data["h12"],
		LastDay:         data["h24"],
	}
}

// WorkerHashrateReport is a nanopool worker report.
type WorkerHashrateReport struct {
	LastHour, LastThreeHours, LastSixHours, LastTwelveHours, LastDay []HashrateItem
}

// User is a nanopool.org user identified by his address. A user can have multiple workers.
type User struct {
	// Account address
	Address string
	// Account balance
	Balance float64
	// Account unconfirmed balance
	UnconfirmedBalance float64
	// Account hashrate [MH/s]
	Hashrate float64
	// Average hashrate [MH/s]
	AverageHashrates HashrateReport
	// Workers
	Workers []Worker
}

// ChartItem stores hashrate metrics of a specific point in time.
type ChartItem struct {
	// Date
	Date Time
	// Number of shares for last 10 minutes
	Shares uint
	// Miner reported hashrate [MH/s]
	Hashrate float64
}

// HistoryItem stores hashrate history metrics.
type HistoryItem struct {
	// Item date
	Date Time
	// Miner hashrate [MH/s]
	Hashrate float64
}

// ShareItem stores share history metrics.
type ShareItem struct {
	// Item date
	Date Time
	// Number of shares for last 10 minutes
	Shares uint
}

// json balance struct
type jsonBalance struct {
	Status  bool    `json:"status"`
	Balance float64 `json:"data"`
}

// json worker hasrate
type jsonWorkerHashrate struct {
	ID       string  `json:"worker"`
	Hashrate float64 `json:"hashrate"`
}

// UserInfo retrieves a complete set of user information including workers and hashrate statistics.
func UserInfo(addr string) (*User, error) {
	var user struct {
		Balance            string            `json:"balance"`
		UnconfirmedBalance string            `json:"unconfirmed_balance"`
		Hashrate           string            `json:"hashrate"`
		AverageHashrates   map[string]string `json:"avghashrate"`
		Workers            []struct {
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
	workers := make([]Worker, len(user.Workers))
	for i, w := range user.Workers {
		averageHashratesMap, err := parseStringMapToFloat(map[string]string{
			"h1":  w.AvgOneHour,
			"h3":  w.AvgThreeHours,
			"h6":  w.AvgSixHours,
			"h12": w.AvgTwelveHours,
			"h24": w.AvgTwentyfourHours,
		})
		if err != nil {
			return nil, err
		}
		currentHashrate, err := strconv.ParseFloat(w.Hashrate, 64)
		if err != nil {
			return nil, err
		}
		workers[i] = Worker{
			ID:               w.ID,
			Hashrate:         currentHashrate,
			LastShare:        w.LastShare,
			AverageHashrates: toHashrateReport(averageHashratesMap),
		}
	}

	averageHashratesMap, err := parseStringMapToFloat(user.AverageHashrates)
	if err != nil {
		return nil, err
	}
	balance, err := strconv.ParseFloat(user.Balance, 64)
	if err != nil {
		return nil, err
	}
	unconfirmedBalance, err := strconv.ParseFloat(user.UnconfirmedBalance, 64)
	if err != nil {
		return nil, err
	}
	currentHashrate, err := strconv.ParseFloat(user.Hashrate, 64)
	if err != nil {
		return nil, err
	}
	return &User{
		Address:            addr,
		Balance:            balance,
		UnconfirmedBalance: unconfirmedBalance,
		Hashrate:           currentHashrate,
		AverageHashrates:   toHashrateReport(averageHashratesMap),
		Workers:            workers,
	}, nil
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
func AverageHashrate(addr string) (HashrateReport, error) {
	avgs := map[string]float64{}
	if err := fetch(&avgs, averageHashrateEndpoint, addr); err != nil {
		return HashrateReport{}, err
	}
	return toHashrateReport(avgs), nil
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
			ID:        w.ID,
			Hashrate:  w.Hashrate,
			LastShare: w.LastShare,
			Rating:    w.Rating,
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
	jsonWorkers := []jsonWorkerHashrate{}
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
func WorkersAverageHashrate(addr string) (WorkerHashrateReport, error) {
	toHashrateItemList := func(jsonWorkers []jsonWorkerHashrate) []HashrateItem {
		workers := make([]HashrateItem, len(jsonWorkers))
		for i, w := range jsonWorkers {
			workers[i] = HashrateItem(w)
		}
		return workers
	}
	jsonIntervals := map[string][]jsonWorkerHashrate{}
	if err := fetch(&jsonIntervals, workersAverageHashrateEndpoint, addr); err != nil {
		return WorkerHashrateReport{}, err
	}
	return WorkerHashrateReport{
		LastHour:        toHashrateItemList(jsonIntervals["h1"]),
		LastThreeHours:  toHashrateItemList(jsonIntervals["h3"]),
		LastSixHours:    toHashrateItemList(jsonIntervals["h6"]),
		LastTwelveHours: toHashrateItemList(jsonIntervals["h12"]),
		LastDay:         toHashrateItemList(jsonIntervals["h24"]),
	}, nil
}

// WorkersReportedHashrate retrieves the last reported hashrate associated with each worker.
func WorkersReportedHashrate(addr string) ([]HashrateItem, error) {
	jsonWorkers := []jsonWorkerHashrate{}
	if err := fetch(&jsonWorkers, workersReportedHashrateEndpoint, addr); err != nil {
		return nil, err
	}
	workers := make([]HashrateItem, len(jsonWorkers))
	for i, w := range jsonWorkers {
		workers[i] = HashrateItem(w)
	}
	return workers, nil
}
