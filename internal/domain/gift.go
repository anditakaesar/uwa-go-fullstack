package domain

type Gift struct {
	Base
	Title       string
	Description string
	Stock       int
	RedeemPoint int
	ImageURL    string
}
