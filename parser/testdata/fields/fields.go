package fields

type fields struct {
	uint8
	i8    int8
	embed struct {
		*int16
		i32    []int32
		nested struct {
			i64 int64
		}
	}
}
