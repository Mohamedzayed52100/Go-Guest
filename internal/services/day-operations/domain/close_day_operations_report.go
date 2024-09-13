package domain

type CloseDayOperationsReport struct {
	BranchName                 string
	Date                       string
	TotalLeftReservations      int64
	TotalLeftGuests            int64
	TotalNoShowReservations    int64
	TotalNoShowGuests          int64
	TotalWalkInReservations    int64
	TotalWalkInGuests          int64
	TotalCancelledReservations int64
	TotalCancelledGuests       int64
	TotalSales                 int64
	AverageCheckPerReservation int64
	AverageCheckPerGuest       int64
}
