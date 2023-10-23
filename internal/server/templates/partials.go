package templates

import (
	"context"
	"photos/internal"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func (s *Service) gtop(_ context.Context) g.Node {
	return Head(
		Meta(Charset("utf-8")),
		TitleEl(g.Text("Yet Another Photo Gallery")),
		Link(
			Href("https://cdn.jsdelivr.net/npm/bootstrap@5.3.1/dist/css/bootstrap.min.css"),
			Rel("stylesheet"),
			g.Attr("integrity", "sha384-4bw+/aepP/YC94hEpVNVgiZdgIC5+VKNBQNGCHeKRQN+PtmoHDEXuppvnDJzQIu9"),
			g.Attr("crossorigin", "anonymous"),
		),
		Link(
			Href("https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.2/css/all.min.css"),
			Rel("stylesheet"),
			g.Attr("integrity", "sha512-z3gLpd7yknf1YoNbCzqRKc4qyor8gaKU1qmn+CShxbuBusANI9QpRohGBreCFkKxLhei6S9CQXFEbbKuqLg0DA=="),
			g.Attr("crossorigin", "anonymous"),
			g.Attr("referrerpolicy", "no-referrer"),
		),
		Link(
			Href("/static/css/stylesheet.css"),
			Rel("stylesheet"),
		),
	)
}

func (s *Service) gnavbar(ctx context.Context) g.Node {

	return Nav(
		Class("navbar navbar-expand-lg bg-dark"),
		DataAttr("bs-theme", "dark"),
		Div(
			Class("container"),
			A(
				Class("navbar-brand"),
				Href(s.buildRoute("home")),
				g.Text("YAPG"),
			),
			Button(
				Class("navbar-toggler"),
				Type("button"),
				DataAttr("bs-toggle", "collapse"),
				DataAttr("bs-target", "#navbarSupportedContent"),
				Aria("controls", "navbarSupportedContent"),
				Aria("expended", "false"),
				Aria("label", "Toggle Navigation"),
				Span(
					Class("navbar-toggler-icon"),
				),
			),
			Div(
				Class("collapse navbar-collapse"),
				ID("navbarSupportedContent"),
				Ul(
					Class("navbar-nav me-auto mb-2 mb-lg-0"),
					Li(
						Class("nav-item"),
						A(
							Class("nav-link active"),
							Aria("current", "page"),
							Href(s.buildRoute("home")),
							g.Text("Home"),
						),
					),
				),
				Ul(
					Class("navbar-nav mb-2 mb-lg-0"),
					s.renderNavbarUserMenu(ctx),
				),
			),
		),
	)
}

func (s *Service) renderNavbarUserMenu(ctx context.Context) g.Node {

	user := internal.UserFromContext(ctx)

	if user != nil {
		return Li(
			Class("nav-item dropdown"),
			A(
				Class("nav-link dropdown-toggle"), Href(s.buildRoute("login")), Role("button"), DataAttr("bs-toggle", "dropdown"),
				g.Textf("Hello %s", user.Name),
			),
			Ul(
				Class("dropdown-menu"),
				Li(A(Class("dropdown-item"), Href(s.buildRoute("dashboard")), g.Text("Dashboard"))),
				Li(Hr(Class("dropdown-divider"))),
				Li(A(Class("dropdown-item"), Href(s.buildRoute("logout")), g.Text("Logout"))),
			),
		)
	}

	return Li(
		Class("nav-item"),
		A(
			Class("nav-link"),
			Href(s.buildRoute("login")),
			g.Text("Login"),
		),
	)

}

func (s *Service) ghtmxDebug() g.Node {
	return Script(
		g.Raw(`
htmx.logger = function (elt, event, data) {
	if (!console) return
	console.log(event,elt.nodeName, elt, data);
}
		`),
	)
}

func (s *Service) gbottom() g.Node {
	return g.Group([]g.Node{
		Script(
			Src("https://unpkg.com/htmx.org@1.9.4"),
			g.Attr("integrity", "sha384-zUfuhFKKZCbHTY6aRR46gxiqszMk5tcHjsVFxnUo8VMus4kHGVdIYVbOYYNlKmHV"),
			g.Attr("crossorigin", "anonymous"),
		),
		Script(
			Src("https://cdn.jsdelivr.net/npm/bootstrap@5.3.1/dist/js/bootstrap.bundle.min.js"),
			g.Attr("integrity", "sha384-HwwvtgBNo3bZJJLYd8oVXjrBZt8cqVSpeBNS5n7C8IVInixGAoxmnlMuBnhbgrkm"),
			g.Attr("crossorigin", "anonymous"),
		),
		s.ghtmxDebug(),
	})
}
