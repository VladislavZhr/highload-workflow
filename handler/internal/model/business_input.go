package model

type BusinessInput struct {
	Counterparties CounterpartiesPayload `json:"counterparties"`
}

type CounterpartiesPayload struct {
	Counterparty []CounterpartyInput `json:"counterparty"`
}

type CounterpartyInput struct {
	CounterpartyID string           `json:"CounterpartyID"`
	Data           CounterpartyData `json:"data"`
}

type CounterpartyData struct {
	Status       string              `json:"status"`
	ResultKVI    []ResultKVIInput    `json:"result_kvi"`
	InfoCreditor []InfoCreditorInput `json:"info_creditor"`
}

type ResultKVIInput struct {
	OrderBank   int               `json:"orderBank"`
	MarkOfOwner int               `json:"markOfOwner"`
	FIO         ResultKVIFIO      `json:"fio"`
	BirthDay    string            `json:"birthDay"`
	INN         string            `json:"inn"`
	Credits     []ResultKVICredit `json:"credits"`
	FlagK060    bool              `json:"flagK060"`
	Pledge      []PledgeInput     `json:"pledge"`
}

type ResultKVIFIO struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

type ResultKVICredit struct {
	TypeCredit     int                `json:"typeCredit"`
	NumberDog      string             `json:"numberDog"`
	DogDay         string             `json:"dogDay"`
	EndDay         string             `json:"endDay"`
	SumZagal       float64            `json:"sumZagal"`
	R030           string             `json:"r030"`
	SumArrears     *float64           `json:"sumArrears"`
	SumArrearsBase float64            `json:"sumArrearsBase"`
	SumArrearsProc float64            `json:"sumArrearsProc"`
	ArrearBase     float64            `json:"arrearBase"`
	ArrearProc     float64            `json:"arrearProc"`
	DayBase        int                `json:"dayBase"`
	DayProc        int                `json:"dayProc"`
	Klass          string             `json:"klass"`
	Tranche        []ResultKVITranche `json:"tranche"`
}

type ResultKVITranche struct {
	NumDogTr         string   `json:"numDogTr"`
	DogDayTr         string   `json:"dogDayTr"`
	EndDayTr         string   `json:"endDayTr"`
	SumZagalTr       float64  `json:"sumZagalTr"`
	R030Tr           string   `json:"r030Tr"`
	SumArrearsTr     *float64 `json:"sumArrearsTr"`
	SumArrearsTrBase float64  `json:"sumArrearsTrBase"`
	SumArrearsTrProc float64  `json:"sumArrearsTrProc"`
	ArrearBaseTr     float64  `json:"arrearBaseTr"`
	ArrearProcTr     float64  `json:"arrearProcTr"`
	DayBaseTr        int      `json:"dayBaseTr"`
	DayProcTr        int      `json:"dayProcTr"`
	KlassTr          string   `json:"klassTr"`
}

type PledgeInput struct {
	PledgeDay    string `json:"pledgeDay"`
	S031         string `json:"s031"`
	OrderZastava int    `json:"orderZastava"`
}

type InfoCreditorInput struct {
	Creditor         string            `json:"creditor"`
	CreditorInfo     string            `json:"creditor_info"`
	IndPerson        IndPersonInput    `json:"ind_person"`
	Credit           []InfoCreditInput `json:"credit"`
	K060RespRelation string            `json:"k060_resp_relation"`
	Comment          string            `json:"comment"`
}

type IndPersonInput struct {
	LastName       string  `json:"last_name"`
	FirstName      string  `json:"first_name"`
	Patronymic     string  `json:"patronymic"`
	BirthDate      *string `json:"birth_date"`
	IndPersonCode  string  `json:"ind_person_code"`
	Passport       string  `json:"passport"`
	DocumentNumber string  `json:"document_number"`
}

type InfoCreditInput struct {
	Liability LiabilityInput `json:"liability"`
	Loan      []LoanInput    `json:"loan"`
}

