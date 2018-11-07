package test

var name *string

func SetName() {
	var n string
	n = "wang"
	name = &n
}

func GetName() *string {
	return name
}
