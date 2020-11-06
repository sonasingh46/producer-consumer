// package widget contains widget specific definitions
package widget

import "time"

// Widget represents a widget object struct.
type Widget struct {
	// CreatedAt has the timestamp of its creation
	CreatedAt time.Time
	/* Similar other relevant fields can go on here */
}

// NewWidget returns a new widget.
func NewWidget() Widget  {
	return Widget{
		CreatedAt:time.Now(),
	}
}

/* Other widget specific methods can go on here */

