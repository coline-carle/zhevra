package toc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Decoder struct {
	r *bufio.Reader
}

var (
	errHashtagNotFound = errors.New("Hashtag not found")
	errNewLine         = errors.New("Newline unexpected")
	errColonNotFound   = errors.New("Colon not found")
)

// NewDecoder dreate decoder struct
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: bufio.NewReader(r),
	}
}

func structFields(v interface{}) map[string]reflect.Value {
	d := reflect.ValueOf(v)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
		if !d.IsValid() {
			return nil
		}
	}
	if d.Kind() == reflect.Struct {
		tags := make(map[string]reflect.Value, 0)
		for i := 0; i < d.NumField(); i++ {
			if field := d.Type().Field(i); field.Tag.Get("toc") == "" {
				tags[field.Name] = d.Field(i)
			} else {
				tags[field.Tag.Get("toc")] = d.Field(i)
			}
		}
		return tags
	}
	return nil
}

// Decode toc file
func (d *Decoder) Decode(v interface{}) error {
	fields := structFields(v)

outer:
	for {
		err := d.scanHashtag()
		switch err {
		case errHashtagNotFound:
			d.skipUntil('\n')
			continue outer
		case errNewLine:
			continue outer
		case nil:
		default:
			return err
		}

		key, err := d.scanKey()
		switch err {
		case errColonNotFound:
			continue outer
		case nil:
		default:
			return err
		}

		value, err := d.scanValue()
		if err != nil {
			return err
		}

		// got key and value let's set them
		if field, ok := fields[key]; ok {
			switch field.Kind() {
			case reflect.Int32:
				intValue, err := strconv.ParseInt(value, 10, 32)
				if err != nil {
					return err
				}
				field.SetInt(intValue)
			case reflect.Int, reflect.Int64:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(intValue)
			case reflect.String:
				field.SetString(value)
			default:
				return err
			}
		}
	}
}

func (d *Decoder) scanHashtag() error {
	ch, _, err := d.r.ReadRune()
	if err != nil {
		return err
	}
	if ch == '\n' {
		return errNewLine
	}
	if ch != '#' {
		return errHashtagNotFound
	}

	ch, _, err = d.r.ReadRune()
	if err != nil {
		return err
	}
	if ch == '\n' {
		return errNewLine
	}
	if ch != '#' {
		return errHashtagNotFound
	}
	return nil
}

func (d *Decoder) scanValue() (string, error) {
	var buf bytes.Buffer

	// skip spaces before
	err := d.skipWhitespace()
	if err != nil {
		return "", err
	}

	for {
		ch, _, err := d.r.ReadRune()
		if err != nil {
			return "", err
		}
		if ch == '\n' {
			break
		}
		buf.WriteRune(ch)
	}

	value := strings.TrimSpace(string(buf.Bytes()))

	return value, nil
}

func (d *Decoder) scanKey() (string, error) {
	var previous rune
	var buf bytes.Buffer
	// skip spaces until firs key value
	err := d.skipWhitespace()
	if err != nil {
		return "", err
	}

	for {
		ch, _, err := d.r.ReadRune()
		if err != nil {
			return "", err
		}
		if unicode.IsSpace(ch) {
			break
		}
		buf.WriteRune(ch)
		previous = ch
	}

	if previous != ':' {
		return "", errColonNotFound
	}

	return string(buf.Bytes()[:buf.Len()-1]), nil
}

func (d *Decoder) skipWhitespace() error {
	for {
		ch, _, err := d.r.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(ch) {
			d.r.UnreadRune()
			return nil
		}
	}
}

func (d *Decoder) skipUntil(r rune) error {
	for {
		ch, _, err := d.r.ReadRune()
		if err != nil {
			return err
		}
		if ch == r {
			return nil
		}
	}
}
