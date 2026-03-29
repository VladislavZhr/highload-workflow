package mapper

import "encoding/xml"

type BusinessOutput struct {
	XMLName           xml.Name          `xml:"Response"`
	NBCR              []NBCR            `xml:"NBCR"`
	ESBResultStatuses ESBResultStatuses `xml:"ESBResultStatuses"`
}

type NBCR struct {
	CounterpartyID string     `xml:"CounterpartyID,attr"`
	Individual     Individual `xml:"Individual"`
}

type Individual struct {
	InquiryList   InquiryList   `xml:"InquiryList"`
	Receipts      Receipts      `xml:"Receipts"`
	InfoCreditors InfoCreditors `xml:"InfoCreditors"`
}

type InquiryList struct {
	DayNumber     int     `xml:"DayNumber,attr"`
	WeekNumber    int     `xml:"WeekNumber,attr"`
	MonthNumber   int     `xml:"MonthNumber,attr"`
	QuarterNumber int     `xml:"QuarterNumber,attr"`
	YearNumber    int     `xml:"YearNumber,attr"`
	Inquiry       Inquiry `xml:"Inquiry"`
}

type Inquiry struct {
	RequestDate string `xml:"RequestDate,attr"`
}

type Receipts struct {
	Receipt []Receipt `xml:"Receipt"`
}

type Receipt struct {
	ReceiptNumber string        `xml:"ReceiptNumber,attr"`
	MarkOfOwner   string        `xml:"MarkOfOwner,attr"`
	LastName      string        `xml:"LastName,attr"`
	FirstName     string        `xml:"FirstName,attr"`
	MiddleName    string        `xml:"MiddleName,attr"`
	BirthDate     string        `xml:"BirthDate,attr"`
	TaxNumber     string        `xml:"TaxNumber,attr"`
	FlagK060      bool          `xml:"FlagK060,attr"`
	LoanContracts LoanContracts `xml:"LoanContracts"`
	PledgeList    PledgeList    `xml:"PledgeList"`
}

type LoanContracts struct {
	LoanContract []LoanContract `xml:"LoanContract"`
}

type LoanContract struct {
	LoanType             string      `xml:"LoanType,attr"`
	ContractReference    string      `xml:"ContractReference,attr"`
	BookingDate          string      `xml:"BookingDate,attr"`
	MaturityDate         string      `xml:"MaturityDate,attr"`
	LimitAmount          string      `xml:"LimitAmount,attr"`
	CurrencyCode         string      `xml:"CurrencyCode,attr"`
	OutstandingAmount    string      `xml:"OutstandingAmount,attr,omitempty"`
	PrincipalAmount      string      `xml:"PrincipalAmount,attr"`
	InterestAmount       string      `xml:"InterestAmount,attr"`
	PrincipalOverdue     string      `xml:"PrincipalOverdue,attr"`
	InterestOverdue      string      `xml:"InterestOverdue,attr"`
	OverdueDaysPrincipal string      `xml:"OverdueDaysPrincipal,attr"`
	OverdueDaysInterest  string      `xml:"OverdueDaysInterest,attr"`
	Class                string      `xml:"Class,attr"`
	TrancheList          TrancheList `xml:"TrancheList"`
}

type TrancheList struct {
	Tranche []Tranche `xml:"Tranche"`
}

type Tranche struct {
	TrancheReference     string `xml:"TrancheReference,attr"`
	BookingDate          string `xml:"BookingDate,attr"`
	MaturityDate         string `xml:"MaturityDate,attr"`
	TrancheAmount        string `xml:"TrancheAmount,attr"`
	CurrencyCode         string `xml:"CurrencyCode,attr"`
	PrincipalAmount      string `xml:"PrincipalAmount,attr"`
	InterestAmount       string `xml:"InterestAmount,attr"`
	PrincipalOverdue     string `xml:"PrincipalOverdue,attr"`
	InterestOverdue      string `xml:"InterestOverdue,attr"`
	OverdueDaysPrincipal string `xml:"OverdueDaysPrincipal,attr"`
	OverdueDaysInterest  string `xml:"OverdueDaysInterest,attr"`
	Class                string `xml:"Class,attr"`
}

type PledgeList struct {
	Pledge []Pledge `xml:"Pledge"`
}

type Pledge struct {
	PledgeDate   string `xml:"PledgeDate,attr"`
	PledgeType   string `xml:"PledgeType,attr"`
	PledgeNumber string `xml:"PledgeNumber,attr"`
}

type InfoCreditors struct {
	InfoCreditor []InfoCreditor `xml:"InfoCreditor"`
}

type InfoCreditor struct {
	Creditor               string      `xml:"Creditor,attr"`
	CreditorInfo           string      `xml:"CreditorInfo,attr"`
	RelationshipPersonK060 string      `xml:"RelationshipPersonK060,attr"`
	DebtorsComment         string      `xml:"DebtorsComment,attr"`
	IdentPerson            IdentPerson `xml:"IdentPerson"`
	Credits                Credits     `xml:"Credits"`
}

type IdentPerson struct {
	LastName       string `xml:"LastName,attr"`
	FirstName      string `xml:"FirstName,attr"`
	MiddleName     string `xml:"MiddleName,attr"`
	BirthDate      string `xml:"BirthDate,attr,omitempty"`
	TaxNumber      string `xml:"TaxNumber,attr"`
	Passport       string `xml:"Passport,attr"`
	DocumentNumber string `xml:"DocumentNumber,attr"`
}

type Credits struct {
	Credit []Credit `xml:"Credit"`
}

type Credit struct {
	Liability   Liability   `xml:"Liability"`
	CreditLoans CreditLoans `xml:"CreditLoans"`
}

