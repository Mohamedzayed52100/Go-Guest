package domain

type GuestBranchVisit struct {
	Name   string `db:"name"`
	Visits int32  `db:"visits"`
}
