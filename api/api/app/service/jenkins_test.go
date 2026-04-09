package service

import (
	"strings"
	"testing"
)

func TestDecodeJobParameters(t *testing.T) {
	jobName, parameters, err := decodeJobParameters(strings.NewReader(`{
		"name": "opsnexus-mock-job",
		"property": [
			{
				"_class": "hudson.model.ParametersDefinitionProperty",
				"parameterDefinitions": [
					{
						"name": "BRANCH",
						"type": "StringParameterDefinition",
						"description": "git branch",
						"defaultParameterValue": {
							"value": "main"
						}
					},
					{
						"name": "DRY_RUN",
						"type": "BooleanParameterDefinition",
						"description": "dry run toggle",
						"defaultParameterValue": {
							"value": true
						}
					},
					{
						"name": "TARGET_ENV",
						"type": "ChoiceParameterDefinition",
						"description": "target env",
						"choices": ["test", "prod"]
					}
				]
			}
		]
	}`))
	if err != nil {
		t.Fatalf("decodeJobParameters() error = %v", err)
	}

	if jobName != "opsnexus-mock-job" {
		t.Fatalf("decodeJobParameters() jobName = %q, want %q", jobName, "opsnexus-mock-job")
	}

	if len(parameters) != 3 {
		t.Fatalf("decodeJobParameters() len(parameters) = %d, want 3", len(parameters))
	}

	if parameters[0].Name != "BRANCH" || parameters[0].DefaultValue != "main" {
		t.Fatalf("decodeJobParameters() first parameter = %#v", parameters[0])
	}

	if parameters[1].Name != "DRY_RUN" || parameters[1].DefaultValue != true {
		t.Fatalf("decodeJobParameters() second parameter = %#v", parameters[1])
	}

	if parameters[2].Name != "TARGET_ENV" {
		t.Fatalf("decodeJobParameters() third parameter name = %q, want %q", parameters[2].Name, "TARGET_ENV")
	}

	if got, want := strings.Join(parameters[2].Choices, ","), "test,prod"; got != want {
		t.Fatalf("decodeJobParameters() third parameter choices = %q, want %q", got, want)
	}
}

func TestDecodeJobParametersNormalizesChoiceText(t *testing.T) {
	_, parameters, err := decodeJobParameters(strings.NewReader(`{
		"name": "opsnexus-mock-job",
		"property": [
			{
				"_class": "hudson.model.ParametersDefinitionProperty",
				"parameterDefinitions": [
					{
						"name": "TARGET_ENV",
						"type": "ChoiceParameterDefinition",
						"choices": "test\nstaging\nprod\n"
					}
				]
			}
		]
	}`))
	if err != nil {
		t.Fatalf("decodeJobParameters() error = %v", err)
	}

	if len(parameters) != 1 {
		t.Fatalf("decodeJobParameters() len(parameters) = %d, want 1", len(parameters))
	}

	if got, want := strings.Join(parameters[0].Choices, ","), "test,staging,prod"; got != want {
		t.Fatalf("decodeJobParameters() choices = %q, want %q", got, want)
	}
}
