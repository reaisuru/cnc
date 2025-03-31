package evaluator

import (
	"cnc/pkg/swash"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

/*
	EvaluatorGo2Swash.go converts from Go types into Swash objects which
	can be interacted and interpolated with via this file.
*/

// Go2Swash will attempt to convert the go reference into swash objects
func (memory *Memory) Go2Swash(header string, reference any) error {
	switch reflect.TypeOf(reference).Kind() {

	case reflect.String:
		if !strings.HasPrefix(header, "$") {
			header = "$" + header
		}

		return memory.allocateVar(&swash.Var{
			Keyword:    &swash.Token{TokenType: swash.INDENT, TokenLiteral: "var"},
			Descriptor: &swash.Token{TokenType: swash.INDENT, TokenLiteral: fmt.Sprint(header)},
			Args: []*swash.Token{
				{TokenType: swash.STRING, TokenLiteral: fmt.Sprint(reference)},
			},
		})

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !strings.HasPrefix(header, "$") {
			header = "$" + header
		}

		return memory.allocateVar(&swash.Var{
			Keyword:    &swash.Token{TokenType: swash.INDENT, TokenLiteral: "var"},
			Descriptor: &swash.Token{TokenType: swash.INDENT, TokenLiteral: fmt.Sprint(header)},
			Args: []*swash.Token{
				{TokenType: swash.NUMBER, TokenLiteral: fmt.Sprint(reference)},
			},
		})

	case reflect.Float64, reflect.Float32:
		if !strings.HasPrefix(header, "$") {
			header = "$" + header
		}

		return memory.allocateVar(&swash.Var{
			Keyword:    &swash.Token{TokenType: swash.INDENT, TokenLiteral: "var"},
			Descriptor: &swash.Token{TokenType: swash.INDENT, TokenLiteral: fmt.Sprint(header)},
			Args: []*swash.Token{
				{TokenType: swash.NUMBER, TokenLiteral: fmt.Sprintf("%.0f", reference)},
			},
		})

	case reflect.Bool:
		if !strings.HasPrefix(header, "$") {
			header = "$" + header
		}

		return memory.allocateVar(&swash.Var{
			Keyword:    &swash.Token{TokenType: swash.INDENT, TokenLiteral: "var"},
			Descriptor: &swash.Token{TokenType: swash.INDENT, TokenLiteral: fmt.Sprint(header)},
			Args: []*swash.Token{
				{TokenType: swash.BOOLEAN, TokenLiteral: fmt.Sprint(reference)},
			},
		})

	case reflect.Struct:
		index := &Object{
			TrueType:   reference,
			Descriptor: &swash.Token{TokenType: swash.INDENT, TokenLiteral: header},
			Values:     NewMemory(memory.wr, memory.rd, memory.packages),
		}

		v, t := reflect.ValueOf(reference), reflect.TypeOf(reference)

		/* registers all the fields */
		for i := 0; i < v.NumField(); i++ {
			if unicode.IsLower(rune(t.Field(i).Name[0])) {
				continue
			}

			name := t.Field(i).Name
			if len(t.Field(i).Tag.Get("swash")) > 0 {
				name = t.Field(i).Tag.Get("swash")
			}

			err := index.Values.Go2Swash(name, v.Field(i).Interface())
			if err != nil {
				return err
			}
		}

		/* registers all the methods */
		for i := 0; i < v.NumMethod(); i++ {
			err := index.Values.Go2Swash(t.Method(i).Name, v.Method(i).Interface())
			if err != nil {
				return err
			}
		}

		memory.memory = append(memory.memory, index)

	case reflect.Func:
		return memory.registerGoFunc(header, reference)

	case reflect.Map:
		index := &Object{
			TrueType:   reference,
			Descriptor: &swash.Token{TokenType: swash.INDENT, TokenLiteral: header},
			Values:     NewMemory(memory.wr, memory.rd, memory.packages),
		}

		v := reflect.ValueOf(reference)
		for _, key := range v.MapKeys() {
			err := index.Values.Go2Swash(key.String(), v.MapIndex(key).Interface())
			if err != nil {
				return err
			}
		}

		memory.memory = append(memory.memory, index)

	case reflect.Pointer:
		if reflect.TypeOf(reference).Elem().Kind() != reflect.Struct {
			return memory.Go2Swash(header, reflect.ValueOf(reference).Elem().Interface())
		}

		index := &Object{
			TrueType:   reference,
			Descriptor: &swash.Token{TokenType: swash.INDENT, TokenLiteral: header},
			Values:     NewMemory(memory.wr, memory.rd, memory.packages),
		}

		value := reflect.ValueOf(reference)
		for i := 0; i < value.NumMethod(); i++ {
			context := reflect.TypeOf(reference).Method(i)
			index.Values.Go2Swash(context.Name, value.Method(i).Interface())
		}

		for i := 0; i < value.Elem().NumField(); i++ {
			name := reflect.TypeOf(reference).Elem().Field(i).Name
			if unicode.IsLower(rune(name[0])) {
				continue
			}

			if len(reflect.TypeOf(reference).Elem().Field(i).Tag.Get("swash")) > 0 {
				name = reflect.TypeOf(reference).Elem().Field(i).Tag.Get("swash")
			}

			index.Values.Go2Swash(name, value.Elem().Field(i).Interface())
		}

		memory.memory = append(memory.memory, index)
	}

	return nil
}
