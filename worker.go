package npapi

// WorkerAverageHashrate fetches the hashrate of a worker in the specified time interval.
func WorkerAverageHashrateIn(addr, worker string, hours uint) (float64, error) {
	var hashrate float64
	if err := fetch(&hashrate, workerAverageHashrateLimitedEndpoint, addr, worker, hours); err != nil {
		return hashrate, err
	}
	return hashrate, nil
}

// WorkerAverageHashrate fetches a collection of average hashrates in different intervals.
func WorkerAverageHashrate(addr, worker string) (HashrateReport, error) {
	jsonHashrates := make(map[string]float64)
	if err := fetch(&jsonHashrates, workerAverageHashrateEndpoint, addr, worker); err != nil {
		return HashrateReport{}, err
	}
	return toHashrateReport(jsonHashrates), nil
}

// WorkerHashrateChart retrieves a hashrate chart specific for the given worker.
func WorkerHashrateChart(addr, worker string) ([]ChartItem, error) {
	jsonChart := []struct {
		Date     Time    `json:"date"`
		Shares   uint    `json:"shares"`
		Hashrate float64 `json:"hashrate"`
	}{}
	if err := fetch(&jsonChart, workerHashrateChartEndpoint, addr, worker); err != nil {
		return nil, err
	}
	chart := make([]ChartItem, len(jsonChart))
	for i, c := range jsonChart {
		chart[i] = ChartItem(c)
	}
	return chart, nil
}

// WorkerCurrentHashrate fetches the current worker hashrate [MH/s].
func WorkerCurrentHashrate(addr, worker string) (float64, error) {
	var hashrate float64
	if err := fetch(&hashrate, workerCurrentHashrateEndpoint, addr, worker); err != nil {
		return hashrate, err
	}
	return hashrate, nil
}

// WorkerHashrateHistory fetches records of hashrates for this specific worker.
func WorkerHashrateHistory(addr, worker string) ([]HistoryItem, error) {
	jsonHistory := []struct {
		Date     Time    `json:"date"`
		Hashrate float64 `json:"hashrate"`
	}{}
	if err := fetch(&jsonHistory, workerHistoryEndpoint, addr, worker); err != nil {
		return nil, err
	}
	history := make([]HistoryItem, len(jsonHistory))
	for i, h := range jsonHistory {
		history[i] = HistoryItem(h)
	}
	return history, nil
}

// WorkerReportedHashrate fetches the hashrate reported by the worker.
func WorkerReportedHashrate(addr, worker string) (float64, error) {
	var hashrate float64
	if err := fetch(&hashrate, workerReportedHashrateEndpoint, addr, worker); err != nil {
		return hashrate, err
	}
	return hashrate, nil
}

// WorkerShareHistory fetches the workers share history.
func WorkerShareHistory(addr, worker string) ([]ShareItem, error) {
	jsonShares := []struct {
		Date   Time `json:"date"`
		Shares uint `json:"shares"`
	}{}
	if err := fetch(&jsonShares, workerShareRateHistoryEndpoint, addr, worker); err != nil {
		return nil, err
	}
	shares := make([]ShareItem, len(jsonShares))
	for i, s := range jsonShares {
		shares[i] = ShareItem(s)
	}
	return shares, nil
}
