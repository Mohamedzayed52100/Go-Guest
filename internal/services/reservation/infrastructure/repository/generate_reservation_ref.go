package repository

import "math/rand"

func (r *ReservationRepository) GenerateReservationRef()string {
	var letters = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}