package infrastructure

import "github.com/labstack/echo/v4"

type CustomBinder struct {
	binder echo.Binder
}

func (b *CustomBinder) Bind(i interface{}, c echo.Context) error {
	if err := b.binder.Bind(i, c); err != nil {
		return err
	}

	if c.Echo().Validator != nil {
		return c.Echo().Validator.Validate(i)
	}

	return nil
}

func NewCustomBinder() *CustomBinder {
	return &CustomBinder{
		binder: new(echo.DefaultBinder),
	}
}
