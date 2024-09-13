package domain

type CoverFlow struct {
	Time         string
	Reservations []*CoverFlowReservation
}

type CoverFlowReservation struct {
	ID           int32
	GuestsNumber int32
	Status       CoverFlowReservationStatus
}

type CoverFlowReservationStatus struct {
	ID    int32
	Name  string
	Color string
	Icon  string
}
