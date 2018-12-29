package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

//get cat-file -p

type node struct {
	mod string //commit,tree,blob,tag
	path string
	children []*node
	name string
	level uint8
}

func parseProject(path string) {


	fis,_:=ioutil.ReadDir(path+"/.git/objects")

	for _,fi := range fis{

		if fi.Name()!="pack" && fi.Name()!="info"{

			ffis,_ := ioutil.ReadDir(path+"/.git/objects/"+fi.Name())
			for _,ffi := range ffis{
				//cts,_ := ioutil.ReadFile(path+"/.git/objects/"+fi.Name()+"/"+ffi.Name())
				//fmt.Println(string(cts))

				rs ,err := exec.Command("/usr/bin/git","cat-file","-p",fi.Name()+ffi.Name()).Output()
				if err !=nil{
					fmt.Println(err.Error())
				}else {

				fmt.Printf("%s",rs)
				}

			}


		}


	}

}


func main(){

}
