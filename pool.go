package npapi

// NumberOfMiners returns the nanopool miners count.
func NumberOfMiners() (uint, error) {
	var size uint
	if err := fetch(&size, activeMinersEndpoint); err != nil {
		return size, err
	}
	return size, nil
}

// NumberOfWorkers returns the nanopool workers count.
func NumberOfWorkers() (uint, error) {
	var size uint
	if err := fetch(&size, activeWorkersEndpoint); err != nil {
		return size, err
	}
	return size, nil
}

// PoolHashrate returns the nanopool hashrate [MH/s].
func PoolHashrate() (float64, error) {
	var hashrate float64
	if err := fetch(&hashrate, poolHashrateEndpoint); err != nil {
		return hashrate, err
	}
	return hashrate, nil
}

// TopMiners returns the top 15 nanopool miners.
func TopMiners() ([]User, error) {
	jsonMiners := []struct {
		Address  string  `json:"address"`
		Hashrate float64 `json:"hashrate"`
	}{}
	if err := fetch(&jsonMiners, topMinersEndpoint); err != nil {
		return nil, err
	}
	miners := make([]User, len(jsonMiners))
	for i, m := range jsonMiners {
		miners[i] = User{
			Address:  m.Address,
			Hashrate: m.Hashrate,
		}
	}
	return miners, nil
}
