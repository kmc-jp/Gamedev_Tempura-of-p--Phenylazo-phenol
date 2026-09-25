package serialize

type Serializable interface {
	// このオブジェクトが表すインスタンスを生成する
	Deserialize() (any, error)
}
