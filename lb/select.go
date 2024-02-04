package lb

//go:generate mockgen -source=select.go -destination=./mock/select.go -package=mock SelectService
type SelectService interface {
	ServerURI() (string, error)
}
