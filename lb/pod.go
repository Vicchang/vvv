package lb

//go:generate mockgen -source=pod.go -destination=./mock/pod.go -package=mock PodService
type PodService interface {
	Add(string)
}
