package swagger

// Tradeable Contracts, Indices, and History
type Instrument struct {
	Symbol                         NullString  `json:"symbol"`
	RootSymbol                     NullString  `json:"rootSymbol,omitempty"`
	State                          NullString  `json:"state,omitempty"`
	Typ                            NullString  `json:"typ,omitempty"`
	Listing                        NullTime    `json:"listing,omitempty"`
	Front                          NullTime    `json:"front,omitempty"`
	Expiry                         NullTime    `json:"expiry,omitempty"`
	Settle                         NullTime    `json:"settle,omitempty"`
	RelistNullInterval             NullTime    `json:"relistNullInterval,omitempty"`
	InverseLeg                     NullString  `json:"inverseLeg,omitempty"`
	SellLeg                        NullString  `json:"sellLeg,omitempty"`
	BuyLeg                         NullString  `json:"buyLeg,omitempty"`
	OptionStrikePcnt               NullFloat64 `json:"optionStrikePcnt,omitempty"`
	OptionStrikeRound              NullFloat64 `json:"optionStrikeRound,omitempty"`
	OptionStrikePrice              NullFloat64 `json:"optionStrikePrice,omitempty"`
	OptionMultiplier               NullFloat64 `json:"optionMultiplier,omitempty"`
	PositionCurrency               NullString  `json:"positionCurrency,omitempty"`
	Underlying                     NullString  `json:"underlying,omitempty"`
	QuoteCurrency                  NullString  `json:"quoteCurrency,omitempty"`
	UnderlyingSymbol               NullString  `json:"underlyingSymbol,omitempty"`
	Reference                      NullString  `json:"reference,omitempty"`
	ReferenceSymbol                NullString  `json:"referenceSymbol,omitempty"`
	CalcNullInterval               NullTime    `json:"calcNullInterval,omitempty"`
	PublishNullInterval            NullTime    `json:"publishNullInterval,omitempty"`
	PublishTime                    NullTime    `json:"publishTime,omitempty"`
	MaxOrderQty                    NullInt     `json:"maxOrderQty,omitempty"`
	MaxPrice                       NullFloat64 `json:"maxPrice,omitempty"`
	LotSize                        NullInt     `json:"lotSize,omitempty"`
	TickSize                       NullFloat64 `json:"tickSize,omitempty"`
	Multiplier                     NullInt     `json:"multiplier,omitempty"`
	SettlCurrency                  NullString  `json:"settlCurrency,omitempty"`
	UnderlyingToPositionMultiplier NullInt     `json:"underlyingToPositionMultiplier,omitempty"`
	UnderlyingToSettleMultiplier   NullInt     `json:"underlyingToSettleMultiplier,omitempty"`
	QuoteToSettleMultiplier        NullInt     `json:"quoteToSettleMultiplier,omitempty"`
	IsQuanto                       NullBool    `json:"isQuanto,omitempty"`
	IsInverse                      NullBool    `json:"isInverse,omitempty"`
	InitMargin                     NullFloat64 `json:"initMargin,omitempty"`
	MaNullIntMargin                NullFloat64 `json:"maNullIntMargin,omitempty"`
	RiskLimit                      NullInt     `json:"riskLimit,omitempty"`
	RiskStep                       NullInt     `json:"riskStep,omitempty"`
	Limit                          NullFloat64 `json:"limit,omitempty"`
	Capped                         NullBool    `json:"capped,omitempty"`
	Taxed                          NullBool    `json:"taxed,omitempty"`
	Deleverage                     NullBool    `json:"deleverage,omitempty"`
	MakerFee                       NullFloat64 `json:"makerFee,omitempty"`
	TakerFee                       NullFloat64 `json:"takerFee,omitempty"`
	SettlementFee                  NullFloat64 `json:"settlementFee,omitempty"`
	InsuranceFee                   NullFloat64 `json:"insuranceFee,omitempty"`
	FundingBaseSymbol              NullString  `json:"fundingBaseSymbol,omitempty"`
	FundingQuoteSymbol             NullString  `json:"fundingQuoteSymbol,omitempty"`
	FundingPremiumSymbol           NullString  `json:"fundingPremiumSymbol,omitempty"`
	FundingTimestamp               NullTime    `json:"fundingTimestamp,omitempty"`
	FundingNullInterval            NullTime    `json:"fundingNullInterval,omitempty"`
	FundingRate                    NullFloat64 `json:"fundingRate,omitempty"`
	IndicativeFundingRate          NullFloat64 `json:"indicativeFundingRate,omitempty"`
	RebalanceTimestamp             NullTime    `json:"rebalanceTimestamp,omitempty"`
	RebalanceNullInterval          NullTime    `json:"rebalanceNullInterval,omitempty"`
	OpeningTimestamp               NullTime    `json:"openingTimestamp,omitempty"`
	ClosingTimestamp               NullTime    `json:"closingTimestamp,omitempty"`
	SessionNullInterval            NullTime    `json:"sessionNullInterval,omitempty"`
	PrevClosePrice                 NullFloat64 `json:"prevClosePrice,omitempty"`
	LimitDownPrice                 NullFloat64 `json:"limitDownPrice,omitempty"`
	LimitUpPrice                   NullFloat64 `json:"limitUpPrice,omitempty"`
	BankruptLimitDownPrice         NullFloat64 `json:"bankruptLimitDownPrice,omitempty"`
	BankruptLimitUpPrice           NullFloat64 `json:"bankruptLimitUpPrice,omitempty"`
	PrevTotalVolume                NullInt     `json:"prevTotalVolume,omitempty"`
	TotalVolume                    NullInt     `json:"totalVolume,omitempty"`
	Volume                         NullInt     `json:"volume,omitempty"`
	Volume24h                      NullInt     `json:"volume24h,omitempty"`
	PrevTotalTurnover              NullInt     `json:"prevTotalTurnover,omitempty"`
	TotalTurnover                  NullInt     `json:"totalTurnover,omitempty"`
	Turnover                       NullInt     `json:"turnover,omitempty"`
	Turnover24h                    NullInt     `json:"turnover24h,omitempty"`
	PrevPrice24h                   NullFloat64 `json:"prevPrice24h,omitempty"`
	Vwap                           NullFloat64 `json:"vwap,omitempty"`
	HighPrice                      NullFloat64 `json:"highPrice,omitempty"`
	LowPrice                       NullFloat64 `json:"lowPrice,omitempty"`
	LastPrice                      NullFloat64 `json:"lastPrice,omitempty"`
	LastPriceProtected             NullFloat64 `json:"lastPriceProtected,omitempty"`
	LastTickDirection              NullString  `json:"lastTickDirection,omitempty"`
	LastChangePcnt                 NullFloat64 `json:"lastChangePcnt,omitempty"`
	BidPrice                       NullFloat64 `json:"bidPrice,omitempty"`
	MidPrice                       NullFloat64 `json:"midPrice,omitempty"`
	AskPrice                       NullFloat64 `json:"askPrice,omitempty"`
	ImpactBidPrice                 NullFloat64 `json:"impactBidPrice,omitempty"`
	ImpactMidPrice                 NullFloat64 `json:"impactMidPrice,omitempty"`
	ImpactAskPrice                 NullFloat64 `json:"impactAskPrice,omitempty"`
	HasLiquidity                   NullBool    `json:"hasLiquidity,omitempty"`
	OpenNullInterest               NullInt     `json:"openNullInterest,omitempty"`
	OpenValue                      NullInt     `json:"openValue,omitempty"`
	FairMethod                     NullString  `json:"fairMethod,omitempty"`
	FairBasisRate                  NullFloat64 `json:"fairBasisRate,omitempty"`
	FairBasis                      NullFloat64 `json:"fairBasis,omitempty"`
	FairPrice                      NullFloat64 `json:"fairPrice,omitempty"`
	MarkMethod                     NullString  `json:"markMethod,omitempty"`
	MarkPrice                      NullFloat64 `json:"markPrice,omitempty"`
	IndicativeTaxRate              NullFloat64 `json:"indicativeTaxRate,omitempty"`
	IndicativeSettlePrice          NullFloat64 `json:"indicativeSettlePrice,omitempty"`
	OptionUnderlyingPrice          NullFloat64 `json:"optionUnderlyingPrice,omitempty"`
	SettledPrice                   NullFloat64 `json:"settledPrice,omitempty"`
	Timestamp                      NullTime    `json:"timestamp,omitempty"`
}
