package relations

type relations struct {
	struct1     // composition
	iface1      // implementation
	association *struct1
	embed       struct {
		struct1 // association
		iface1  // association
		nested  struct {
			if1 []iface1
		}
	}
}

type struct1 struct {
}

type iface1 interface {
}
