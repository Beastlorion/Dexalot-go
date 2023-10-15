package maker

type CumulativeQtySpreadModel interface {
	BidSpreadBps(totalQty float64) float64
	OfferSpreadBps(totalQty float64) float64
}

type SingleLayerSpreadModel struct {
	spreadBps float64
}

func (slsm *SingleLayerSpreadModel) BidSpreadBps(totalQty float64) float64 {
	return slsm.spreadBps
}

func (slsm *SingleLayerSpreadModel) OfferSpreadBps(totalQty float64) float64 {
	return slsm.spreadBps
}

func (slsm *SingleLayerSpreadModel) SetSpreadBps(spreadBps float64) {
	slsm.spreadBps = spreadBps
}

func NewSingleLayerSpreadModel() *SingleLayerSpreadModel {
	return &SingleLayerSpreadModel{}
}

type WithBiasSpreadModel[M CumulativeQtySpreadModel] struct {
	UnderlyingModel M
	bidBiasBps      float64
	offerBiasBps    float64
}

func (wbsm *WithBiasSpreadModel[M]) BidSpreadBps(totalQty float64) float64 {
	return wbsm.UnderlyingModel.BidSpreadBps(totalQty) + wbsm.bidBiasBps
}

func (wbsm *WithBiasSpreadModel[M]) OfferSpreadBps(totalQty float64) float64 {
	return wbsm.UnderlyingModel.OfferSpreadBps(totalQty) + wbsm.offerBiasBps
}

func (wbsm *WithBiasSpreadModel[M]) SetBidBiasBps(bidBiasBps float64) {
	wbsm.bidBiasBps = bidBiasBps
}

func (wbsm *WithBiasSpreadModel[M]) SetOfferBiasBps(offerBiasBps float64) {
	wbsm.offerBiasBps = offerBiasBps
}

func NewWithBiasSpreadModel[M CumulativeQtySpreadModel](underlyingModel M) *WithBiasSpreadModel[M] {
	return &WithBiasSpreadModel[M]{
		UnderlyingModel: underlyingModel,
	}
}

type TwoFactorSymmetricCurveModel struct {
	insideQty              float64
	insideSpreadBps        float64
	qtySpreadMultiplierBps float64
}

func (tfscm *TwoFactorSymmetricCurveModel) BidSpreadBps(totalQty float64) float64 {
	if totalQty < tfscm.insideQty {
		return tfscm.insideSpreadBps
	}
	return tfscm.insideSpreadBps + tfscm.qtySpreadMultiplierBps*(totalQty-tfscm.insideQty)
}

func (tfscm *TwoFactorSymmetricCurveModel) OfferSpreadBps(totalQty float64) float64 {
	if totalQty < tfscm.insideQty {
		return tfscm.insideSpreadBps
	}
	return tfscm.insideSpreadBps + tfscm.qtySpreadMultiplierBps*(totalQty-tfscm.insideQty)
}

func (tfscm *TwoFactorSymmetricCurveModel) SetInsideQty(insideQty float64) {
	tfscm.insideQty = insideQty
}

func (tfscm *TwoFactorSymmetricCurveModel) SetInsideSpreadBps(insideSpreadBps float64) {
	tfscm.insideSpreadBps = insideSpreadBps
}

func (tfscm *TwoFactorSymmetricCurveModel) SetQtySpreadMultiplierBps(qtySpreadMultiplierBps float64) {
	tfscm.qtySpreadMultiplierBps = qtySpreadMultiplierBps
}

func NewTwoFactorSymmetricCurveModel() *TwoFactorSymmetricCurveModel {
	return &TwoFactorSymmetricCurveModel{}
}

type StaticSymmetricCurveModel struct {
	modelSpreadBpsFn func(totalQty float64) float64
}

func (scm *StaticSymmetricCurveModel) BidSpreadBps(totalQty float64) float64 {
	return scm.modelSpreadBpsFn(totalQty)
}

func (scm *StaticSymmetricCurveModel) OfferSpreadBps(totalQty float64) float64 {
	return scm.modelSpreadBpsFn(totalQty)
}

func NewStaticSymmetricCurveModel(modelSpreadBpsFn func(totalQty float64) float64) *StaticSymmetricCurveModel {
	return &StaticSymmetricCurveModel{
		modelSpreadBpsFn: modelSpreadBpsFn,
	}
}
