package remote

type ErrorResponse struct {
	Res         bool
	Code        string
	Description string
}

type SuccessResponse struct {
	Res         bool
	Code        string
	Description string
}
