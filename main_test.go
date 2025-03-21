package main

import (
	"testing"

	"buf.build/go/bufplugin/check/checktest"
)

func TestSpec(t *testing.T) {
	t.Parallel()
	checktest.SpecTest(t, spec)
}

func TestFailure(t *testing.T) {
	t.Parallel()

	checktest.CheckTest{
		Request: &checktest.RequestSpec{
			Options: map[string]any{
				"locked_services": []string{"test.LockedService"},
			},
			Files: &checktest.ProtoFileSpec{
				DirPaths: []string{
					"testdata/failure/current",
				},
				FilePaths: []string{
					"service.proto",
				},
			},
			AgainstFiles: &checktest.ProtoFileSpec{
				DirPaths: []string{
					"testdata/failure/previous",
				},
				FilePaths: []string{
					"service.proto",
				},
			},
		},
		Spec: spec,
		ExpectedAnnotations: []checktest.ExpectedAnnotation{
			{
				RuleID:  lockedServiceRuleID,
				Message: "Method \"MethodTwo\" on locked service \"test.LockedService\" is new.",
				FileLocation: &checktest.ExpectedFileLocation{
					FileName:    "service.proto",
					StartLine:   6,
					StartColumn: 2,
					EndLine:     6,
					EndColumn:   62,
				},
				AgainstFileLocation: &checktest.ExpectedFileLocation{
					FileName:    "service.proto",
					StartLine:   4,
					StartColumn: 0,
					EndLine:     6,
					EndColumn:   1,
				},
			},
		},
	}.Run(t)

	checktest.CheckTest{
		Request: &checktest.RequestSpec{
			Options: map[string]any{},
			Files: &checktest.ProtoFileSpec{
				DirPaths: []string{
					"testdata/failure/current",
				},
				FilePaths: []string{
					"service.proto",
				},
			},
			AgainstFiles: &checktest.ProtoFileSpec{
				DirPaths: []string{
					"testdata/failure/previous",
				},
				FilePaths: []string{
					"service.proto",
				},
			},
		},
		Spec:                spec,
		ExpectedAnnotations: []checktest.ExpectedAnnotation{},
	}.Run(t)

}
