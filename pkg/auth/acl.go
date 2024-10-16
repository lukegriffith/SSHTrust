package auth

type ACL struct {
	CA         string
	Principals []string
	Permission bool
}

var _ = []ACL{
	ACL{
		CA:         "*",
		Principals: []string{},
		Permission: false,
	},
	ACL{
		CA: "TestCA",
		Principals: []string{
			"group/test",
			"user/admin",
		},
		Permission: true,
	},
}
