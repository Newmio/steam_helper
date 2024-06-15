package steam_helper

import (
	"fmt"
	"runtime"

	"github.com/tebeka/selenium"
)

func Trace(err error, any ...interface{}) error {
	if err == nil {
		return nil
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}

	var str string

	for _, value := range any {

		switch v := value.(type) {

		case selenium.WebDriver:
			html, err := v.PageSource()
			if err != nil {
				continue
			}

			str += html

		case selenium.WebElement:
			html, err := v.GetAttribute("outerHTML")
			if err != nil {
				continue
			}

			str += html

		default:
			str += fmt.Sprint(v)
		}
	}

	return fmt.Errorf("%s%s%s%s(*_*) %s:%d (*_*)", err.Error(), "\n", str, "\n", file, line)
}