package appprotect

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestValidateRequiredFields(t *testing.T) {
	tests := []struct {
		obj        *unstructured.Unstructured
		fieldsList [][]string
		expectErr  bool
		msg        string
	}{
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": map[string]interface{}{},
					"b": map[string]interface{}{},
				},
			},
			fieldsList: [][]string{{"a"}, {"b"}},
			expectErr:  false,
			msg:        "valid object with 2 fields",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": map[string]interface{}{},
				},
			},
			fieldsList: [][]string{{"a"}, {"b"}},
			expectErr:  true,
			msg:        "invalid object with a missing field",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": map[string]interface{}{},
					"x": map[string]interface{}{},
				},
			},
			fieldsList: [][]string{{"a"}, {"b"}},
			expectErr:  true,
			msg:        "invalid object with a wrong field",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{},
					},
				},
			},
			fieldsList: [][]string{{"a", "b"}},
			expectErr:  false,
			msg:        "valid object with nested field",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": map[string]interface{}{
						"x": map[string]interface{}{},
					},
				},
			},
			fieldsList: [][]string{{"a", "b"}},
			expectErr:  true,
			msg:        "invalid object with a wrong nested field",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			fieldsList: nil,
			expectErr:  false,
			msg:        "valid object with no validation",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": "wrong-type", // must be map[string]interface{}
				},
			},
			fieldsList: [][]string{{"a"}},
			expectErr:  true,
			msg:        "invalid object with a field of wrong type",
		},
	}

	for _, test := range tests {
		err := validateRequiredFields(test.obj, test.fieldsList)
		if test.expectErr && err == nil {
			t.Errorf("validateRequiredFields() returned no error for the case of %s", test.msg)
		}
		if !test.expectErr && err != nil {
			t.Errorf("validateRequiredFields() returned unexpected error %v for the case of %s", err, test.msg)
		}
	}
}

func TestValidateRequiredSlices(t *testing.T) {
	tests := []struct {
		obj        *unstructured.Unstructured
		fieldsList [][]string
		expectErr  bool
		msg        string
	}{
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": []interface{}{},
					"b": []interface{}{},
				},
			},
			fieldsList: [][]string{{"a"}, {"b"}},
			expectErr:  false,
			msg:        "valid object with 2 fields",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": []interface{}{},
				},
			},
			fieldsList: [][]string{{"a"}, {"b"}},
			expectErr:  true,
			msg:        "invalid object with a field",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": []interface{}{},
					"x": []interface{}{},
				},
			},
			fieldsList: [][]string{{"a"}, {"b"}},
			expectErr:  true,
			msg:        "invalid object with a wrong field",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": map[string]interface{}{
						"b": []interface{}{},
					},
				},
			},
			fieldsList: [][]string{{"a", "b"}},
			expectErr:  false,
			msg:        "valid object with nested field",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": map[string]interface{}{
						"x": []interface{}{},
					},
				},
			},
			fieldsList: [][]string{{"a", "b"}},
			expectErr:  true,
			msg:        "invalid object with a wrong nested field",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			fieldsList: nil,
			expectErr:  false,
			msg:        "valid object with no validation",
		},
		{
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"a": "wrong-type", // must be [string]interface{}
				},
			},
			fieldsList: [][]string{{"a"}},
			expectErr:  true,
			msg:        "invalid object with a field of wrong type",
		},
	}

	for _, test := range tests {
		err := validateRequiredSlices(test.obj, test.fieldsList)
		if test.expectErr && err == nil {
			t.Errorf("validateRequiredSlices() returned no error for the case of %s", test.msg)
		}
		if !test.expectErr && err != nil {
			t.Errorf("validateRequiredSlices() returned unexpected error %v for the case of %s", err, test.msg)
		}
	}
}

func TestValidateAppProtectPolicy(t *testing.T) {
	tests := []struct {
		policy    *unstructured.Unstructured
		expectErr bool
		msg       string
	}{
		{
			policy: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"policy": map[string]interface{}{},
					},
				},
			},
			expectErr: false,
			msg:       "valid policy",
		},
		{
			policy: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"something": map[string]interface{}{},
					},
				},
			},
			expectErr: true,
			msg:       "invalid policy with no policy field",
		},
		{
			policy: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"something": map[string]interface{}{
						"policy": map[string]interface{}{},
					},
				},
			},
			expectErr: true,
			msg:       "invalid policy with no spec field",
		},
	}

	for _, test := range tests {
		err := validateAppProtectPolicy(test.policy)
		if test.expectErr && err == nil {
			t.Errorf("validateAppProtectPolicy() returned no error for the case of %s", test.msg)
		}
		if !test.expectErr && err != nil {
			t.Errorf("validateAppProtectPolicy() returned unexpected error %v for the case of %s", err, test.msg)
		}
	}
}

