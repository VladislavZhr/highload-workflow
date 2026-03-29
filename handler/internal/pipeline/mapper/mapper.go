package mapper

import (
	"strconv"
)

const (
	dataSourceCode    = "NBU"
	serviceCode       = "NBCRPPLoanList"
	requestResultCode = "1"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Map(input BusinessInput, requestDate string) BusinessOutput {
	output := BusinessOutput{
		NBCR: make([]NBCR, 0, len(input.Counterparties.Counterparty)),
		ESBResultStatuses: ESBResultStatuses{
			RequestResultStatus: make([]RequestResultStatus, 0, len(input.Counterparties.Counterparty)),
		},
	}

	for _, counterparty := range input.Counterparties.Counterparty {
		output.NBCR = append(output.NBCR, m.mapNBCR(counterparty, requestDate))
		output.ESBResultStatuses.RequestResultStatus = append(
			output.ESBResultStatuses.RequestResultStatus,
			m.mapRequestResultStatus(counterparty.CounterpartyID),
		)
	}

	return output
}

func (m *Mapper) mapNBCR(counterparty CounterpartyInput, requestDate string) NBCR {
	return NBCR{
		CounterpartyID: counterparty.CounterpartyID,
		Individual: Individual{
			InquiryList: InquiryList{
				DayNumber:     0,
				WeekNumber:    0,
				MonthNumber:   0,
				QuarterNumber: 0,
				YearNumber:    0,
				Inquiry: Inquiry{
					RequestDate: requestDate,
				},
			},
			Receipts:      m.mapReceipts(counterparty.Data.ResultKVI),
			InfoCreditors: m.mapInfoCreditors(counterparty.Data.InfoCreditor),
		},
	}
}

func (m *Mapper) mapReceipts(items []ResultKVIInput) Receipts {
	receipts := Receipts{
		Receipt: make([]Receipt, 0, len(items)),
	}

	for _, item := range items {
		receipts.Receipt = append(receipts.Receipt, Receipt{
			ReceiptNumber: strconv.Itoa(item.OrderBank),
			MarkOfOwner:   strconv.Itoa(item.MarkOfOwner),
			LastName:      item.FIO.LastName,
			FirstName:     item.FIO.FirstName,
			MiddleName:    item.FIO.MiddleName,
			BirthDate:     item.BirthDay,
			TaxNumber:     item.INN,
			FlagK060:      item.FlagK060,
			LoanContracts: m.mapLoanContracts(item.Credits),
			PledgeList:    m.mapPledgeList(item.Pledge),
		})
	}

	return receipts
}

func (m *Mapper) mapLoanContracts(items []ResultKVICredit) LoanContracts {
	contracts := LoanContracts{
		LoanContract: make([]LoanContract, 0, len(items)),
	}

	for _, item := range items {
		contracts.LoanContract = append(contracts.LoanContract, LoanContract{
			LoanType:             strconv.Itoa(item.TypeCredit),
			ContractReference:    item.NumberDog,
			BookingDate:          item.DogDay,
			MaturityDate:         item.EndDay,
			LimitAmount:          formatScaledAmount(item.SumZagal),
			CurrencyCode:         item.R030,
			OutstandingAmount:    formatOptionalScaledAmount(item.SumArrears),
			PrincipalAmount:      formatScaledAmount(item.SumArrearsBase),
			InterestAmount:       formatScaledAmount(item.SumArrearsProc),
			PrincipalOverdue:     formatScaledAmount(item.ArrearBase),
			InterestOverdue:      formatScaledAmount(item.ArrearProc),
			OverdueDaysPrincipal: strconv.Itoa(item.DayBase),
			OverdueDaysInterest:  strconv.Itoa(item.DayProc),
			Class:                item.Klass,
			TrancheList:          m.mapTrancheList(item.Tranche),
		})
	}

	return contracts
}

func (m *Mapper) mapTrancheList(items []ResultKVITranche) TrancheList {
	tranches := TrancheList{
		Tranche: make([]Tranche, 0, len(items)),
	}

	for _, item := range items {
		tranches.Tranche = append(tranches.Tranche, Tranche{
			TrancheReference:     item.NumDogTr,
			BookingDate:          item.DogDayTr,
			MaturityDate:         item.EndDayTr,
			TrancheAmount:        formatScaledAmount(item.SumZagalTr),
			CurrencyCode:         item.R030Tr,
			PrincipalAmount:      formatScaledAmount(item.SumArrearsTrBase),
			InterestAmount:       formatScaledAmount(item.SumArrearsTrProc),
			PrincipalOverdue:     formatScaledAmount(item.ArrearBaseTr),
			InterestOverdue:      formatScaledAmount(item.ArrearProcTr),
			OverdueDaysPrincipal: strconv.Itoa(item.DayBaseTr),
			OverdueDaysInterest:  strconv.Itoa(item.DayProcTr),
			Class:                item.KlassTr,
		})
	}

	return tranches
}

func (m *Mapper) mapPledgeList(items []PledgeInput) PledgeList {
	pledges := PledgeList{
		Pledge: make([]Pledge, 0, len(items)),
	}

	for _, item := range items {
		pledges.Pledge = append(pledges.Pledge, Pledge{
			PledgeDate:   item.PledgeDay,
			PledgeType:   item.S031,
			PledgeNumber: strconv.Itoa(item.OrderZastava),
		})
	}

	return pledges
}

func (m *Mapper) mapInfoCreditors(items []InfoCreditorInput) InfoCreditors {
	infoCreditors := InfoCreditors{
		InfoCreditor: make([]InfoCreditor, 0, len(items)),
	}

	for _, item := range items {
		infoCreditors.InfoCreditor = append(infoCreditors.InfoCreditor, InfoCreditor{
			Creditor:               item.Creditor,
			CreditorInfo:           item.CreditorInfo,
			RelationshipPersonK060: item.K060RespRelation,
			DebtorsComment:         item.Comment,
			IdentPerson:            m.mapIdentPerson(item.IndPerson),
			Credits:                m.mapCredits(item.Credit),
		})
	}

	return infoCreditors
}

func (m *Mapper) mapIdentPerson(item IndPersonInput) IdentPerson {
	return IdentPerson{
		LastName:       item.LastName,
		FirstName:      item.FirstName,
		MiddleName:     item.Patronymic,
		BirthDate:      derefString(item.BirthDate),
		TaxNumber:      item.IndPersonCode,
		Passport:       item.Passport,
		DocumentNumber: item.DocumentNumber,
	}
}

func (m *Mapper) mapCredits(items []InfoCreditInput) Credits {
	credits := Credits{
		Credit: make([]Credit, 0, len(items)),
	}

	for _, item := range items {
		credits.Credit = append(credits.Credit, Credit{
			Liability:   m.mapLiability(item.Liability),
			CreditLoans: m.mapCreditLoans(item.Loan),
		})
	}

	return credits
}

func (m *Mapper) mapLiability(item LiabilityInput) Liability {
	return Liability{
		LoanType:           item.F037LoanType,
		AgreementNumber:    item.AgreemNo,
		AgreementStartDate: item.AgreemStartDate,
		AgreementEndDate:   derefString(item.ObligationEndDate),
		CurrencyCode:       item.R030Currency,
		TotalAmount:        formatAmount(item.TotalAmount),
		BalanceAmount:      formatAmount(item.BalanceOff),
		OverdueAmount:      formatAmount(item.OverdueAmount),
		OverdueDays:        derefInt(item.DaysOverdue),
		CRTypeAssessment:   item.S083RiskTypeAssessment,
		Class:              item.Class,
		ClassAdjusted:      item.KorrClass,
		Factors:            m.mapFactors(item.Factors),
		CreditTrancheList:  m.mapLiabilityCreditTrancheList(item.Tranche),
		Collaterals:        m.mapCollaterals(item.Collateral),
	}
}

func (m *Mapper) mapLiabilityCreditTrancheList(items []LiabilityTrancheInput) CreditTrancheList {
	list := CreditTrancheList{
		CreditTranche: make([]CreditTranche, 0, len(items)),
	}

	for _, item := range items {
		list.CreditTranche = append(list.CreditTranche, CreditTranche{
			LoanType:           item.F037LoanType,
			AgreementStartDate: item.AgreemStartDate,
			AgreementEndDate:   derefString(item.ObligationEndDate),
			CurrencyCode:       item.R030Currency,
			TotalAmount:        formatAmount(item.TotalAmount),
			BalanceAmount:      formatAmount(item.BalanceOff),
			OverdueAmount:      formatAmount(item.OverdueAmount),
			OverdueDays:        derefInt(item.DaysOverdue),
			CRTypeAssessment:   item.S083RiskTypeAssessment,
			Class:              item.Class,
			ClassAdjusted:      item.KorrClass,
			Factors:            m.mapFactors(item.Factors),
		})
	}

	return list
}

func (m *Mapper) mapFactors(items []FactorInput) Factors {
	factors := Factors{
		Factor: make([]Factor, 0, len(items)),
	}

	for _, item := range items {
		factors.Factor = append(factors.Factor, Factor{
			CRDefaultEvent:     item.F075GRiskEventList,
			CRDefaultEventCode: item.RiskEventCode,
		})
	}

	return factors
}

func (m *Mapper) mapCollaterals(items []CollateralInput) Collaterals {
	collaterals := Collaterals{
		Collateral: make([]Collateral, 0, len(items)),
	}

	for _, item := range items {
		collaterals.Collateral = append(collaterals.Collateral, Collateral{
			AgreementStartDate: item.AgreemStartDate,
			Movables:           m.mapMovables(item.Movable),
			Immovables:         m.mapImmovables(item.Immovable),
			CollatDeposits:     m.mapDeposits(item.Deposit),
		})
	}

	return collaterals
}

func (m *Mapper) mapMovables(items []CollateralMovableInput) Movables {
	movables := Movables{
		CollatList: make([]CollateralListItem, 0, len(items)),
	}

	for _, item := range items {
		movables.CollatList = append(movables.CollatList, CollateralListItem{
			CollateralType:   item.S031ColType,
			CurrencyCode:     item.R030Currency,
			CollateralAmount: item.CollateralAmount,
		})
	}

	return movables
}

func (m *Mapper) mapImmovables(items []CollateralMovableInput) Immovables {
	immovables := Immovables{
		CollatList: make([]CollateralListItem, 0, len(items)),
	}

	for _, item := range items {
		immovables.CollatList = append(immovables.CollatList, CollateralListItem{
			CollateralType:   item.S031ColType,
			CurrencyCode:     item.R030Currency,
			CollateralAmount: item.CollateralAmount,
		})
	}

	return immovables
}

func (m *Mapper) mapDeposits(items []CollateralDepositInput) CollatDeposits {
	deposits := CollatDeposits{
		CollatList: make([]DepositCollatListItem, 0, len(items)),
	}

	for _, item := range items {
		deposits.CollatList = append(deposits.CollatList, DepositCollatListItem{
			CurrencyCode:     item.R030Currency,
			CollateralAmount: item.CollateralAmount,
		})
	}

	return deposits
}

func (m *Mapper) mapCreditLoans(items []LoanInput) CreditLoans {
	loans := CreditLoans{
		CreditLoan: make([]CreditLoan, 0, len(items)),
	}

	for _, item := range items {
		loans.CreditLoan = append(loans.CreditLoan, CreditLoan{
			LoanType:              item.F037LoanType,
			AgreementNumber:       item.AgreemNo,
			AgreementStartDate:    derefString(item.AgreemStartDate),
			AgreementEndDate:      derefString(item.AgreemEndDate),
			CurrencyCode:          item.R030Currency,
			TotalAmount:           formatAmount(item.TotalAmount),
			BalanceAmount:         formatAmount(item.BalanceOn),
			OverdueAmount:         formatAmount(item.OverdueAmount),
			OverdueDays:           strconv.Itoa(item.DaysOverdue),
			DebtReservesWritten:   formatAmount(item.DebtWrOn),
			IncomeReservesWritten: formatAmount(item.IncWrOn),
			OffBalanceWritten:     derefString(item.WriteOffDate),
			CRTypeAssessment:      item.S083RiskTypeAssessment,
			Class:                 item.Class,
			ClassAdjusted:         item.KorrClass,
			StateLoanProgram:      item.D180StateLoanProg,
			Factors:               m.mapFactors(item.Factors),
			CreditTrancheList:     m.mapLoanCreditTrancheList(item.Tranche),
			Collaterals:           m.mapCollaterals(item.Collateral),
		})
	}

	return loans
}

func (m *Mapper) mapLoanCreditTrancheList(items []LoanTrancheInput) CreditTrancheList {
	list := CreditTrancheList{
		CreditTranche: make([]CreditTranche, 0, len(items)),
	}

	for _, item := range items {
		list.CreditTranche = append(list.CreditTranche, CreditTranche{
			LoanType:           item.F037LoanType,
			AgreementStartDate: item.AgreemStartDate,
			AgreementEndDate:   derefString(item.AgreemEndDate),
			CurrencyCode:       item.R030Currency,
			TotalAmount:        formatAmount(item.TotalAmount),
			BalanceAmount:      formatAmount(item.BalanceOn),
			OverdueAmount:      formatAmount(item.OverdueAmount),
			OverdueDays:        strconv.Itoa(item.DaysOverdue),
			CRTypeAssessment:   item.S083RiskTypeAssessment,
			Class:              item.Class,
			ClassAdjusted:      item.KorrClass,
			Factors:            m.mapFactors(item.Factors),
		})
	}

	return list
}

func (m *Mapper) mapRequestResultStatus(counterpartyID string) RequestResultStatus {
	return RequestResultStatus{
		CounterpartyID:    counterpartyID,
		DataSourceCode:    dataSourceCode,
		ServiceCode:       serviceCode,
		RequestResultCode: requestResultCode,
	}
}

func formatScaledAmount(value float64) string {
	return strconv.FormatFloat(value/100, 'f', -1, 64)
}

func formatOptionalScaledAmount(value *float64) string {
	if value == nil {
		return ""
	}

	return formatScaledAmount(*value)
}

func formatAmount(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func derefInt(value *int) string {
	if value == nil {
		return ""
	}

	return strconv.Itoa(*value)
}
