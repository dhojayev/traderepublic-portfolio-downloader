package details

const (
	SectionTitleOverview                  = "Übersicht"
	SectionTitlePerformance               = "Performance"
	SectionTitleTransaction               = "Transaktion"
	sectionTypeHeader                     = "header"
	sectionTypeTable                      = "table"
	sectionTypeHorizontalTable            = "horizontalTable"
	sectionTypeDocuments                  = "documents"
	overviewDataTitleStatus               = "Status"
	overviewDataTitleEvent                = "Ereignis"
	overviewDataTitleOrderType            = "Orderart"
	overviewDataTitleOrderTypeAlt         = "Auftragsart"
	overviewDataTitleOrderTypeAlt2        = "Ordertyp"
	OverviewDataTitleAsset                = "Asset"
	overviewDataTitleProduct              = "Produkt"
	OverviewDataTitleUnderlyingAsset      = "Basiswert"
	overviewDataTitleReceivedFrom         = "Von"
	overviewDataTitleIBAN                 = "IBAN"
	overviewDataTitleDeposit              = "Zahlung"
	overviewDataTitleYoY                  = "Jahressatz"
	TransactionDataTitleShares            = "Anteile"
	TransactionDataTitleStocks            = "Aktien"
	TransactionDataTitleRate              = "Aktienkurs"
	TransactionDataTitlePrice             = "Anteilspreis"
	TransactionDataTitleDividendsPerStock = "Dividende je Aktie"
	TransactionDataTitleCommission        = "Gebühr"
	TransactionDataTitleTotal             = "Gesamt"
	PerformanceDataTitleYield             = "Rendite"
	PerformanceDataTitleProfit            = "Gewinn"
	PerformanceDataTitleLoss              = "Verlust"
	eventTypeTextPayout                   = "Ausschüttung"
	orderTypeTextsSale                    = "Verkauf"
	orderTypeTextsPurchase                = "Kauf"
	orderTypeTextsSavingsPlan             = "Sparplan"
	orderTypeTextsRoundUp                 = "Round up"
	orderTypeTextsSaveback                = "Saveback"
	TrendNegative                         = "negative"
)

// var (
// 	ErrSectionNotFound          = errors.New("section not found")
// 	ErrSectionDataEntryNotFound = errors.New("section data entry not found")
// )

