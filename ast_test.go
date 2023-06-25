package boule

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestBinaryExpression_Evaluate(t *testing.T) {

	t.Run("evaluate bool <-> bool", func(t *testing.T) {

		t.Run("equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right string) bool {
				result, err := (&BinaryExpression{
					token: EQUAL,
					left: &LiteralIdent{
						identifier: "true",
					},
					right: &LiteralIdent{
						identifier: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, "true"))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, "false"))
			})
		})

		t.Run("not equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right string) bool {
				result, err := (&BinaryExpression{
					token: NOT_EQUAL,
					left: &LiteralIdent{
						identifier: "true",
					},
					right: &LiteralIdent{
						identifier: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, "false"))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, "true"))
			})
		})
	})

	t.Run("evaluate string <-> string", func(t *testing.T) {

		t.Run("equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right string) bool {
				result, err := (&BinaryExpression{
					token: EQUAL,
					left: &LiteralString{
						value: "mars",
					},
					right: &LiteralString{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, "mars"))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, "neptune"))
			})
		})

		t.Run("not equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right string) bool {
				result, err := (&BinaryExpression{
					token: NOT_EQUAL,
					left: &LiteralString{
						value: "mars",
					},
					right: &LiteralString{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, "neptune"))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, "mars"))
			})
		})
	})

	t.Run("evaluate float64 <-> float64", func(t *testing.T) {

		t.Run("equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: EQUAL,
					left: &LiteralFloat{
						value: 481.4,
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, 481.4))
			})
			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, 892673.8976))
			})
			t.Run("false if greater", func(t *testing.T) {
				assert.False(t, evaluate(t, 22.4))
			})
		})

		t.Run("not equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: NOT_EQUAL,
					left: &LiteralFloat{
						value: 481.4,
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if less", func(t *testing.T) {
				assert.True(t, evaluate(t, 892673.8976))
			})
			t.Run("false if equal", func(t *testing.T) {
				assert.False(t, evaluate(t, 481.4))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, 22.4))
			})
		})

		t.Run("less", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: LESS,
					left: &LiteralFloat{
						value: 481.4,
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if less", func(t *testing.T) {
				assert.True(t, evaluate(t, 9087235.98))
			})
			t.Run("false if equal", func(t *testing.T) {
				assert.False(t, evaluate(t, 481.4))
			})
			t.Run("false if greater", func(t *testing.T) {
				assert.False(t, evaluate(t, 234.9))
			})
		})

		t.Run("less or equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: LESS_OR_EQUAL,
					left: &LiteralFloat{
						value: 481.4,
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if less", func(t *testing.T) {
				assert.True(t, evaluate(t, 9087235.98))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, 481.4))
			})
			t.Run("false if greater", func(t *testing.T) {
				assert.False(t, evaluate(t, 234.9))
			})
		})

		t.Run("greater", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: GREATER,
					left: &LiteralFloat{
						value: 3245635.456,
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, 9087235.98))
			})
			t.Run("false if equal", func(t *testing.T) {
				assert.False(t, evaluate(t, 3245635.456))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, 234.9))
			})
		})

		t.Run("greater or equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: GREATER_OR_EQUAL,
					left: &LiteralFloat{
						value: 3245635.456,
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, 9087235.98))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, 3245635.456))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, 234.9))
			})
		})

		t.Run("negative float", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: GREATER_OR_EQUAL,
					left: &LiteralFloat{
						value: -234.4,
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, 9087235.98))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, -234.4))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, -9872635.234))
			})
		})
	})

	t.Run("evaluate big.Int <-> big.Int", func(t *testing.T) {

		t.Run("equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right *big.Int) bool {
				result, err := (&BinaryExpression{
					token: EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(983467234523),
					},
					right: &LiteralInteger{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(983467234523)))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, big.NewInt(23423)))
			})
		})

		t.Run("not equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right *big.Int) bool {
				result, err := (&BinaryExpression{
					token: NOT_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(983467234523),
					},
					right: &LiteralInteger{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(23423)))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, big.NewInt(983467234523)))
			})
		})

		t.Run("less", func(t *testing.T) {

			evaluate := func(t *testing.T, right *big.Int) bool {
				result, err := (&BinaryExpression{
					token: LESS,
					left: &LiteralInteger{
						value: big.NewInt(345),
					},
					right: &LiteralInteger{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if less", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(23423)))
			})
			t.Run("false if equal", func(t *testing.T) {
				assert.False(t, evaluate(t, big.NewInt(345)))
			})
			t.Run("false if greater", func(t *testing.T) {
				assert.False(t, evaluate(t, big.NewInt(34)))
			})
		})

		t.Run("less or equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right *big.Int) bool {
				result, err := (&BinaryExpression{
					token: LESS_OR_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(23423),
					},
					right: &LiteralInteger{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if less", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(983467234523)))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(23423)))
			})
			t.Run("false if greater", func(t *testing.T) {
				assert.False(t, evaluate(t, big.NewInt(34)))
			})
		})

		t.Run("greater", func(t *testing.T) {

			evaluate := func(t *testing.T, right *big.Int) bool {
				result, err := (&BinaryExpression{
					token: GREATER,
					left: &LiteralInteger{
						value: big.NewInt(234234),
					},
					right: &LiteralInteger{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, big.NewInt(235456)))
			})
			t.Run("false if equal", func(t *testing.T) {
				assert.False(t, evaluate(t, big.NewInt(234234)))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(34)))
			})
		})

		t.Run("greater or equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right *big.Int) bool {
				result, err := (&BinaryExpression{
					token: GREATER_OR_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(235456),
					},
					right: &LiteralInteger{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, big.NewInt(983467234523)))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(235456)))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(34)))
			})
		})

		t.Run("negative float", func(t *testing.T) {

			evaluate := func(t *testing.T, right *big.Int) bool {
				result, err := (&BinaryExpression{
					token: GREATER_OR_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(-235456),
					},
					right: &LiteralInteger{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, big.NewInt(983467234523)))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(-235456)))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, big.NewInt(-2354562222)))
			})
		})
	})

	t.Run("evaluate big.Int <-> float64", func(t *testing.T) {

		t.Run("equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(983467234523),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, 983467234523.0))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, 983467234522.0))
				assert.False(t, evaluate(t, 983467234522.3))
				assert.False(t, evaluate(t, 983467234522.7))
				assert.False(t, evaluate(t, 983467234523.3))
				assert.False(t, evaluate(t, 983467234523.7))
				assert.False(t, evaluate(t, 983467234524.0))
			})
		})

		t.Run("not equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: NOT_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(983467234523),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, 983467234522.0))
				assert.True(t, evaluate(t, 983467234522.3))
				assert.True(t, evaluate(t, 983467234522.7))
				assert.True(t, evaluate(t, 983467234523.3))
				assert.True(t, evaluate(t, 983467234523.7))
				assert.True(t, evaluate(t, 983467234524.0))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, 983467234523.0))
			})
		})

		t.Run("less", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: LESS,
					left: &LiteralInteger{
						value: big.NewInt(345),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if less", func(t *testing.T) {
				assert.True(t, evaluate(t, 345.3))
				assert.True(t, evaluate(t, 345.7))
				assert.True(t, evaluate(t, 346.0))
			})
			t.Run("false if equal", func(t *testing.T) {
				assert.False(t, evaluate(t, 345.0))
			})
			t.Run("false if greater", func(t *testing.T) {
				assert.False(t, evaluate(t, 344.0))
				assert.False(t, evaluate(t, 344.3))
				assert.False(t, evaluate(t, 344.7))
			})
		})

		t.Run("less or equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: LESS_OR_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(23423),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if less", func(t *testing.T) {
				assert.True(t, evaluate(t, 23423.3))
				assert.True(t, evaluate(t, 23423.7))
				assert.True(t, evaluate(t, 23424.0))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, 23423.0))
			})
			t.Run("false if greater", func(t *testing.T) {
				assert.False(t, evaluate(t, 23422.0))
				assert.False(t, evaluate(t, 23422.3))
				assert.False(t, evaluate(t, 23422.7))
			})
		})

		t.Run("greater", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: GREATER,
					left: &LiteralInteger{
						value: big.NewInt(234234),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, 234234.3))
				assert.False(t, evaluate(t, 234234.7))
				assert.False(t, evaluate(t, 234235.0))
			})
			t.Run("false if equal", func(t *testing.T) {
				assert.False(t, evaluate(t, 234234.0))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, 234233.0))
				assert.True(t, evaluate(t, 234233.3))
				assert.True(t, evaluate(t, 234233.7))
			})
		})

		t.Run("greater or equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: GREATER_OR_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(235456),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, 235456.3))
				assert.False(t, evaluate(t, 235456.7))
				assert.False(t, evaluate(t, 235457.0))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, 235456.0))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, 235455.0))
				assert.True(t, evaluate(t, 235455.3))
				assert.True(t, evaluate(t, 235455.7))
			})
		})

		t.Run("negative float", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: GREATER_OR_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(-235456),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, -235455.0))
				assert.False(t, evaluate(t, -235455.3))
				assert.False(t, evaluate(t, -235455.7))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, -235456.0))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, -235456.3))
				assert.True(t, evaluate(t, -235456.7))
				assert.True(t, evaluate(t, -235457.0))
			})
		})
	})

	t.Run("evaluate float64 <-> big.Int", func(t *testing.T) {

		t.Run("equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(983467234523),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, 983467234523.0))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, 983467234522.0))
				assert.False(t, evaluate(t, 983467234522.3))
				assert.False(t, evaluate(t, 983467234522.7))
				assert.False(t, evaluate(t, 983467234523.3))
				assert.False(t, evaluate(t, 983467234523.7))
				assert.False(t, evaluate(t, 983467234524.0))
			})
		})

		t.Run("not equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: NOT_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(983467234523),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true", func(t *testing.T) {
				assert.True(t, evaluate(t, 983467234522.0))
				assert.True(t, evaluate(t, 983467234522.3))
				assert.True(t, evaluate(t, 983467234522.7))
				assert.True(t, evaluate(t, 983467234523.3))
				assert.True(t, evaluate(t, 983467234523.7))
				assert.True(t, evaluate(t, 983467234524.0))
			})
			t.Run("false", func(t *testing.T) {
				assert.False(t, evaluate(t, 983467234523.0))
			})
		})

		t.Run("less", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: LESS,
					left: &LiteralInteger{
						value: big.NewInt(345),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if less", func(t *testing.T) {
				assert.True(t, evaluate(t, 345.3))
				assert.True(t, evaluate(t, 345.7))
				assert.True(t, evaluate(t, 346.0))
			})
			t.Run("false if equal", func(t *testing.T) {
				assert.False(t, evaluate(t, 345.0))
			})
			t.Run("false if greater", func(t *testing.T) {
				assert.False(t, evaluate(t, 344.0))
				assert.False(t, evaluate(t, 344.3))
				assert.False(t, evaluate(t, 344.7))
			})
		})

		t.Run("less or equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: LESS_OR_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(23423),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("true if less", func(t *testing.T) {
				assert.True(t, evaluate(t, 23423.3))
				assert.True(t, evaluate(t, 23423.7))
				assert.True(t, evaluate(t, 23424.0))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, 23423.0))
			})
			t.Run("false if greater", func(t *testing.T) {
				assert.False(t, evaluate(t, 23422.0))
				assert.False(t, evaluate(t, 23422.3))
				assert.False(t, evaluate(t, 23422.7))
			})
		})

		t.Run("greater", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: GREATER,
					left: &LiteralInteger{
						value: big.NewInt(234234),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, 234234.3))
				assert.False(t, evaluate(t, 234234.7))
				assert.False(t, evaluate(t, 234235.0))
			})
			t.Run("false if equal", func(t *testing.T) {
				assert.False(t, evaluate(t, 234234.0))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, 234233.0))
				assert.True(t, evaluate(t, 234233.3))
				assert.True(t, evaluate(t, 234233.7))
			})
		})

		t.Run("greater or equal", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: GREATER_OR_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(235456),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, 235456.3))
				assert.False(t, evaluate(t, 235456.7))
				assert.False(t, evaluate(t, 235457.0))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, 235456.0))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, 235455.0))
				assert.True(t, evaluate(t, 235455.3))
				assert.True(t, evaluate(t, 235455.7))
			})
		})

		t.Run("negative float", func(t *testing.T) {

			evaluate := func(t *testing.T, right float64) bool {
				result, err := (&BinaryExpression{
					token: GREATER_OR_EQUAL,
					left: &LiteralInteger{
						value: big.NewInt(-235456),
					},
					right: &LiteralFloat{
						value: right,
					},
				}).Evaluate(NewData())
				assert.NoError(t, err)
				return result.(bool)
			}

			t.Run("false if less", func(t *testing.T) {
				assert.False(t, evaluate(t, -235455.0))
				assert.False(t, evaluate(t, -235455.3))
				assert.False(t, evaluate(t, -235455.7))
			})
			t.Run("true if equal", func(t *testing.T) {
				assert.True(t, evaluate(t, -235456.0))
			})
			t.Run("true if greater", func(t *testing.T) {
				assert.True(t, evaluate(t, -235456.3))
				assert.True(t, evaluate(t, -235456.7))
				assert.True(t, evaluate(t, -235457.0))
			})
		})
	})
}
