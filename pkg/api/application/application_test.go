package application

import (
	"testing"
)

func TestString(t *testing.T) {
	app := new(App)
	t.Error(app.String())
}
