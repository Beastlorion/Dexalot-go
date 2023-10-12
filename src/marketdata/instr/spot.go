package instr

import (
	"fmt"
	"strings"
)

type Spot struct {
	Base Asset
	Term Asset
}

func (s Spot) BaseAsset() Asset {
	return s.Base
}

func (s Spot) TermAsset() Asset {
	return s.Term
}

func (s Spot) String() string {
	return string(s.Base) + "-" + string(s.Term)
}

func (s Spot) NegativePriceAllowed() bool {
	return false
}

func ToSpotUnderscoreSeparated(pair string) (Spot, error) {
	assetPair := strings.Split(pair, "_")
	if len(assetPair) != 2 {
		return Spot{}, fmt.Errorf("invalid asset pair format: %s", pair)
	}
	return Spot{
		Base: ToInternalAssetName(assetPair[0]),
		Term: ToInternalAssetName(assetPair[1]),
	}, nil
}

func ToSpotSlashSeparated(pair string) (Spot, error) {
	assetPair := strings.Split(pair, "/")
	if len(assetPair) != 2 {
		return Spot{}, fmt.Errorf("invalid asset pair format: %s", pair)
	}
	return Spot{
		Base: ToInternalAssetName(assetPair[0]),
		Term: ToInternalAssetName(assetPair[1]),
	}, nil
}

func SpotToUpperCaseUnderscoreSeparated(pair Spot) string {
	return strings.ToUpper(string(pair.Base) + "_" + string(pair.Term))
}

func SpotToLowerCaseUnderscoreSeparated(pair Spot) string {
	return strings.ToLower(string(pair.Base) + "_" + string(pair.Term))
}

func SpotToUpperCaseDashSeparated(pair Spot) string {
	return strings.ToUpper(string(pair.Base) + "-" + string(pair.Term))
}

func SpotToLowerCaseDashSeparated(pair Spot) string {
	return strings.ToLower(string(pair.Base) + "-" + string(pair.Term))
}

func SpotToUpperCaseSlashSeparated(pair Spot) string {
	return strings.ToUpper(string(pair.Base) + "/" + string(pair.Term))
}

func SpotToLowerCaseSlashSeparated(pair Spot) string {
	return strings.ToLower(string(pair.Base) + "/" + string(pair.Term))
}