type LiabilityInput struct {
	F037LoanType           string                  `json:"f037_loan_type"`
	AgreemNo               string                  `json:"agreem_no"`
	AgreemStartDate        string                  `json:"agreem_start_date"`
	ObligationEndDate      *string                 `json:"obligation_end_date"`
	R030Currency           string                  `json:"r030_currency"`
	TotalAmount            float64                 `json:"total_amount"`
	BalanceOff             float64                 `json:"balance_off"`
	OverdueAmount          float64                 `json:"overdue_amount"`
	DaysOverdue            *int                    `json:"days_overdue"`
	S083RiskTypeAssessment string                  `json:"s083_risk_type_assessment"`
	Class                  string                  `json:"class"`
	KorrClass              string                  `json:"korr_class"`
	Factors                []FactorInput           `json:"factors"`
	Tranche                []LiabilityTrancheInput `json:"tranche"`
	Collateral             []CollateralInput       `json:"collateral"`
}

type LiabilityTrancheInput struct {
	F037LoanType           string        `json:"f037_loan_type"`
	AgreemStartDate        string        `json:"agreem_start_date"`
	ObligationEndDate      *string       `json:"obligation_end_date"`
	R030Currency           string        `json:"r030_currency"`
	TotalAmount            float64       `json:"total_amount"`
	BalanceOff             float64       `json:"balance_off"`
	OverdueAmount          float64       `json:"overdue_amount"`
	DaysOverdue            *int          `json:"days_overdue"`
	S083RiskTypeAssessment string        `json:"s083_risk_type_assessment"`
	Class                  string        `json:"class"`
	KorrClass              string        `json:"korr_class"`
	Factors                []FactorInput `json:"factors"`
}

type LoanInput struct {
	F037LoanType           string             `json:"f037_loan_type"`
	AgreemNo               string             `json:"agreem_no"`
	AgreemStartDate        *string            `json:"agreem_start_date"`
	AgreemEndDate          *string            `json:"agreem_end_date"`
	R030Currency           string             `json:"r030_currency"`
	TotalAmount            float64            `json:"total_amount"`
	BalanceOn              float64            `json:"balance_on"`
	OverdueAmount          float64            `json:"overdue_amount"`
	DaysOverdue            int                `json:"days_overdue"`
	DebtWrOn               float64            `json:"debt_wr_on"`
	IncWrOn                float64            `json:"inc_wr_on"`
	WriteOffDate           *string            `json:"write_off_date"`
	S083RiskTypeAssessment string             `json:"s083_risk_type_assessment"`
	Class                  string             `json:"class"`
	KorrClass              string             `json:"korr_class"`
	D180StateLoanProg      string             `json:"d180_state_loan_prog"`
	Factors                []FactorInput      `json:"factors"`
	Tranche                []LoanTrancheInput `json:"tranche"`
	Collateral             []CollateralInput  `json:"collateral"`
}

type LoanTrancheInput struct {
	F037LoanType           string        `json:"f037_loan_type"`
	AgreemStartDate        string        `json:"agreem_start_date"`
	AgreemEndDate          *string       `json:"agreem_end_date"`
	R030Currency           string        `json:"r030_currency"`
	TotalAmount            float64       `json:"total_amount"`
	BalanceOn              float64       `json:"balance_on"`
	OverdueAmount          float64       `json:"overdue_amount"`
	DaysOverdue            int           `json:"days_overdue"`
	DebtWrOn               float64       `json:"debt_wr_on"`
	IncWrOn                float64       `json:"inc_wr_on"`
	WriteOffDate           *string       `json:"write_off_date"`
	S083RiskTypeAssessment string        `json:"s083_risk_type_assessment"`
	Class                  string        `json:"class"`
	KorrClass              string        `json:"korr_class"`
	Factors                []FactorInput `json:"factors"`
}

type FactorInput struct {
	F075GRiskEventList string `json:"f075g_risk_event_list"`
	RiskEventCode      string `json:"risk_event_code"`
}

type CollateralInput struct {
	AgreemStartDate string                   `json:"agreem_start_date"`
	Movable         []CollateralMovableInput `json:"movable"`
	Immovable       []CollateralMovableInput `json:"immovable"`
	Deposit         []CollateralDepositInput `json:"deposit"`
}

type CollateralMovableInput struct {
	S031ColType      string `json:"s031_col_type"`
	R030Currency     string `json:"r030_currency"`
	CollateralAmount string `json:"collateral_amount"`
}

type CollateralDepositInput struct {
	R030Currency     string `json:"r030_currency"`
	CollateralAmount string `json:"collateral_amount"`
}
