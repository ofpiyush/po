package interpreter

import "testing"

func TestSimpleAssign(t *testing.T) {
	global := NewScope("root", nil)
	ast := &Assign{
		LHS: &Object{
			Name:  "yo yo",
			Type:  INT_TYPE,
			Value: 20,
		},
		RHS: &Object{
			Type:  STRING_TYPE,
			Value: "honey singh",
		},
	}
	_, err := Eval(global, ast)
	if err != nil {
		t.Error(err)
	}

	obj := global.Lookup("yo yo")
	if obj.Type != STRING_TYPE {
		t.Error("Expected assigned object type to be STRING_TYPE, got", TypeName[obj.Type])
	}
	if obj.Value != "honey singh" {
		t.Error("Expected assigned object value to be 'honey singh', got", obj.Value)
	}

}

func TestSimpleLookup(t *testing.T) {

	global := NewScope("root", nil)
	global.Update("yo yo", &Object{Type: STRING_TYPE, Value: "Honey singh"})
	ast := &Lookup{Name: "yo yo"}
	obj, err := Eval(global, ast)
	if err != nil {
		t.Error(err)
	}

	if obj.Type != STRING_TYPE {
		t.Error("Expected assigned object type to be STRING_TYPE, got", TypeName[obj.Type])
	}
	if obj.Value != "Honey singh" {
		t.Error("Expected assigned object value to be 'Honey singh', got", obj.Value)
	}
}

func TestSimpleConcatExpr(t *testing.T) {
	global := NewScope("root", nil)
	ast := &Concat{
		Source: []Expression{
			&Object{
				Type:  STRING_TYPE,
				Value: "Yahan se ",
			},
			&Object{
				Type:  INT_TYPE,
				Value: 50,
			},
			&Object{
				Type:  STRING_TYPE,
				Value: " kos door",
			},
		},
	}
	obj, err := Eval(global, ast)
	if err != nil {
		t.Error(err)
	}
	if obj.Type != STRING_TYPE {
		t.Error("Expected assigned object type to be STRING_TYPE, got", TypeName[obj.Type])
	}
	if obj.Value != "Yahan se 50 kos door" {
		t.Error("Expected assigned object value to be 'Yahan se 50 kos door', got", obj.Value)
	}

}

func TestBlockExpression(t *testing.T) {
	global := NewScope("root", nil)
	ast := &BlockExpr{
		Expressions: []Expression{
			&Assign{
				LHS: &Object{Name: "mom"},
				RHS: &Object{Type: STRING_TYPE, Value: "maa"},
			},
			&Assign{
				LHS: &Object{Name: "pachas"},
				RHS: &Object{Type: INT_TYPE, Value: 50},
			},
			&Assign{
				LHS: &Object{Name: "s1"},
				RHS: &Concat{
					Source: []Expression{
						&Object{Type: STRING_TYPE, Value: "Yahan se "},
						&Lookup{Name: "pachas"},
						&Object{Type: STRING_TYPE, Value: " kos door "},
					},
				},
			},
			&Assign{
				LHS: &Object{Name: "s2"},
				RHS: &Concat{
					Source: []Expression{
						&Object{Type: STRING_TYPE, Value: "jab koi "},
						&Object{Type: STRING_TYPE, Value: "mai ka laal "},
						&Object{Type: STRING_TYPE, Value: "rota hai "},
					},
				},
			},
			&Assign{
				LHS: &Object{Name: "s3"},
				RHS: &Concat{
					Source: []Expression{
						&Object{Type: STRING_TYPE, Value: "tab uski "},
						&Lookup{Name: "mom"},
						&Object{Type: STRING_TYPE, Value: " kehti hai, "},
					},
				},
			},
			&Assign{
				LHS: &Object{Name: "s4"},
				RHS: &Concat{
					Source: []Expression{
						&Lookup{Name: "s1"},
						&Lookup{Name: "s2"},
						&Lookup{Name: "s3"},
						&Object{Type: STRING_TYPE, Value: "beta so jaa "},
						&Object{Type: STRING_TYPE, Value: "nahi to Gabbar singh aa jaayega"},
					},
				},
			},
		},
	}
	_, err := Eval(global, ast)
	if err != nil {
		t.Error(err)
	}

	s := global.Lookup("s4")
	if s.Type != STRING_TYPE {
		t.Error("Expected assigned object type to be STRING_TYPE, got", TypeName[s.Type])
	}
	if s.Value != "Yahan se 50 kos door jab koi mai ka laal rota hai tab uski maa kehti hai, beta so jaa nahi to Gabbar singh aa jaayega" {
		t.Error("Expected assigned object value to be 'Yahan se 50 kos door jab koi bachha rota hai tab uski maa kehti hai, beta so jaa nahi to Gabbar singh aa jaayega', got", s.Value)
	}

}

