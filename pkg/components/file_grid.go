package components

import (
	"os"
	"path/filepath"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type FileGrid struct {
	app.Compo

	Index []os.FileInfo

	SelectedPath    string
	SetSelectedPath func(string)

	CurrentPath    string
	SetCurrentPath func(string)

	Standalone bool
}

func (c *FileGrid) Render() app.UI {
	return app.Div().
		Class(func() string {
			classes := "pf-l-grid pf-m-gutter"
			if c.Standalone {
				classes += " pf-m-all-4-col-on-md pf-m-all-3-col-on-xl pf-u-py-md"
			} else {
				classes += " pf-m-all-4-col-on-sm pf-m-all-4-col-on-md pf-m-all-3-col-on-lg pf-m-all-2-col-on-xl"
			}

			return classes
		}()).
		Body(
			app.Range(c.Index).Slice(func(i int) app.UI {
				return app.Div().
					Class("pf-l-grid__item pf-u-text-align-center").
					Body(
						app.Div().
							Class(
								func() string {
									classes := "pf-c-card pf-m-plain pf-m-hoverable pf-m-selectable"
									if c.SelectedPath == filepath.Join(c.CurrentPath, c.Index[i].Name()) {
										classes += " pf-m-selected"
									}

									return classes
								}()).
							OnClick(func(ctx app.Context, e app.Event) {
								newSelectedPath := filepath.Join(c.CurrentPath, c.Index[i].Name())
								if c.SelectedPath == newSelectedPath {
									newSelectedPath = ""
								}

								c.SetSelectedPath(newSelectedPath)
							}).
							OnDblClick(func(ctx app.Context, e app.Event) {
								if c.Index[i].IsDir() {
									e.PreventDefault()

									c.SetCurrentPath(filepath.Join(c.CurrentPath, c.Index[i].Name()))

									c.SetSelectedPath("")
								}
							}).
							Aria("role", "button").
							TabIndex(0).
							Body(
								app.Div().
									Class("pf-c-card__body").
									Body(
										app.I().
											Class(func() string {
												classes := "fas pf-u-font-size-3xl"
												if c.Index[i].IsDir() {
													classes += " fa-folder"
												} else {
													classes += " fa-file-alt"
												}

												return classes
											}()).
											Aria("hidden", true),
									),
								app.Div().
									Class("pf-c-card__footer").
									Body(
										app.Text(c.Index[i].Name()),
									),
							),
					)
			}),
		)
}
