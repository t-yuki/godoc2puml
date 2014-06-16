package relations

type relations struct {
	struct1        // composition
	iface1         // implementation
	xassociation   *struct1
	yembed, zembed struct {
		*struct1 // association
		iface1   // association
		nested   struct {
			if1 []iface1
		}
	}
	//zzembed []struct { // TODO: what we should do???
	// 	iface1
	//	i64 int64
	//}
}

type struct1 struct {
}

type iface1 interface {
}
