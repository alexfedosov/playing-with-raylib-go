package lib

type System interface {
	Update(w *World, deltaTime float64)
}
