package doc

import (
	"ele/constants"
	"fmt"

	c "github.com/logrusorgru/aurora"
)

// Title ...
func Title(title string) c.Value {

	return c.Yellow(title)
}

// Titlef ...
func Titlef(format string, a ...interface{}) c.Value {

	text := fmt.Sprintf(format, a...)
	return c.Yellow(text)
}

// Info ...
func Info(info string) c.Value {

	return c.White(info)
}

// Infof ...
func Infof(format string, a ...interface{}) c.Value {

	text := fmt.Sprintf(format, a...)
	return c.White(text)
}

// Success ...
func Success(info string) c.Value {

	return c.Green(info)
}

// Successf ...
func Successf(format string, a ...interface{}) c.Value {

	text := fmt.Sprintf(format, a...)
	return c.Green(text)
}

// Data ...
func Data(data interface{}) c.Value {

	return c.Cyan(data)
}

// Dataf ...
func Dataf(format string, a ...interface{}) c.Value {

	text := fmt.Sprintf(format, a...)
	return c.Cyan(text)
}

// Hint ...
func Hint(hint string) c.Value {

	return c.White(c.Italic(hint))
}

// Hintf ...
func Hintf(format string, a ...interface{}) c.Value {

	text := fmt.Sprintf(format, a...)
	return c.White(c.Italic(text))
}

// Error ...
func Error(errorvalue string) c.Value {
	return c.Red(errorvalue)
}

// Errorf ...
func Errorf(format string, a ...interface{}) c.Value {
	text := fmt.Sprintf(format, a...)
	return c.Red(text)
}

// Option ...
func Option(value string) c.Value {
	format := fmt.Sprintf("[%s]", value)
	return c.White(format)
}

// Tag ...
func Tag(start string, end string) c.Value {
	if start == end {
		if len(start) == 0 {
			return c.BgCyan(c.Black(" ALL "))
		}
		if start == constants.REPEATABLE {
			return c.BgCyan(c.Black(" REPEATABLE "))
		}
		return c.BgCyan(c.Black(" " + start + " "))
	}

	return c.BgCyan(c.Black(" " + start + ":" + end + " "))
}

// Env ...
func Env(name string) c.Value {
	if name == "" {
		return c.White("")
	}
	return c.BgBlue(c.Black(" " + name + " "))
}

// YesNo ...
func YesNo(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}
