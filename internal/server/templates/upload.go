package templates

import (
	"context"
	"photos"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

const petiteVueScript = `
import {createApp} from 'https://unpkg.com/petite-vue?module'
createApp({
	files: [],
	progress: 0,
	onFileChange(e) {
		console.log(e)
		this.files = e.target.files
	},
	submit(e) {
		e.preventDefault()
		console.log(this.files)


		for (let i = 0; i<this.files.length; i++) {
			const file = this.files[i]
			if (file.type.indexOf('image/') !== 0) {
				console.log('Not an image file.')
				continue
			}
			fetch('https://upload.photos.onetwentyseven.dev/' + file.name, {
				credentials: 'include',
				method: 'PUT',
				headers: {
					'Content-Type': file.type
				},
				body: file,
			}).then(r => {
				console.log(r)
			}).catch(e => {
				console.log('Error Processing request', e.message)
			})
		}

		// htmx.
	}
}).mount("#upload-image")
`

func (s *Service) Upload(ctx context.Context, user *photos.User) g.Node {

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
							g.Text("Upload"),
						),
					),

					Hr(),
				),
			),
			Div(
				Class("row border border-secondary"),
				StyleAttr("height: 15vh;"),
				Div(
					ID("upload-image"),
					Class("col-12 d-flex justify-content-center align-items-center"),

					FormEl(
						EncType("multipart/form-data"),
						Input(
							Type("file"),
							Multiple(),
							Name("image"),
							Class("form-control-file"),
							Accept("image/*"),
							g.Attr("@change", "onFileChange"),
						),
						Button(
							Class("mb-0 btn btn-primary m-0"),
							Type("button"),
							g.Attr("@click", "submit"),
							g.Text("Upload Image"),
						),
					),
				),
			),
			Script(
				Type("module"),
				g.Raw(petiteVueScript),
			),
		),
	)
}
