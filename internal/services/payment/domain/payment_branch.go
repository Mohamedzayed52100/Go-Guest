package domain

type PaymentBranch struct {
	Name          string
	Address       string
	VatPercent    float32
	ServiceCharge float32
	CrNumber      string
	VatRegNumber  string
}
