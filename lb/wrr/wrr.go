package wrr

type WRR interface {
	Add(item any, weight uint32)
	Update(item any, weight uint32)
	Next() any
}
