package instr

type Instrument interface {
	comparable
	BaseAsset() Asset
	TermAsset() Asset
	String() string
	NegativePriceAllowed() bool
}
