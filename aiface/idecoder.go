package aiface

type IDecoder interface {
	Interceptor
	GetLengthField() *LengthField
}
