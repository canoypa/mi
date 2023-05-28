package utils

import "github.com/AlecAivazis/survey/v2"

func Confirm(Message string, Default bool) bool {
	result := Default

	prompt := &survey.Confirm{
		Message: Message,
		Default: Default,
	}
	survey.AskOne(prompt, &result)

	return result
}

func Input(Message string) string {
	result := ""

	prompt := &survey.Input{
		Message: Message,
	}
	survey.AskOne(prompt, &result)

	return result
}

func Multiline(Message string) string {
	result := ""

	prompt := &survey.Multiline{
		Message: Message,
	}
	survey.AskOne(prompt, &result)

	return result
}
