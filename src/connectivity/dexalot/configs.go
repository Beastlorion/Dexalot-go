package dexalot

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
)

type PairConfig struct {
	Env                  string `json:"env"`
	Pair                 string `json:"pair"`
	Base                 string `json:"base"`
	Quote                string `json:"quote"`
	BaseDisplayDecimals  int    `json:"basedisplaydecimals"`
	QuoteDisplayDecimals int    `json:"quotedisplaydecimals"`
	BaseEVMDecimals      int    `json:"base_evmdecimals"`
	QuoteEVMDecimals     int    `json:"quote_evmdecimals"`
	MinTradeAmount       string `json:"mintrade_amnt"`
	MaxTradeAmount       string `json:"maxtrade_amnt"`
	Status               string `json:"status"`
}

func GetPairConfigs(apiKey string, env Network) (map[instr.Spot]PairConfig, error) {
	var u string
	if env == MAINNET {
		u = "https://api.dexalot.com/privapi/trading/pairs"
	} else {
		u = "https://api.dexalot-test.com/privapi/trading/pairs"
	}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	req.Header.Add("x-apikey", apiKey)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		var temp []PairConfig
		if err := json.NewDecoder(res.Body).Decode(&temp); err != nil {
			return nil, err
		}
		config := make(map[instr.Spot]PairConfig)
		for _, pair := range temp {
			spot, err := ConvertToSpot(pair.Pair)
			if err != nil {
				return nil, err
			}
			config[spot] = pair
		}
		return config, err
	}
	return nil, err
}

// TODO - feel free to extract min quote time, post only, and replenishment rate as configs
func PairConfigToReferenceData(config *PairConfig) (*refdata.Composite, error) {
    minTradeSize, err := strconv.ParseFloat(config.MinTradeAmount, 64)
    if err != nil {
        return nil, err
    }
    maxTradeSize, err := strconv.ParseFloat(config.MaxTradeAmount, 64)
    if err != nil {
        return nil, err
    }

	priceQtyRefData := &refdata.Precision{
		PricePrecision: config.QuoteDisplayDecimals,
		QtyPrecision: config.BaseDisplayDecimals,
	}

	tradeSizeRefData := &refdata.TermTradeSizeLimit{
		TermMinTradeSize: minTradeSize,
		TermMaxTradeSize: maxTradeSize,
	}

    return &refdata.Composite{
		PriceQty: priceQtyRefData,
		TradeSizeLimit: tradeSizeRefData,
        MinQuoteTime: time.Duration(2 * time.Second),
        ReplenishmentRate: time.Duration(10 * time.Second),
        UsePostOnly: true,
    }, nil
}

type TokenConfig struct {
	Env         string `json:"env"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	IsNative    bool   `json:"isnative"`
	Address     string `json:"address"`
	EVMDecimals int    `json:"evmdecimals"`
	Status      string `json:"status"`
	AuctionMode int    `json:"auctionmode"`
}

func GetTokenConfig(apiKey string, env Network) (map[instr.Asset]TokenConfig, error) {
	var u string
	if env == MAINNET {
		u = "https://api.dexalot.com/privapi/trading/tokens"
	} else {
		u = "https://api.dexalot-test.com/privapi/trading/tokens"
	}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	req.Header.Add("x-apikey", apiKey)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		var temp []TokenConfig
		if err := json.NewDecoder(res.Body).Decode(&temp); err != nil {
			return nil, err
		}
		config := make(map[instr.Asset]TokenConfig)
		for _, token := range temp {
			sym := ConvertFromAPIAsset(token.Symbol)
			config[instr.Asset(sym)] = token
		}
		return config, err
	}
	return nil, err
}
