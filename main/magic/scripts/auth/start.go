package magic_auth

import (
	"fmt"
	"log"

	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/station/main/integration"
	"github.com/bytedance/sonic"
)

func GetStartForm(p *mconfig.Plan, print bool) {
	basePath := p.Environment["BASE_PATH"]

	// Make a request to the auth
	form, err := integration.PostRequest("http://"+basePath, "/v1/account/auth/form", map[string]interface{}{})
	if err != nil {
		log.Fatalln("couldn't get form:", err)
	}

	// Parse to JSON and print (if desired)
	if print {
		encoded, err := sonic.MarshalIndent(form, "", "  ")
		if err != nil {
			log.Fatalln("couldn't encode to json:", err)
		}
		fmt.Println(string(encoded))
	}

}