func TestIfAsStatement(t *testing.T) {
	global := NewScope("root", nil)
	ast := &If{
		Condition: &Lookup{Name: "gabbar_or_samba"},
		Success: &BlockExpr{
			Expressions: []Expression{&Assign{
				LHS: &Object{Name: "dialogue"},
				RHS: &Object{Type: STRING_TYPE, Value: "Kitne aadmi the?"},
			}},
		},
		Fail: &BlockExpr{
			Expressions: []Expression{&Assign{
				LHS: &Object{Name: "dialogue"},
				RHS: &Object{Type: STRING_TYPE, Value: "Do sarkaar"},
			}},
		},
	}
	TrueInt := &Object{Type: INT_TYPE, Value: 1}
	FalseInt := &Object{Type: INT_TYPE, Value: 0}
	TrueStr := &Object{Type: STRING_TYPE, Value: "yes"}
	FalseStr := &Object{Type: STRING_TYPE, Value: ""}

	global.Update("gabbar_or_samba", TrueInt)
	Eval(global, ast)

	k := global.Lookup("dialogue")
	if k.Type != STRING_TYPE {
		t.Error("Expected assigned object type to be STRING_TYPE, got", TypeName[k.Type])
	}
	if k.Value != "Kitne aadmi the?" {
		t.Error("Expected assigned object value to be 'Kitne aadmi the?', got", k.Value)
	}

	global.Update("gabbar_or_samba", FalseInt)
	Eval(global, ast)

	k = global.Lookup("dialogue")
	if k.Type != STRING_TYPE {
		t.Error("Expected assigned object type to be STRING_TYPE, got", TypeName[k.Type])
	}
	if k.Value != "Do sarkaar" {
		t.Error("Expected assigned object value to be 'Do sarkaar', got", k.Value)
	}

	global.Update("gabbar_or_samba", TrueStr)
	Eval(global, ast)

	k = global.Lookup("dialogue")
	if k.Type != STRING_TYPE {
		t.Error("Expected assigned object type to be STRING_TYPE, got", TypeName[k.Type])
	}
	if k.Value != "Kitne aadmi the?" {
		t.Error("Expected assigned object value to be 'Kitne aadmi the?', got", k.Value)
	}

	global.Update("gabbar_or_samba", FalseStr)
	Eval(global, ast)

	k = global.Lookup("dialogue")
	if k.Type != STRING_TYPE {
		t.Error("Expected assigned object type to be STRING_TYPE, got", TypeName[k.Type])
	}
	if k.Value != "Do sarkaar" {
		t.Error("Expected assigned object value to be 'Do sarkaar', got", k.Value)
	}
}

func TestLambda(t *testing.T) {
	global := NewScope("root", nil)
	global.Update("pachas", &Object{Type: INT_TYPE, Value: 50})

	ast := &Lambda{
		Params: map[string]Expression{
			"mom": &Object{Type: STRING_TYPE, Value: "maa"},
			"s1": &Concat{
				Source: []Expression{
					&Object{Type: STRING_TYPE, Value: "Yahan se "},
					&Lookup{Name: "pachas"},
					&Object{Type: STRING_TYPE, Value: " kos door "},
				},
			},
		},
		Body: &BlockExpr{
			Expressions: []Expression{
				&Assign{
					LHS: &Object{Name: "s2"},
					RHS: &Concat{
						Source: []Expression{
							&Object{Type: STRING_TYPE, Value: "jab koi "},
							&Object{Type: STRING_TYPE, Value: "mai ka laal "},
							&Object{Type: STRING_TYPE, Value: "rota hai "},
						},
					},
				},
				&Assign{
					LHS: &Object{Name: "s3"},
					RHS: &Concat{
						Source: []Expression{
							&Object{Type: STRING_TYPE, Value: "tab uski "},
							&Lookup{Name: "mom"},
							&Object{Type: STRING_TYPE, Value: " kehti hai, "},
						},
					},
				},
				&Return{
					Expr: &Concat{
						Source: []Expression{
							&Lookup{Name: "s1"},
							&Lookup{Name: "s2"},
							&Lookup{Name: "s3"},
							&Object{Type: STRING_TYPE, Value: "beta so jaa "},
							&Object{Type: STRING_TYPE, Value: "nahi to Gabbar singh aa jaayega"},
						},
					},
				},
			},
		},
	}
	out, err := Eval(global, ast)
	if err != nil {
		t.Fatal(err)
	}

	if out.Type != STRING_TYPE {
		t.Error("Expected assigned object type to be STRING_TYPE, got", TypeName[out.Type])
	}
	if out.Value != "Yahan se 50 kos door jab koi mai ka laal rota hai tab uski maa kehti hai, beta so jaa nahi to Gabbar singh aa jaayega" {
		t.Error("Expected assigned object value to be 'Yahan se 50 kos door jab koi bachha rota hai tab uski maa kehti hai, beta so jaa nahi to Gabbar singh aa jaayega', got", out.Value)
	}

}

func TestInfix(t *testing.T) {

}
