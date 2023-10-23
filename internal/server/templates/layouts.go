package templates

import (
	"context"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func (s *Service) layoutCore(ctx context.Context, node g.Node) g.Node {

	return Doctype(
		HTML(
			Lang("en"),
			s.gtop(ctx),
			Body(
				s.gnavbar(ctx),
				node,
				s.gbottom(),
			),
		),
	)

}
