package fields

type fields struct {
	uint8
	xi8             int8
	yembed, yyembed struct {
		*int16
		xi32, yi32 []int32
		znested    struct {
			i64 int64
		}
	}
	zembed *struct {
		*bool
	}
	//zzzembed []struct { // TODO: what we should do???
	//	i64 int64
	//}
}
