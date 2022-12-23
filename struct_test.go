package gocast

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type customInt int

func (c *customInt) CastSet(v any) error {
	*c = customInt(Int(v))
	return nil
}

func TestStructGetSetFieldValue(t *testing.T) {
	st := &struct {
		Name        string
		Value       int64
		Count       customInt
		CreatedAt   time.Time
		UpdatedAt   time.Time
		PublishedAt time.Time
		AnyTarget   any
	}{}

	assert.NoError(t, SetStructFieldValue(&st, "Name", "TestName"), "set Name field value")
	assert.NoError(t, SetStructFieldValue(&st, "Value", int64(127)), "set Value field value")
	assert.NoError(t, SetStructFieldValue(&st, "Count", 3.1), "set Count field value")
	assert.NoError(t, SetStructFieldValue(&st, "CreatedAt", "2020/10/10"), "set CreatedAt field value")
	assert.NoError(t, SetStructFieldValue(&st, "UpdatedAt", time.Now().Unix()), "set UpdatedAt field value")
	assert.NoError(t, SetStructFieldValue(&st, "PublishedAt", time.Now()), "set PublishedAt field value")
	assert.NoError(t, SetStructFieldValue(&st, "AnyTarget", "hi"), "set AnyTarget field value")
	assert.Error(t, SetStructFieldValue(&st, "UndefinedField", int64(127)), "set UndefinedField field value must be error")

	name, err := StructFieldValue(st, "Name")
	assert.NoError(t, err, "get Name value")
	assert.Equal(t, "TestName", name)

	value, err := StructFieldValue(st, "Value")
	assert.NoError(t, err, "get Value value")
	assert.Equal(t, int64(127), value)

	count, err := StructFieldValue(st, "Count")
	assert.NoError(t, err, "get Count value")
	assert.Equal(t, customInt(3), count)

	createdAt, err := StructFieldValue(st, "CreatedAt")
	assert.NoError(t, err, "get CreatedAt value")
	assert.Equal(t, 2020, createdAt.(time.Time).Year())

	updatedAt, err := StructFieldValue(st, "UpdatedAt")
	assert.NoError(t, err, "get UpdatedAt value")
	assert.Equal(t, time.Now().Year(), updatedAt.(time.Time).Year())

	publishedAt, err := StructFieldValue(st, "PublishedAt")
	assert.NoError(t, err, "get PublishedAt value")
	assert.Equal(t, time.Now().Year(), publishedAt.(time.Time).Year())

	anyTarget, err := StructFieldValue(st, "AnyTarget")
	assert.NoError(t, err, "get AnyTarget value")
	assert.Equal(t, "hi", anyTarget)

	_, err = StructFieldValue(st, "UndefinedField")
	assert.Error(t, err, "get UndefinedField value must be error")
}

func TestStructCast(t *testing.T) {
	type testStruct struct {
		Name        string
		Value       int64
		Count       customInt
		CreatedAt   time.Time
		UpdatedAt   time.Time
		PublishedAt time.Time
		AnyTarget   any
	}
	res, err := Struct[testStruct](map[string]any{
		"Name":        "test",
		"Value":       "1900",
		"Count":       112.2,
		"CreatedAt":   "2020/10/10",
		"UpdatedAt":   time.Now().Unix(),
		"PublishedAt": time.Now(),
		"AnyTarget":   "hi",
	})
	assert.NoError(t, err)
	assert.Equal(t, "test", res.Name)
	assert.Equal(t, int64(1900), res.Value)
	assert.Equal(t, customInt(112), res.Count)
	assert.Equal(t, 2020, res.CreatedAt.Year())
	assert.Equal(t, time.Now().Year(), res.UpdatedAt.Year())
	assert.Equal(t, time.Now().Year(), res.PublishedAt.Year())
	assert.Equal(t, "hi", res.AnyTarget)
}

func BenchmarkGetSetFieldValue(b *testing.B) {
	st := &struct{ Name string }{}

	b.Run("set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := SetStructFieldValue(st, "Name", "value")
			if err != nil {
				b.Error("set field error", err.Error())
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := StructFieldValue(st, "Name")
			if err != nil {
				b.Error("get field error", err.Error())
			}
		}
	})
}
