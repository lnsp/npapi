// Package npapi provides a lightweight wrapper for the Nanopool Ethereum API.
//
// See https://eth.nanopool.org/api for more information.
package npapi

const (
	apiAddress                            = "https://api.nanopool.org/v1/eth"
	accountBalanceEndpoint                = "%s/balance/%s"
	averageHashrateLimitedEndpoint        = "%s/avghashratelimited/%s/%d"
	averageHashrateEndpoint               = "%s/avghashrate/%s"
	hashrateChartEndpoint                 = "%s/hashratechart/%s"
	accountExistEndpoint                  = "%s/accountexist/%s"
	currentHashrateEndpoint               = "%s/hashrate/%s"
	userEndpoint                          = "%s/user/%s"
	historyEndpoint                       = "%s/history/%s"
	balanceHashrateEndpoint               = "%s/balance_hashrate/%s"
	reportedHashrateEndpoint              = "%s/reportedhashrate/%s"
	workersEndpoint                       = "%s/workers/%s"
	paymentsEndpoint                      = "%s/payments/%s"
	sharerateHistoryEndpoint              = "%s/shareratehistory/%s"
	workersAverageHashrateLimitedEndpoint = "%s/avghashrateworkers/%s/%d"
	workersAverageHashrateEndpoint        = "%s/avghashrateworkers/%s"
	workersReportedHashrateEndpoint       = "%s/reportedhashrates/%s"
	averageBlocktimeEndpoint              = "%s/network/avgblocktime"
	blockStatsEndpoint                    = "%s/block_stats/%d/%d"
	blockEndpoint                         = "%s/blocks/%d/%d"
	lastBlockNumberEndpoint               = "%s/network/lastblocknumber"
	timeToNextEpochEndpoint               = "%s/network/timetonextepoch"
	approximatedEarningsEndpoint          = "%s/approximated_earnings/%f"
	pricesEndpoint                        = "%s/prices"
	activeMinersEndpoint                  = "%s/pool/activeminers"
	activeWorkersEndpoint                 = "%s/pool/activeworkers"
	poolHashrateEndpoint                  = "%s/pool/hashrate"
	topMinersEndpoint                     = "%s/pool/topminers"
)
