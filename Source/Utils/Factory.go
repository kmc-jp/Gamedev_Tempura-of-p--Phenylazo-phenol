package Utils

type Factory[T any] func() (T, error)

func CastFactory[FromType any, ToType any](factoryA Factory[FromType]) Factory[ToType] {
	return func() (ToType, error) {
		res, err := factoryA()
		if err != nil {
			var zero ToType
			return zero, err
		}
		// FromType が ToType を満たしている場合、ここで暗黙的に型変換されて返されます
		// (満たしていない場合はコンパイルエラーになります)
		return any(res).(ToType), nil
	}
}
