package components

import (
	"os"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type PathPickerToolbar struct {
	app.Compo

	Index        []os.FileInfo
	RefreshIndex func()

	PathComponents []string

	CurrentPath    string
	SetCurrentPath func(string)

	SelectedPath    string
	SetSelectedPath func(string)

	OpenCreateDirectoryModal func()
}

func (c *PathPickerToolbar) Render() app.UI {
	return app.Div().
		Class("pf-c-toolbar pf-u-pt-0").
		Body(
			app.Div().
				Class("pf-c-toolbar__content").
				Body(
					app.Div().
						Class("pf-c-toolbar__content-section").
						Body(
							app.Div().
								Class("pf-c-toolbar__item pf-m-overflow-menu").
								Body(
									app.Div().
										Class("pf-c-overflow-menu").
										Body(
											app.Div().
												Class("pf-c-overflow-menu__content").
												Body(
													app.Div().
														Class("pf-c-overflow-menu__group pf-m-button-group").
														Body(
															app.Div().
																Class("pf-c-overflow-menu__item").
																Body(
																	&Breadcrumbs{
																		PathComponents: c.PathComponents,

																		CurrentPath:    c.CurrentPath,
																		SetCurrentPath: c.SetCurrentPath,

																		SelectedPath: c.SelectedPath,
																		SetSelectedPath: func(s string) {
																			c.SetSelectedPath(s)
																		},
																	},
																),
														),
												),
										),
								),
							app.Div().
								Class("pf-c-toolbar__item pf-m-pagination").
								Body(
									app.Div().
										Class("pf-c-pagination pf-m-compact").
										Body(
											app.Div().
												Class("pf-c-pagination pf-m-compact pf-m-compact").
												Body(
													app.Div().
														Class("pf-c-overflow-menu").
														Body(
															app.Div().
																Class("pf-c-overflow-menu__content").
																Body(
																	app.Div().
																		Class("pf-c-overflow-menu__group pf-m-button-group").
																		Body(
																			app.Div().
																				Class("pf-c-overflow-menu__item").
																				Body(
																					app.Button().
																						Type("button").
																						Aria("label", "Create directory").
																						Title("Create directory").
																						Class("pf-c-button pf-m-plain").
																						OnClick(func(ctx app.Context, e app.Event) {
																							c.OpenCreateDirectoryModal()
																						}).
																						Body(
																							app.I().
																								Class("fas fa-folder-plus").
																								Aria("hidden", true),
																						),
																				),
																		),
																),
															app.Div().
																Class("pf-c-divider pf-m-vertical pf-m-inset-md").
																Aria("role", "separator"),
															app.Div().
																Class("pf-c-overflow-menu__group pf-m-button-group").
																Body(
																	app.Div().
																		Class("pf-c-overflow-menu__item").
																		Body(
																			app.Button().
																				Type("button").
																				Aria("label", "Refresh").
																				Title("Refresh").
																				Class("pf-c-button pf-m-plain").
																				OnClick(func(ctx app.Context, e app.Event) {
																					c.RefreshIndex()
																				}).
																				Body(
																					app.I().
																						Class("fas fas fa-sync").
																						Aria("hidden", true),
																				),
																		),
																),
														),
												),
										),
								),
						),
				),
		)
}
