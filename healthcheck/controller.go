package healthcheck

import (
	"net/http"

	"github.com/go-chi/render"
)

func (c *controller) Healthcheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.metrics.EmitDist(ctx, "healthcheck", 1.0, nil)
	render.PlainText(w, r, "OK")
}
