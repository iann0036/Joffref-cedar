package cedar

import (
	"context"
	"testing"
)

func TestCedarEngine_IsAuthorized(t *testing.T) {
	policy := `
	permit(
		principal == User::"alice",
		action    == Action::"update",
		resource  == Photo::"VacationPhoto94.jpg"
	);
	`
	engine, err := NewCedarEngine(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close(context.Background())
	err = engine.SetEntitiesFromJson(context.Background(), "[]")
	if err != nil {
		t.Fatal(err)
	}
	err = engine.SetPolicies(context.Background(), policy)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("is authorized must return allow", func(t *testing.T) {
		isAuthorizedMustReturnAllow(t, engine, "User::\"alice\"", "Action::\"update\"", "Photo::\"VacationPhoto94.jpg\"")
	})
	t.Run("is authorized must return deny", func(t *testing.T) {
		isAuthorizedMustReturnDeny(t, engine, "User::\"alice\"", "Action::\"update\"", "Photo::\"VacationPhoto95.jpg\"")
	})
}
func TestCedarEngine_IsAuthorizedJson(t *testing.T) {
	policy := `
	permit(
		principal == User::"alice",
		action    == Action::"update",
		resource  == Photo::"VacationPhoto94.jpg"
	);
	`
	engine, err := NewCedarEngine(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close(context.Background())
	err = engine.SetEntitiesFromJson(context.Background(), "[]")
	if err != nil {
		t.Fatal(err)
	}
	err = engine.SetPolicies(context.Background(), policy)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("is authorized json", func(t *testing.T) {
		res, err := engine.IsAuthorizedJson(context.Background(), EvalRequest{
			Principal: "User::\"alice\"",
			Action:    "Action::\"update\"",
			Resource:  "Photo::\"VacationPhoto94.jpg\"",
			Context:   "{}",
		})
		if err != nil {
			t.Fatal(err)
		}
		if !res {
			t.Fatal("expected permit")
		}
	})
}
func TestCedarEngine_Validate(t *testing.T) {
	policy := `
	permit(
		principal,
		action,
		resource
	);
	`
	schema := `{
		"": {
			"entityTypes": {
				"User": {
					"shape": {
						"type": "Record",
						"attributes": {
							"department": { "type": "String" }
						}
					}
				},
				"Photo": {
					"shape": {
						"type": "Record",
						"attributes": {
							"private": { "type": "Boolean" }
						}
					}
				}
			},
			"actions": {
				"viewPhoto": {
					"appliesTo": {
						"principalTypes": [ "User" ],
						"resourceTypes": [ "Photo" ]
					}
				}
			}
		}
	}
	`
	engine, err := NewCedarEngine(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close(context.Background())
	err = engine.SetEntitiesFromJson(context.Background(), "[]")
	if err != nil {
		t.Fatal(err)
	}
	err = engine.SetPolicies(context.Background(), policy)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("validate pass", func(t *testing.T) {
		res, err := engine.Validate(context.Background(), schema, ValidationModePermissive)
		if err != nil {
			t.Fatal(err)
		}
		if !res.Passed {
			t.Fatal("expected validate pass but got failure")
		}
	})
	invalidPolicy := `
	permit(
		principal == Banana::"invalid",
		resource  == Apple::"foo"
	);
	`
	err = engine.SetPolicies(context.Background(), invalidPolicy)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("validate fail", func(t *testing.T) {
		res, err := engine.Validate(context.Background(), schema, ValidationModeStrict)
		if err != nil {
			t.Fatal(err)
		}
		if res.Passed {
			t.Fatal("expected validate failure but got pass")
		}
	})
}
func isAuthorizedMustReturnAllow(t *testing.T, engine *CedarEngine, principal, action, resource string) {
	res, err := engine.IsAuthorized(context.Background(), EvalRequest{
		Principal: principal,
		Action:    action,
		Resource:  resource,
		Context:   "{}",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("expected permit")
	}
}

func isAuthorizedMustReturnDeny(t *testing.T, engine *CedarEngine, principal, action, resource string) {
	res, err := engine.IsAuthorized(context.Background(), EvalRequest{
		Principal: principal,
		Action:    action,
		Resource:  resource,
		Context:   "{}",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res {
		t.Fatal("expected deny")
	}
}
