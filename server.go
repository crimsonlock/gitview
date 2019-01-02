package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/", htmlProject)
}

func htmlProject(response http.ResponseWriter, request *http.Request) {

	//return content head
	cts := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<style type="text/css">
div{padding:0}	
.tag{height:50px;width:150px;border:1px solid black;position:absolute}
.pointer{height:50px;width:150px;border:1px solid black;position:absolute}
.commit{height:50px;width:150px;border:1px solid black;position:absolute}
.tree{height:50px;width:150px;border:1px solid black;position:absolute}
.blob{height:50px;width:150px;border:1px solid black;position:absolute}	
.line{background-color:black;height:1px;position:absolute}	
p{margin:5px;padding:0px;text-align:center;font-size:14px;}	
</style>
</head>
<body>`

	//max level for display
	stack := [10]int{}

	//get real path
	request.ParseForm()
	path, ok := request.Form["path"]
	pp := ""
	if ok {
		pp = path[0]
	}

	//read objects
	mp, hashs := parseProject(string(os.PathSeparator) + pp)
	if mp == nil {
		fmt.Fprintf(response, "%s<p>Bad project path : %c%s</p><p>Usage sample: http://127.0.0.1:8802/?path=/Users/lishaopeng/go/src/github.com/crimsonlock/gitview</p></body></html>", cts, os.PathSeparator, pp)
		return
	}

	//show objects rects
	for _, hs := range hashs {
		v := mp[hs]
		v.top = float64(100 - 200*v.level)
		v.left = float64(100 + 200*stack[-1*v.level])
		stack[-1*v.level]++
		cts = fmt.Sprintf("%s<div class='%s' style='top:%dpx;left:%dpx'><p>%s<br>%s</p></div>\n", cts, v.mod, int(v.top), int(v.left), v.hash[0:6], v.name)
	}

	//draw lines between rects
	for _, hs := range hashs {
		v := mp[hs]
		for _, chd := range v.children {
			cts = cts + drawLine(v, chd)
		}
	}

	fmt.Fprintf(response, "%s</body></html>", cts)

}

//draw lines between o1 and o2
func drawLine(o1 *object, o2 *object) string {

	//line length
	length := math.Sqrt(float64(22500 + (o1.left-o2.left)*(o1.left-o2.left)))

	//line top
	top := int((o1.top+o2.top)/2) + 25

	//line left
	left := int((o1.left+o2.left-length)/2) + 75

	//line rotation
	rot := 0.0
	if o1.left == o2.left {
		rot = 90.0
	} else {
		rot = 180.0 * math.Atan(-150/(o1.left-o2.left)) / math.Pi
	}

	return fmt.Sprintf("<div class='line' style='top:%dpx;left:%dpx;width:%dpx;-webkit-transform:rotate(%fdeg);'></div>\n", top, left, int(length), rot)
}

func main() {

	http.ListenAndServe(":8802", nil)

}