func TestValidateAppProtectLogConf(t *testing.T) {
	tests := []struct {
		logConf   *unstructured.Unstructured
		expectErr bool
		msg       string
	}{
		{
			logConf: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"content": map[string]interface{}{},
						"filter":  map[string]interface{}{},
					},
				},
			},
			expectErr: false,
			msg:       "valid log conf",
		},
		{
			logConf: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"filter": map[string]interface{}{},
					},
				},
			},
			expectErr: true,
			msg:       "invalid log conf with no content field",
		},
		{
			logConf: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"content": map[string]interface{}{},
					},
				},
			},
			expectErr: true,
			msg:       "invalid log conf with no filter field",
		},
		{
			logConf: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"something": map[string]interface{}{
						"content": map[string]interface{}{},
						"filter":  map[string]interface{}{},
					},
				},
			},
			expectErr: true,
			msg:       "invalid log conf with no spec field",
		},
	}

	for _, test := range tests {
		err := validateAppProtectLogConf(test.logConf)
		if test.expectErr && err == nil {
			t.Errorf("validateAppProtectLogConf() returned no error for the case of %s", test.msg)
		}
		if !test.expectErr && err != nil {
			t.Errorf("validateAppProtectLogConf() returned unexpected error %v for the case of %s", err, test.msg)
		}
	}
}

func TestValidateAppProtectLogDestinationAnnotation(t *testing.T) {
	// Positive test cases
	var posDstAntns = []string{"stderr", "syslog:server=localhost:9000", "syslog:server=10.1.1.2:9000", "/var/log/ap.log"}

	// Negative test cases item, expected error message
	var negDstAntns = [][]string{
		{"stdout", "Log Destination did not follow format"},
		{"syslog:server=localhost:99999", "not a valid port number"},
		{"syslog:server=999.99.99.99:5678", "is not a valid ip address"},
	}

	for _, tCase := range posDstAntns {
		err := ValidateAppProtectLogDestination(tCase)
		if err != nil {
			t.Errorf("got %v expected nil", err)
		}
	}
	for _, nTCase := range negDstAntns {
		err := ValidateAppProtectLogDestination(nTCase[0])
		if err == nil {
			t.Errorf("got no error expected error containing %s", nTCase[1])
		} else {
			if !strings.Contains(err.Error(), nTCase[1]) {
				t.Errorf("got %v expected to contain: %s", err, nTCase[1])
			}
		}
	}
}

func TestValidateAppProtectUserSig(t *testing.T) {
	tests := []struct {
		userSig   *unstructured.Unstructured
		expectErr bool
		msg       string
	}{
		{
			userSig: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"signatures": []interface{}{},
					},
				},
			},
			expectErr: false,
			msg:       "valid user sig",
		},
		{
			userSig: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"something": []interface{}{},
					},
				},
			},
			expectErr: true,
			msg:       "invalid user sig with no signatures",
		},
		{
			userSig: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"something": map[string]interface{}{
						"signatures": []interface{}{},
					},
				},
			},
			expectErr: true,
			msg:       "invalid user sign with no spec field",
		},
	}

	for _, test := range tests {
		err := validateAppProtectUserSig(test.userSig)
		if test.expectErr && err == nil {
			t.Errorf("validateAppProtectUserSig() returned no error for the case of %s", test.msg)
		}
		if !test.expectErr && err != nil {
			t.Errorf("validateAppProtectUserSig() returned unexpected error %v for the case of %s", err, test.msg)
		}
	}
}

func TestParseResourceReferenceAnnotation(t *testing.T) {
	tests := []struct {
		ns, antn, expected string
	}{
		{
			ns:       "default",
			antn:     "resource",
			expected: "default/resource",
		},
		{
			ns:       "default",
			antn:     "ns-1/resource",
			expected: "ns-1/resource",
		},
	}

	for _, test := range tests {
		result := ParseResourceReferenceAnnotation(test.ns, test.antn)
		if result != test.expected {
			t.Errorf("ParseResourceReferenceAnnotation(%q,%q) returned %q but expected %q", test.ns, test.antn, result, test.expected)
		}
	}
}

func TestGenNsName(t *testing.T) {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"namespace": "default",
				"name":      "resource",
			},
		},
	}

	expected := "default/resource"

	result := GetNsName(obj)
	if result != expected {
		t.Errorf("GetNsName() returned %q but expected %q", result, expected)
	}
}
