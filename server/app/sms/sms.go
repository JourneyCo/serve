package sms

import (
	"log"
	"regexp"
	"strings"
)

// phoneRegex accepts only XXX-XXX-XXXX where 0<=X<=9
const phoneRegex = "^[0-9]{3}[-]?[0-9]{3}[-]?[0-9]{4}$"

// Send composes and sends a text message out to the given number
// via twillio
func Send(number string) {

	matched, err := regexp.MatchString(phoneRegex, number)
	if err != nil {
		log.Print("error matching regez to number: ", err)
		return
	}
	if !matched {
		log.Printf("number doesnt fit regex: %s", number)
		return
	}

	// append +1 and remove the hyphens from the phone number
	number = "+1" + strings.ReplaceAll(number, "+", "")

	// client := twilio.NewRestClientWithParams(
	// 	twilio.ClientParams{
	// 		Username: accountSid,
	// 		Password: authToken,
	// 	},
	// )
	//
	// params := &twilioApi.CreateMessageParams{}
	// params.SetTo("+15558675309")
	// params.SetFrom("+15017250604")
	// params.SetBody("Hello from Go!")
	//
	// resp, err := client.Api.CreateMessage(params)
	// if err != nil {
	// 	fmt.Println("Error sending SMS message: " + err.Error())
	// } else {
	// 	response, _ := json.Marshal(*resp)
	// 	fmt.Println("Response: " + string(response))
	// }

	log.Printf("sms services has been called - mock send to %s", number)

}