// type Response struct {
// 	ID       string            `json:"id"`
// 	Sections []ResponseSection `json:"sections"`
// }
//
// func (r Response) HeaderSection() (ResponseSectionTypeHeader, error) {
// 	var section ResponseSectionTypeHeader
//
// 	err := r.UnmarshalSection("", sectionTypeHeader, &section)
//
// 	return section, err
// }
//
// func (r Response) OverviewSection() (ResponseSectionTypeTable, error) {
// 	var section ResponseSectionTypeTable
//
// 	err := r.UnmarshalSection(SectionTitleOverview, sectionTypeTable, &section)
//
// 	return section, err
// }
//
// func (r Response) PerformanceSection() (ResponseSectionTypeTable, error) {
// 	var section ResponseSectionTypeTable
//
// 	err := r.UnmarshalSection(SectionTitlePerformance, sectionTypeHorizontalTable, &section)
//
// 	return section, err
// }
//
// func (r Response) TransactionSection() (ResponseSectionTypeTable, error) {
// 	var section ResponseSectionTypeTable
//
// 	err := r.UnmarshalSection(SectionTitleTransaction, sectionTypeTable, &section)
//
// 	return section, err
// }
//
// func (r Response) DocumentsSection() (ResponseSectionTypeDocuments, error) {
// 	var section ResponseSectionTypeDocuments
//
// 	err := r.UnmarshalSection("", sectionTypeDocuments, &section)
//
// 	return section, err
// }
//
// func (r Response) UnmarshalSection(sectionTitle, sectionType string, v any) error {
// 	for _, s := range r.Sections {
// 		if (sectionTitle != "" && s.Title != sectionTitle) || s.Type != sectionType {
// 			continue
// 		}
//
// 		sectionBytes, err := json.Marshal(s)
// 		if err != nil {
// 			return fmt.Errorf("could not marshal %s section: %w", sectionTitle, err)
// 		}
//
// 		if err := json.Unmarshal(sectionBytes, v); err != nil {
// 			return fmt.Errorf("could not unmarshal %s section: %w", sectionTitle, err)
// 		}
//
// 		return nil
// 	}
//
// 	return ErrSectionNotFound
// }
//
// type ResponseSection struct {
// 	Action ResponseAction `json:"action,omitempty"`
// 	Data   any            `json:"data"`
// 	Title  string         `json:"title"`
// 	Type   string         `json:"type"`
// }
//
// type ResponseAction struct {
// 	Payload any    `json:"payload"`
// 	Type    string `json:"type"`
// }
//
// type ResponseSectionTypeHeader struct {
// 	ResponseSection
// 	Data ResponseSectionTypeHeaderData `json:"data"`
// }
//
// type ResponseSectionTypeHeaderData struct {
// 	Icon         string `json:"icon"`
// 	Status       string `json:"status"`
// 	SubtitleText string `json:"subtitleText,omitempty"`
// 	Timestamp    string `json:"timestamp"`
// }
//
// type ResponseSectionTypeTable struct {
// 	ResponseSection
// 	Data []ResponseSectionTypeTableData `json:"data"`
// }
//
// func (r ResponseSectionTypeTable) Status() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(overviewDataTitleStatus)
// }
//
// func (r ResponseSectionTypeTable) Event() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(overviewDataTitleEvent)
// }
//
// func (r ResponseSectionTypeTable) OrderType() (ResponseSectionTypeTableData, error) {
// 	data, err := r.findDataByTitle(overviewDataTitleOrderType)
// 	if err == nil {
// 		return data, nil
// 	}
//
// 	data, err = r.findDataByTitle(overviewDataTitleOrderTypeAlt)
// 	if err == nil {
// 		return data, nil
// 	}
//
// 	return r.findDataByTitle(overviewDataTitleOrderTypeAlt2)
// }
//
// func (r ResponseSectionTypeTable) Asset() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(OverviewDataTitleAsset)
// }
//
// func (r ResponseSectionTypeTable) Product() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(overviewDataTitleProduct)
// }
//
// func (r ResponseSectionTypeTable) UnderlyingAsset() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(OverviewDataTitleUnderlyingAsset)
// }
//
// func (r ResponseSectionTypeTable) ReceivedFrom() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(overviewDataTitleReceivedFrom)
// }
//
// func (r ResponseSectionTypeTable) IBAN() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(overviewDataTitleIBAN)
// }
//
// func (r ResponseSectionTypeTable) Deposit() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(overviewDataTitleDeposit)
// }
//
// func (r ResponseSectionTypeTable) YoY() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(overviewDataTitleYoY)
// }
//
// func (r ResponseSectionTypeTable) Shares() (ResponseSectionTypeTableData, error) {
// 	data, err := r.findDataByTitle(TransactionDataTitleShares)
// 	if err == nil {
// 		return data, nil
// 	}
//
// 	return r.findDataByTitle(TransactionDataTitleStocks)
// }
//
// func (r ResponseSectionTypeTable) Rate() (ResponseSectionTypeTableData, error) {
// 	data, err := r.findDataByTitle(TransactionDataTitleRate)
// 	if err == nil {
// 		return data, nil
// 	}
//
// 	data, err = r.findDataByTitle(TransactionDataTitlePrice)
// 	if err == nil {
// 		return data, nil
// 	}
//
// 	return r.findDataByTitle(TransactionDataTitleDividendsPerStock)
// }
//
// func (r ResponseSectionTypeTable) Commission() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(TransactionDataTitleCommission)
// }
//
// func (r ResponseSectionTypeTable) Total() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(TransactionDataTitleTotal)
// }
//
// func (r ResponseSectionTypeTable) Yield() (ResponseSectionTypeTableData, error) {
// 	return r.findDataByTitle(PerformanceDataTitleYield)
// }
//
// func (r ResponseSectionTypeTable) Profit() (ResponseSectionTypeTableData, error) {
// 	data, err := r.findDataByTitle(PerformanceDataTitleProfit)
// 	if err == nil {
// 		return data, nil
// 	}
//
// 	return r.findDataByTitle(PerformanceDataTitleLoss)
// }
//
// func (r ResponseSectionTypeTable) findDataByTitle(title string) (ResponseSectionTypeTableData, error) {
// 	for _, data := range r.Data {
// 		if data.Title != title {
// 			continue
// 		}
//
// 		return data, nil
// 	}
//
// 	return ResponseSectionTypeTableData{}, fmt.Errorf("%w (%s)", ErrSectionDataEntryNotFound, title)
// }
//
// type ResponseSectionTypeTableData struct {
// 	Detail ResponseSectionTypeTableDataDetail `json:"detail"`
// 	Style  string                             `json:"style"`
// 	Title  string                             `json:"title"`
// }
//
// func (r ResponseSectionTypeTableData) IsOrderTypeSale() bool {
// 	return strings.Contains(r.Detail.Text, orderTypeTextsSale)
// }
//
// func (r ResponseSectionTypeTableData) IsOrderTypePurchase() bool {
// 	return strings.Contains(r.Detail.Text, orderTypeTextsPurchase) || r.Detail.Text == orderTypeTextsSavingsPlan
// }
//
// func (r ResponseSectionTypeTableData) IsOrderTypeRoundUp() bool {
// 	return r.Detail.Text == orderTypeTextsRoundUp
// }
//
// func (r ResponseSectionTypeTableData) IsOrderTypeSaveback() bool {
// 	return r.Detail.Text == orderTypeTextsSaveback
// }
//
// func (r ResponseSectionTypeTableData) IsEventPayout() bool {
// 	return r.Detail.Text == eventTypeTextPayout
// }
//
// func (r ResponseSectionTypeTableData) HasSharesWithPeriod() bool {
// 	return r.Title == TransactionDataTitleStocks
// }
//
// type ResponseSectionTypeTableDataDetail struct {
// 	Action ResponseAction `json:"action,omitempty"`
// 	Text   string         `json:"text"`
// 	Trend  string         `json:"trend,omitempty"`
// 	Type   string         `json:"type"`
// }
//
// func (r ResponseSectionTypeTableDataDetail) IsTrendNegative() bool {
// 	return r.Trend == TrendNegative
// }
//
// type ResponseSectionTypeDocuments struct {
// 	ResponseSection
// 	Data []ResponseSectionTypeDocumentsData `json:"data"`
// }
//
// type ResponseSectionTypeDocumentsData struct {
// 	Action      ResponseAction `json:"action,omitempty"`
// 	Detail      string         `json:"detail"`
// 	ID          string         `json:"id"`
// 	PostboxType string         `json:"postboxType"`
// 	Title       string         `json:"title"`
// }
