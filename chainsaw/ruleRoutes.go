package chainsaw

var functionMap = map[string]func(){
	"78d5cab4-557e-454f-9fb9-a222bd0d5edc": Ident78d5cab4557e454f9fb9a222bd0d5edc,
}

var count int

func RunEventResponse(id string) {
	function, exists := functionMap[id]
	if exists {
		function()
	} else {
		//fmt.Println("Function does not exist.")
		return
	}
}

// RDP Connection
func Ident78d5cab4557e454f9fb9a222bd0d5edc() {

}
