package templates

import (
	"context"
	"photos"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func (s *Service) Dashboard(ctx context.Context, user *photos.User) g.Node {

	return s.layoutCore(
		ctx,
		Div(
			Class("container"),
			Div(
				Class("row"),
				Div(
					Class("col-12 mt-2"),
					Div(
						Class("d-flex justify-content-between align-items-center"),

						H1(
							Class("mb-0"),
							g.Text("Dashboard"),
						),
						A(
							Href(s.buildRoute("upload")),
							Class("mb-0 btn btn-primary m-0"),
							g.Text("Upload Image"),
						),
					),

					Hr(),
					P(
						g.Textf("Welcome, %s", user.Name),
					),
				),
			),
		),
	)
}
