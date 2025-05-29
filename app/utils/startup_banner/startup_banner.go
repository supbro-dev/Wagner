package startup_banner

import "github.com/common-nighthawk/go-figure"

func Run() {
	// 瓦格纳，全军出击！
	myFigure := figure.NewFigure("Wagner, all units engage!", "", true)
	myFigure.Print()
}
