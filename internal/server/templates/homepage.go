package templates

import (
	"context"
	"photos"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func (s *Service) Homepage(ctx context.Context, user *photos.User) g.Node {

	return s.layoutCore(
		ctx,
		Div(
			Class("banner"),
			Div(
				Class("overlay"),
				Div(
					Class("mt-5"),
					H1(g.Text("Yet Another Photo Gallery")),
					Hr(),
					P(
						Class("text-center"),
						g.Text("Upload and Store Your Photos, backed by AWS S3"),
					),
					P(
						Class("text-center"),
						Button(
							Class("btn btn-primary"),
							g.Text("Get Started"),
						),
					),
				),
			),
		),
	)
}
