package main

import (
	"fmt"
	"syscall/js"

	"github.com/maxence-charriere/go-app/v6/pkg/app"
	"github.com/suyashkumar/dicom"
)

type hello struct {
	app.Compo

	Name string
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.Div().
			Class("menu-button").
			OnClick(h.OnMenuClick).
			Body(
				app.Text("☰"),
			),
		app.Main().
			Class("hello").
			Body(
				app.H1().
					Class("hello-title").
					Body(
						app.Text("test, "),
						app.If(h.Name != "",
							app.Text(h.Name),
						).Else(
							app.Text("dicom"),
						),
					),
				app.Input().
					Class("hello-input").
					Value(h.Name).
					Placeholder("input aaaa?").
					AutoFocus(true).
					OnChange(h.OnInputChange),
				app.Input().Type("file").OnChange(h.OnFileChange),
			),
	)
}

func (h *hello) OnMenuClick(src app.Value, e app.Event) {
	app.NewContextMenu(
		app.MenuItem().
			Label("Reload").
			Keys("cmdorctrl+r").
			OnClick(func(src app.Value, e app.Event) {
				app.Reload()
			}),
		app.MenuItem().Separator(),
		app.MenuItem().
			Label("City demo").
			OnClick(func(src app.Value, e app.Event) {
				app.Navigate("/city")
			}),
		app.MenuItem().Separator(),
		app.MenuItem().
			Icon(icon).
			Label("Go to repository").
			OnClick(func(src app.Value, e app.Event) {
				app.Navigate("https://github.com/maxence-charriere/go-app")
			}),
		app.MenuItem().
			Icon(icon).
			Label("Demo sources").
			OnClick(func(src app.Value, e app.Event) {
				app.Navigate("https://github.com/maxence-charriere/go-app-demo/tree/v6/demo")
			}),
	)
}

func (h *hello) OnInputChange(src app.Value, e app.Event) {

	h.Name = src.Get("value").String()
	h.Update()
}

func OnLoadData(this js.Value, args []js.Value) interface{} {
	array := js.Global().Get("Uint8Array").New(args[0].Get("target").Get("result"))
	fmt.Println("size: ", array.Get("length"))
	inBuf := make([]byte, array.Get("length").Int())
	js.CopyBytesToGo(inBuf, array)
	fmt.Println(inBuf)
	// do something with inBuf

	// opts := dicom.ReadOptions{DropPixelData: true}
	// element, err := dicom.ReadDataSetInBytes(inBuf, opts)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(element)

	p, err := dicom.NewParserFromBytes(inBuf, nil)
	opts := dicom.ParseOptions{DropPixelData: true}
	if err != nil {
		fmt.Println(err)
	}
	element, err := p.Parse(opts)
	fmt.Println(element)
	return nil
}
func (h *hello) OnFileChange(src app.Value, e app.Event) {

	reader := js.Global().Get("FileReader").New()

	reader.Set("onload", js.FuncOf(OnLoadData))

	reader.Call("readAsArrayBuffer", app.JsVal(src.Get("files").Get("0")))

	/*
		let reader = new FileReader();
		reader.onload = (ev) => {
			bytes = new Uint8Array(ev.target.result);
			loadImage(bytes);
			let blob = new Blob([bytes], {'type': imageType});
			document.getElementById("sourceImg").src = URL.createObjectURL(blob);
		};
		imageType = this.files[0].type;
		reader.readAsArrayBuffer(this.files[0]);
	*/

	// js.ValueOf(app.ValueOf(src))
	// app.ValueOf(src)

	// src.Get("value")
	// src.JSValue()
	// js.Global().Get("console").Call("log", "time takena1:", src.JSValue())
	fmt.Println("===== ", e.JSValue().Type(), src.Get("files").Get("0").Get("size"))

}
