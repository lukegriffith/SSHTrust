package auth

type ACL struct {
	CA         string
	Principals []string
	Permission bool
}

var _ = []ACL{
	{
		CA:         "*",
		Principals: []string{},
		Permission: false,
	},
	{
		CA: "TestCA",
		Principals: []string{
			"group/test",
			"user/admin",
		},
		Permission: true,
	},
}
