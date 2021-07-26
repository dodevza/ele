package versionparsers

import "testing"

func Test_SemanticVersion_Matches(t *testing.T) {
	vp := NewSemanticParser("V")

	_, isVersion := vp.Parse("V1.0.0")

	if isVersion == false {
		t.Fatalf("Did not pass as version")
	}
}

func Test_SemanticVersionWithBetaRelease_Matches(t *testing.T) {
	vp := NewSemanticParser("V")

	_, isVersion := vp.Parse("V1.0.0-Beta1")

	if isVersion == false {
		t.Fatalf("Did not pass as version")
	}
}

func Test_NonVersion_DoesnotMatches(t *testing.T) {
	vp := NewSemanticParser("V")

	_, isVersion := vp.Parse("Vehicle")

	if isVersion == true {
		t.Fatalf("Vehicle is not a semantic version")
	}
}

func Test_SingleDigitVersion_Matches(t *testing.T) {
	vp := NewSemanticParser("V")

	_, isVersion := vp.Parse("V1")

	if isVersion == false {
		t.Fatalf("Single digit did not match")
	}
}

func Test_PrefixWithPeriod_Matches(t *testing.T) {
	vp := NewSemanticParser("V.")

	_, isVersion := vp.Parse("V.1")

	if isVersion == false {
		t.Fatalf("Prefix with period did not match")
	}
}

func Test_PrefixWithPeriodNonSemanticVersion_DoesNotMatch(t *testing.T) {
	vp := NewSemanticParser("V.")

	_, isVersion := vp.Parse("Vehicles")

	if isVersion == true {
		t.Fatalf("Vehicles is not a sematic mersion")
	}
}
