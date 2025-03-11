package sms

import (
	"log"
)

func Send() {
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

	log.Printf("sms services has been called - not active at this time")

}
