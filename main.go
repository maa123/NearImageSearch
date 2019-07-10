package main

import(
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var distanceValue uint8 = 4
var hdb [][]uint64

func main() {
	hdb = LoadDB("hash.db")
	imgv := new(walk.ImageView)
	var sliderLabel *walk.TextLabel
	var textLabel *walk.TextEdit
	var slider *walk.Slider
	var preSearchHash uint64
	var searchRes string
	mw := new(walk.MainWindow)
	MainWindow{
		AssignTo: &mw,
		Title:   "Img Search",
		MinSize: Size{400, 640},
		Size: Size{400, 640},
		Layout:  VBox{MarginsZero: true},
		OnDropFiles: func(files []string) {
			im, err := walk.NewImageFromFile(files[0])
			if err == nil {
				imgv.SetImage(im)
				imgv.SetMode(walk.ImageViewModeShrink)
				searchRes, preSearchHash = SearchFromPath(files[0], distanceValue, hdb)
				textLabel.SetText(searchRes)
			}
		},
		Children: []Widget{
			TextLabel{
				AssignTo: &sliderLabel,
				Text: "検索範囲:4",
			},
			Slider{
				AssignTo:    &slider,
				MinValue:    1,
				MaxValue:    20,
				Value:       4,
				OnValueChanged: func() {
					distanceValue = uint8(slider.Value())
					sliderLabel.SetText(fmt.Sprintf("検索範囲:%d", distanceValue))
					if preSearchHash != 0 {
						textLabel.SetText(SearchFromHash(preSearchHash, distanceValue, hdb))
					}
				},
			},
			ImageView{
				AssignTo: &imgv,
				MinSize: Size{380, 300},
				MaxSize: Size{380, 300},
			},
			TextEdit{
				AssignTo: &textLabel,
				Text:     "\r\n\r\n\r\n\r\n\r\n",
				MinSize: Size{380, 150},
				MaxSize: Size{380, 150},
				ReadOnly: true,
				VScroll:  true,
				Row: 8,
			},
		},
	}.Run()
}