package main

import (
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

//get cat-file -p

type object struct {
	mod string //commit,tree,blob,tag
	hash string //hash string
	children []*object //child
	regCts [][]string //regexp rs for text
	length int //text length
	level int //tree level
	name string //display name when draw
	top float64 //top px when draw
	left float64 //left px when draw
}

/*
func findChildren(obj *object, objs []*object){

	switch obj.mod {

	case "commit":
		for _,o := range objs{
			if obj.regCts[0][1] == o.hash {
				obj.children = append(obj.children,o)
				break
			}
		}
	case "tree":


	}

}
*/

func setChildrenLevel(obj *object){

	if len(obj.children)==0 {
		return
	}

	i := obj.level-1
	for _,chi := range obj.children {
		chi.level = i
		setChildrenLevel(chi)
	}

}

func fixMap(objs map[string]*object){

	for _,obj := range objs{

		switch obj.mod {

		case "commit":
			obj.name = "[commit]"
			obj.children = append(obj.children, objs[obj.regCts[0][1]])
			objs[obj.regCts[0][1]].name = "[root]"
		case "tree":
			for _,al := range obj.regCts {
				obj.children = append(obj.children,objs[al[2]])
				objs[al[2]].name = al[3]
			}
		case "tag":

		}


	}

}

func readObjects(path string) (map[string]*object,[]string) {

	//
	objs := make(map[string]*object)
	hashs := []string{}

	//Check git path exsit

	_, err := os.Stat(path+string(os.PathSeparator)+".git")
	if err != nil && os.IsNotExist(err) {
		return nil,nil
	}

	//Change dir for git commands
	os.Chdir(path)

	//Cat all obejcts
	rs , err := exec.Command("git","cat-file","--batch-check","--batch-all-objects").Output()
	if err !=  nil{
		return nil,nil
	}
	cts := string(rs)

	//Read each line and insert to array
	regLines,_:= regexp.Compile(`([0-9a-z]+)\s(commit|blob|tree|tag)\s([\d]+)\n`)
	regCommit,_:= regexp.Compile(`tree\s([0-9a-z]+)\n`)
	regTree,_ := regexp.Compile(`[\d]+\s(commit|blob|tree|tag)\s([0-9a-z]+)[\s]+([\S]+)\n`)
	//regBlob,_ := regexp.Compile(``)
	//regTag,_ := regexp.Compile(``)

	lines := regLines.FindAllStringSubmatch(cts,-1)

	for _,line := range lines{

		hashs = append(hashs,line[1])
		ot,_ := exec.Command("git","cat-file","-p",line[1]).Output()
		content := string(ot)

		switch line[2] {

		case "commit":
			i,_ := strconv.Atoi(line[3])
			s := regCommit.FindAllStringSubmatch(content,-1)
			obj := &object{
				mod:"commit",
				hash:line[1],
				regCts:s,
				length:i,
			}
			objs[line[1]]=obj

		case "tree":
			i,_ := strconv.Atoi(line[3])
			s := regTree.FindAllStringSubmatch(content,-1)
			obj := &object{
				mod:"tree",
				hash:line[1],
				regCts:s,
				length:i,
			}
			objs[line[1]]=obj
		case "blob":
			i,_ := strconv.Atoi(line[3])
			obj := &object{
				mod:"blob",
				hash:line[1],
				length:i,
			}
			objs[line[1]]=obj

		case "tag":

		}


	}

	return objs,hashs

}


func parseProject(path string) (map[string]*object,[]string){

	mp,hashs := readObjects(path)
	if mp==nil{
		return nil,nil
	}

	fixMap(mp)
	for _,m:=range mp{
		setChildrenLevel(m)
	}

	return mp,hashs

}


