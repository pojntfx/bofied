package components

import (
	"os"
	"path/filepath"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type FileGrid struct {
	app.Compo

	Index []os.FileInfo

	SelectedPath    string
	SetSelectedPath func(string)

	CurrentPath    string
	SetCurrentPath func(string)

	Standalone bool

	hasInitiatedClick bool
}

func (c *FileGrid) Render() app.UI {
	return app.Div().
		Class(func() string {
			classes := "pf-v6-l-grid pf-m-gutter"
			if c.Standalone {
				classes += " pf-m-all-4-col-on-md pf-m-all-3-col-on-xl pf-v6-u-py-md"
			} else {
				classes += " pf-m-all-4-col-on-sm pf-m-all-4-col-on-md pf-m-all-3-col-on-lg pf-m-all-3-col-on-xl"
			}

			return classes
		}()).
		Body(
			app.Range(c.Index).Slice(func(i int) app.UI {
				selectCard := func() {
					newSelectedPath := filepath.Join(c.CurrentPath, c.Index[i].Name())
					if c.SelectedPath == newSelectedPath {
						// Handle double click
						if c.hasInitiatedClick && c.Index[i].IsDir() {
							c.SetCurrentPath(filepath.Join(c.CurrentPath, c.Index[i].Name()))

							c.SetSelectedPath("")

							return
						}

						newSelectedPath = ""
					}

					c.SetSelectedPath(newSelectedPath)

					// Prepare for double click
					c.hasInitiatedClick = true
					time.AfterFunc(time.Second, func() {
						c.hasInitiatedClick = false
					})
				}

				return app.Div().
					Class("pf-v6-l-grid__item pf-v6-u-text-align-center").
					Body(
						app.Div().
							Class(
								func() string {
									classes := "pf-v6-c-card pf-m-plain pf-m-selectable"
									if c.SelectedPath == filepath.Join(c.CurrentPath, c.Index[i].Name()) {
										classes += " pf-m-selected"
									}

									return classes
								}()).
							On("keyup", func(ctx app.Context, e app.Event) {
								if e.Get("key").String() == "Enter" || e.Get("key").String() == " " {
									selectCard()
								}
							}).
							OnClick(func(ctx app.Context, e app.Event) {
								selectCard()
							}).
							Aria("role", "button").
							TabIndex(0).
							Body(
								app.Div().
									Class("pf-v6-c-card__body").
									Body(
										app.I().
											Class(func() string {
												classes := "fas pf-v6-u-font-size-3xl"
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
									Class("pf-v6-c-card__footer").
									Body(
										app.Text(c.Index[i].Name()),
									),
							),
					)
			}),
		)
}
