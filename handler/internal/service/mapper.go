package service

import (
	"strconv"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
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

func (m *Mapper) Map(input model.BusinessInput, requestDate string) model.BusinessOutput {
	output := model.BusinessOutput{
		NBCR: make([]model.NBCR, 0, len(input.Counterparties.Counterparty)),
		ESBResultStatuses: model.ESBResultStatuses{
			RequestResultStatus: make([]model.RequestResultStatus, 0, len(input.Counterparties.Counterparty)),
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

func (m *Mapper) mapNBCR(counterparty model.CounterpartyInput, requestDate string) model.NBCR {
	return model.NBCR{
		CounterpartyID: counterparty.CounterpartyID,
		Individual: model.Individual{
			InquiryList: model.InquiryList{
				DayNumber:     0,
				WeekNumber:    0,
				MonthNumber:   0,
				QuarterNumber: 0,
				YearNumber:    0,
				Inquiry: model.Inquiry{
					RequestDate: requestDate,
				},
			},
			Receipts:      m.mapReceipts(counterparty.Data.ResultKVI),
			InfoCreditors: m.mapInfoCreditors(counterparty.Data.InfoCreditor),
		},
	}
}

func (m *Mapper) mapReceipts(items []model.ResultKVIInput) model.Receipts {
	receipts := model.Receipts{
		Receipt: make([]model.Receipt, 0, len(items)),
	}

	for _, item := range items {
		receipts.Receipt = append(receipts.Receipt, model.Receipt{
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

func (m *Mapper) mapLoanContracts(items []model.ResultKVICredit) model.LoanContracts {
	contracts := model.LoanContracts{
		LoanContract: make([]model.LoanContract, 0, len(items)),
	}

	for _, item := range items {
		contracts.LoanContract = append(contracts.LoanContract, model.LoanContract{
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

func (m *Mapper) mapTrancheList(items []model.ResultKVITranche) model.TrancheList {
	tranches := model.TrancheList{
		Tranche: make([]model.Tranche, 0, len(items)),
	}

	for _, item := range items {
		tranches.Tranche = append(tranches.Tranche, model.Tranche{
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

func (m *Mapper) mapPledgeList(items []model.PledgeInput) model.PledgeList {
	pledges := model.PledgeList{
		Pledge: make([]model.Pledge, 0, len(items)),
	}

	for _, item := range items {
		pledges.Pledge = append(pledges.Pledge, model.Pledge{
			PledgeDate:   item.PledgeDay,
			PledgeType:   item.S031,
			PledgeNumber: strconv.Itoa(item.OrderZastava),
		})
	}

	return pledges
}

func (m *Mapper) mapInfoCreditors(items []model.InfoCreditorInput) model.InfoCreditors {
	infoCreditors := model.InfoCreditors{
		InfoCreditor: make([]model.InfoCreditor, 0, len(items)),
	}

	for _, item := range items {
		infoCreditors.InfoCreditor = append(infoCreditors.InfoCreditor, model.InfoCreditor{
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

func (m *Mapper) mapIdentPerson(item model.IndPersonInput) model.IdentPerson {
	return model.IdentPerson{
		LastName:       item.LastName,
		FirstName:      item.FirstName,
		MiddleName:     item.Patronymic,
		BirthDate:      derefString(item.BirthDate),
		TaxNumber:      item.IndPersonCode,
		Passport:       item.Passport,
		DocumentNumber: item.DocumentNumber,
	}
}

func (m *Mapper) mapCredits(items []model.InfoCreditInput) model.Credits {
	credits := model.Credits{
		Credit: make([]model.Credit, 0, len(items)),
	}

	for _, item := range items {
		credits.Credit = append(credits.Credit, model.Credit{
			Liability:   m.mapLiability(item.Liability),
			CreditLoans: m.mapCreditLoans(item.Loan),
		})
	}

	return credits
}

func (m *Mapper) mapLiability(item model.LiabilityInput) model.Liability {
	return model.Liability{
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

func (m *Mapper) mapLiabilityCreditTrancheList(items []model.LiabilityTrancheInput) model.CreditTrancheList {
	list := model.CreditTrancheList{
		CreditTranche: make([]model.CreditTranche, 0, len(items)),
	}

	for _, item := range items {
		list.CreditTranche = append(list.CreditTranche, model.CreditTranche{
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

func (m *Mapper) mapFactors(items []model.FactorInput) model.Factors {
	factors := model.Factors{
		Factor: make([]model.Factor, 0, len(items)),
	}

	for _, item := range items {
		factors.Factor = append(factors.Factor, model.Factor{
			CRDefaultEvent:     item.F075GRiskEventList,
			CRDefaultEventCode: item.RiskEventCode,
		})
	}

	return factors
}

func (m *Mapper) mapCollaterals(items []model.CollateralInput) model.Collaterals {
	collaterals := model.Collaterals{
		Collateral: make([]model.Collateral, 0, len(items)),
	}

	for _, item := range items {
		collaterals.Collateral = append(collaterals.Collateral, model.Collateral{
			AgreementStartDate: item.AgreemStartDate,
			Movables:           m.mapMovables(item.Movable),
			Immovables:         m.mapImmovables(item.Immovable),
			CollatDeposits:     m.mapDeposits(item.Deposit),
		})
	}

	return collaterals
}

func (m *Mapper) mapMovables(items []model.CollateralMovableInput) model.Movables {
	movables := model.Movables{
		CollatList: make([]model.CollateralListItem, 0, len(items)),
	}

	for _, item := range items {
		movables.CollatList = append(movables.CollatList, model.CollateralListItem{
			CollateralType:   item.S031ColType,
			CurrencyCode:     item.R030Currency,
			CollateralAmount: item.CollateralAmount,
		})
	}

	return movables
}

func (m *Mapper) mapImmovables(items []model.CollateralMovableInput) model.Immovables {
	immovables := model.Immovables{
		CollatList: make([]model.CollateralListItem, 0, len(items)),
	}

	for _, item := range items {
		immovables.CollatList = append(immovables.CollatList, model.CollateralListItem{
			CollateralType:   item.S031ColType,
			CurrencyCode:     item.R030Currency,
			CollateralAmount: item.CollateralAmount,
		})
	}

	return immovables
}

func (m *Mapper) mapDeposits(items []model.CollateralDepositInput) model.CollatDeposits {
	deposits := model.CollatDeposits{
		CollatList: make([]model.DepositCollatListItem, 0, len(items)),
	}

	for _, item := range items {
		deposits.CollatList = append(deposits.CollatList, model.DepositCollatListItem{
			CurrencyCode:     item.R030Currency,
			CollateralAmount: item.CollateralAmount,
		})
	}

	return deposits
}

func (m *Mapper) mapCreditLoans(items []model.LoanInput) model.CreditLoans {
	loans := model.CreditLoans{
		CreditLoan: make([]model.CreditLoan, 0, len(items)),
	}

	for _, item := range items {
		loans.CreditLoan = append(loans.CreditLoan, model.CreditLoan{
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

func (m *Mapper) mapLoanCreditTrancheList(items []model.LoanTrancheInput) model.CreditTrancheList {
	list := model.CreditTrancheList{
		CreditTranche: make([]model.CreditTranche, 0, len(items)),
	}

	for _, item := range items {
		list.CreditTranche = append(list.CreditTranche, model.CreditTranche{
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

func (m *Mapper) mapRequestResultStatus(counterpartyID string) model.RequestResultStatus {
	return model.RequestResultStatus{
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
