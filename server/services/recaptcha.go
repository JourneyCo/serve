package services

import (
	"context"
	"fmt"

	recaptcha "cloud.google.com/go/recaptchaenterprise/v2/apiv1"
	recaptchapb "cloud.google.com/go/recaptchaenterprise/v2/apiv1/recaptchaenterprisepb"
	"serve/config"
)

/**
 * Create an assessment to analyze the risk of a UI action.
 *
 * @param projectID: Your Google Cloud Project ID.
 * @param recaptchaKey: The reCAPTCHA key associated with the site/app
 * @param token: The generated token obtained from the client.
 * @param recaptchaAction: Action name corresponding to the token.
 */
func CreateAssessment(cfg *config.Config, token string) error {

	// Create the reCAPTCHA client.
	ctx := context.Background()
	cli, err := recaptcha.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error creating reCAPTCHA client\n")
	}
	defer cli.Close()

	// Set the properties of the event to be tracked.
	event := &recaptchapb.Event{
		Token:   token,
		SiteKey: cfg.RecaptchaKey,
	}

	assessment := &recaptchapb.Assessment{
		Event: event,
	}

	// Build the assessment request.
	request := &recaptchapb.CreateAssessmentRequest{
		Assessment: assessment,
		Parent:     fmt.Sprintf("projects/%s", cfg.RecaptchaProject),
	}

	response, err := cli.CreateAssessment(
		ctx,
		request,
	)

	if err != nil {
		fmt.Printf("Error calling CreateAssessment: %v", err.Error())
		return err
	}

	// Check if the token is valid.
	if !response.TokenProperties.Valid {
		return fmt.Errorf(
			"call failed because the token was invalid for the following reasons: %v",
			response.TokenProperties.InvalidReason,
		)
	}

	// Check if the expected action was executed.
	if response.TokenProperties.Action != cfg.RecaptchaAction {
		return fmt.Errorf("action attribute in your reCAPTCHA tag does not match the action you are expecting to score")
	}

	// Get the risk score and the reason(s).
	// For more information on interpreting the assessment, see:
	// https://cloud.google.com/recaptcha-enterprise/docs/interpret-assessment
	// fmt.Printf("The reCAPTCHA score for this token is:  %v", response.RiskAnalysis.Score)

	// for _, reason := range response.RiskAnalysis.Reasons {
	// 	fmt.Printf(reason.String() + "\n")
	// }
	return nil
}
