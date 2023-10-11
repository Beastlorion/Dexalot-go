package instr

type Instrument interface {
	comparable
	BaseAsset() Asset
	TermAsset() Asset
	String() string
	NegativePriceAllowed() bool
}

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

