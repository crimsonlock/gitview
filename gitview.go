package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//get cat-file -p

type object struct {
	mod string //commit,tree,blob,tag
	hash string
	children []*object
	length int
	level int
}

func readObjects(path string) []*object {

	//Check git path exsit
	_, err := os.Stat(path+"/.git")
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	//Change dir for git commands
	os.Chdir(path)

	//Cat all obejcts
	rs , err := exec.Command("git","cat-file","--batch-check","--batch-all-objects").Output()
	if err !=  nil{
		return nil
	}
	cts := string(rs)

	//Read each line and insert to array
	lines := strings.Split(cts,"\n")
	for _,line := range lines{
		wds := strings.Split(line," ")
		info,_ := exec.Command("git","cat-file","-p",wds[0]).Output()
		fmt.Printf("%s\n",info)
	}

	return nil

}

func parseProject(path string) {




	//fis,_:=ioutil.ReadDir(path+"/.git/objects")

	/*
	for _,fi := range fis{

		if fi.Name()!="pack" && fi.Name()!="info"{

			ffis,_ := ioutil.ReadDir(path+"/.git/objects/"+fi.Name())
			for _,ffi := range ffis{
				//cts,_ := ioutil.ReadFile(path+"/.git/objects/"+fi.Name()+"/"+ffi.Name())
				//fmt.Println(string(cts))
				//exec.CommandContext()
				rs ,err := exec.Command("git","cat-file","-p",fi.Name()+ffi.Name()).Output()
				if err !=nil{
					fmt.Println(err.Error())
				}else {
					fmt.Printf("%s\n=========\n",rs)
				}

			}


		}


	}
	*/

}


func main(){

}