type Liability struct {
	LoanType           string            `xml:"LoanType,attr"`
	AgreementNumber    string            `xml:"AgreementNumber,attr"`
	AgreementStartDate string            `xml:"AgreementStartDate,attr"`
	AgreementEndDate   string            `xml:"AgreementEndDate,attr,omitempty"`
	CurrencyCode       string            `xml:"CurrencyCode,attr"`
	TotalAmount        string            `xml:"TotalAmount,attr"`
	BalanceAmount      string            `xml:"BalanceAmount,attr"`
	OverdueAmount      string            `xml:"OverdueAmount,attr"`
	OverdueDays        string            `xml:"OverdueDays,attr,omitempty"`
	CRTypeAssessment   string            `xml:"CRTypeAssessment,attr"`
	Class              string            `xml:"Class,attr"`
	ClassAdjusted      string            `xml:"ClassAdjusted,attr"`
	Factors            Factors           `xml:"Factors"`
	CreditTrancheList  CreditTrancheList `xml:"CreditTrancheList"`
	Collaterals        Collaterals       `xml:"Collaterals"`
}

type CreditTrancheList struct {
	CreditTranche []CreditTranche `xml:"CreditTranche"`
}

type CreditTranche struct {
	LoanType           string  `xml:"LoanType,attr"`
	AgreementStartDate string  `xml:"AgreementStartDate,attr"`
	AgreementEndDate   string  `xml:"AgreementEndDate,attr,omitempty"`
	CurrencyCode       string  `xml:"CurrencyCode,attr"`
	TotalAmount        string  `xml:"TotalAmount,attr"`
	BalanceAmount      string  `xml:"BalanceAmount,attr"`
	OverdueAmount      string  `xml:"OverdueAmount,attr"`
	OverdueDays        string  `xml:"OverdueDays,attr,omitempty"`
	CRTypeAssessment   string  `xml:"CRTypeAssessment,attr"`
	Class              string  `xml:"Class,attr"`
	ClassAdjusted      string  `xml:"ClassAdjusted,attr"`
	Factors            Factors `xml:"Factors"`
}

type Factors struct {
	Factor []Factor `xml:"Factor"`
}

type Factor struct {
	CRDefaultEvent     string `xml:"CRDefaultEvent,attr"`
	CRDefaultEventCode string `xml:"CRDefaultEventCode,attr"`
}

type Collaterals struct {
	Collateral []Collateral `xml:"Collateral"`
}

type Collateral struct {
	AgreementStartDate string         `xml:"AgreementStartDate,attr"`
	Movables           Movables       `xml:"Movables"`
	Immovables         Immovables     `xml:"Immovables"`
	CollatDeposits     CollatDeposits `xml:"CollatDeposits"`
}

type Movables struct {
	CollatList []CollateralListItem `xml:"CollatList"`
}

type Immovables struct {
	CollatList []CollateralListItem `xml:"CollatList"`
}

type CollatDeposits struct {
	CollatList []DepositCollatListItem `xml:"CollatList"`
}

type CollateralListItem struct {
	CollateralType   string `xml:"CollateralType,attr"`
	CurrencyCode     string `xml:"CurrencyCode,attr"`
	CollateralAmount string `xml:"CollateralAmount,attr"`
}

type DepositCollatListItem struct {
	CurrencyCode     string `xml:"CurrencyCode,attr"`
	CollateralAmount string `xml:"CollateralAmount,attr"`
}

type CreditLoans struct {
	CreditLoan []CreditLoan `xml:"CreditLoan"`
}

type CreditLoan struct {
	LoanType              string            `xml:"LoanType,attr"`
	AgreementNumber       string            `xml:"AgreementNumber,attr"`
	AgreementStartDate    string            `xml:"AgreementStartDate,attr,omitempty"`
	AgreementEndDate      string            `xml:"AgreementEndDate,attr,omitempty"`
	CurrencyCode          string            `xml:"CurrencyCode,attr"`
	TotalAmount           string            `xml:"TotalAmount,attr"`
	BalanceAmount         string            `xml:"BalanceAmount,attr"`
	OverdueAmount         string            `xml:"OverdueAmount,attr"`
	OverdueDays           string            `xml:"OverdueDays,attr"`
	DebtReservesWritten   string            `xml:"DebtReservesWritten,attr"`
	IncomeReservesWritten string            `xml:"IncomeReservesWritten,attr"`
	OffBalanceWritten     string            `xml:"OffBalanceWritten,attr,omitempty"`
	CRTypeAssessment      string            `xml:"CRTypeAssessment,attr"`
	Class                 string            `xml:"Class,attr"`
	ClassAdjusted         string            `xml:"ClassAdjusted,attr"`
	StateLoanProgram      string            `xml:"StateLoanProgram,attr,omitempty"`
	Factors               Factors           `xml:"Factors"`
	CreditTrancheList     CreditTrancheList `xml:"CreditTrancheList"`
	Collaterals           Collaterals       `xml:"Collaterals"`
}

type ESBResultStatuses struct {
	RequestResultStatus []RequestResultStatus `xml:"RequestResultStatus"`
}

type RequestResultStatus struct {
	CounterpartyID    string `xml:"CounterpartyID,attr"`
	DataSourceCode    string `xml:"DataSourceCode,attr"`
	ServiceCode       string `xml:"ServiceCode,attr"`
	RequestResultCode string `xml:"RequestResultCode,attr"`
	ErrorMessage      string `xml:"ErrorMessage,attr,omitempty"`
	ErrorDescription  string `xml:"ErrorDescription,attr,omitempty"`
}
