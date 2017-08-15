package npapi

// BlockStatItem is block metric measuring the difficulty and block time.
type BlockStatItem struct {
	// Block date
	Date Time
	// Block difficulty
	Difficulty uint64
	// Block time
	BlockTime float64
}

// BlockItem is a nanopool.org object representing a generated block.
type BlockItem struct {
	// Block number
	Number uint
	// Block hash
	Hash string
	// Block date
	Date Time
	// Block difficulty
	Difficulty uint64
	// Block miner address
	Miner string
}

// AverageBlocktime fetches the average time needed to create a block.
func AverageBlocktime() (float64, error) {
	var blocktime float64
	if err := fetch(&blocktime, averageBlocktimeEndpoint); err != nil {
		return blocktime, err
	}
	return blocktime, nil
}

// BlockStats fetches the blocks stats for the given block interval.
func BlockStats(offset, count uint) ([]BlockStatItem, error) {
	jsonStats := []struct {
		Date       Time    `json:"date"`
		Difficulty uint64  `json:"difficulty"`
		BlockTime  float64 `json:"block_time"`
	}{}
	if err := fetch(&jsonStats, blockStatsEndpoint, offset, count); err != nil {
		return nil, err
	}
	stats := make([]BlockStatItem, len(jsonStats))
	for i, s := range jsonStats {
		stats[i] = BlockStatItem(s)
	}
	return stats, nil
}

// Blocks fetches the latest blocks provided by the nanopool network.
func Blocks(offset, count uint) ([]BlockItem, error) {
	jsonBlocks := []struct {
		Number     uint   `json:"number"`
		Hash       string `json:"hash"`
		Date       Time   `json:"date"`
		Difficulty uint64 `json:"difficulty"`
		Miner      string `json:"miner"`
	}{}
	if err := fetch(&jsonBlocks, blockEndpoint, offset, count); err != nil {
		return nil, err
	}
	blocks := make([]BlockItem, len(jsonBlocks))
	for i, b := range jsonBlocks {
		blocks[i] = BlockItem(b)
	}
	return blocks, nil
}

// LastBlockNumber fetches the latest block number.
func LastBlockNumber() (uint, error) {
	var number uint
	if err := fetch(&number, lastBlockNumberEndpoint); err != nil {
		return number, err
	}
	return number, nil
}

// TimeToNextEpoch returns the time in seconds until the next epoch.
func TimeToNextEpoch() (float64, error) {
	var seconds float64
	if err := fetch(&seconds, timeToNextEpochEndpoint); err != nil {
		return seconds, err
	}
	return seconds, nil
}
