package dexalot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Network int
type contractType string
type Contract string

const (
	MAINNET = Network(iota + 1)
	TESTNET
)

const (
	PortfolioSub = Contract("PortfolioSub")
	TradePairs   = Contract("TradePairs")
)

const (
	portfolio  = contractType("Portfolio")
	tradePairs = contractType("TradePairs")
	orderBooks = contractType("OrderBooks")
	exchange   = contractType("Exchange")
)

type ContractData struct {
	Address common.Address
	ABI     abi.ABI
}

func convertContractToRequestString(contract Contract) string {
	switch contract {
	case PortfolioSub:
		return "Portfolio"
	case TradePairs:
		return "TradePairs"
	default:
		panic(fmt.Sprintf("unknown contract %s", contract))
	}
}

type payload struct {
	ABI    abi.ABI `json:"abi"`
	Format string  `json:"_format"`
}

type deployment struct {
	ParentEnvironment string  `json:"parentenv"`
	Environment       string  `json:"env"`
	EnvironmentType   string  `json:"env_type"`
	ContractName      string  `json:"contract_name"`
	ContractType      string  `json:"contract_type"`
	Address           string  `json:"address"`
	ImplAddress       string  `json:"impl_address"`
	Version           string  `json:"version"`
	Owner             string  `json:"owner"`
	Status            string  `json:"status"`
	Action            *string `json:"action"`
	Payload           payload `json:"abi"`
}

type RESTHelper struct {
	apiKey string
	env    Network
}

func (r *RESTHelper) GetContractData(contract Contract) (*ContractData, error) {
	var apiURL url.URL
	if r.env == MAINNET {
		apiURL = url.URL{
			Scheme: "https",
			Host:   "api.dexalot.com",
			Path:   "/privapi/trading/deployment",
			RawQuery: url.Values{
				"contracttype": []string{convertContractToRequestString(contract)},
				"env":          []string{"production-multi-subnet"},
				"returnabi":    []string{"true"},
			}.Encode(),
		}
	} else {
		apiURL = url.URL{
			Scheme: "https",
			Host:   "api.dexalot-test.com",
			Path:   "/privapi/trading/deployment",
			RawQuery: url.Values{
				"contracttype": []string{convertContractToRequestString(contract)},
				"env":          []string{"fuji-multi-subnet"},
				"returnabi":    []string{"true"},
			}.Encode(),
		}
	}

	req, err := http.NewRequest(http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-apikey", r.apiKey)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", response.StatusCode)
	}
	defer response.Body.Close()

	var tmp []deployment
	err = json.NewDecoder(response.Body).Decode(&tmp)
	if err != nil {
		return nil, err
	}

	for _, d := range tmp {
		if d.ContractName == string(contract) {
			return &ContractData{
				Address: common.HexToAddress(d.Address),
				ABI:     d.Payload.ABI,
			}, nil
		}
	}

	return nil, fmt.Errorf("contract %s not found", contract)
}

func NewRESTHelper(apiKey string, env Network) *RESTHelper {
	return &RESTHelper{
		apiKey: apiKey,
		env:    env,
	}
}
