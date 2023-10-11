package instr

import "strings"

type Asset string 

const (
    AVAX Asset = "AVAX"
    ALOT Asset = "ALOT"
	LOST Asset = "LOST"
	sAVAX Asset = "sAVAX"
	STL Asset = "STL"
    BTCb Asset = "BTC.b"
    WETHe Asset = "WETH.e"
    USDC Asset = "USDC"
    USDT Asset = "USDT"
    EUROC Asset = "EUROC"
)

func ToInternalAssetName(asset string) Asset {
    if strings.Contains(asset, ".") {
        prePost := strings.Split(asset, ".")
        return Asset(strings.ToUpper(prePost[0]) + "." + strings.ToLower(prePost[1]))
    }
	return Asset(strings.ToUpper(asset))
}
