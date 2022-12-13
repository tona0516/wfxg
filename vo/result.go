package vo

type Result[T any] struct {
	Result T
	Error  error
}
