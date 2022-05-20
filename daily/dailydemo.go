package daily

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendFile() {
	uri:="http://127.0.0.1:8080/file/upload"
	byte,err:=ioutil.ReadFile("text.txt")

	res,err :=http.Post(uri,"text/html",bytes.NewReader(byte))
	if err !=nil {
		fmt.Println("err=",err)
	}
	//http返回的response的body必须close,否则就会有内存泄露
	defer func() {
		res.Body.Close()
		fmt.Println("finish")
	}()
	//读取body
	body,err:=ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Println(" post err=",err)
	}
	fmt.Println(string(body))
}



var jobList []string
var jobTable = make(map[string]int, 100)
type Job struct {
	name string
	deeps []Job
}

func makeJobTable(jobs []Job) {
	for _, job := range jobs {
		jobTable[job.name] = 0
	}
}

//查找整个链路的启动顺序 问题：不一定同级链路都需要启动，但是启动顺序上同级链路启动时都需要启动其他同级以及前置任务
func sortJob(jobs []Job)  {

	for _, job := range jobs {
		if jobTable[job.name] == 0 {
			jobList = append(jobList, job.name)
			jobTable[job.name] = 1
			sortJobChild(job, len(jobList))
		} else {
			for index, str := range jobList {
				if str == job.name {
					sortJobChild(job,index)
				}
			}
		}
	}
}

func sortJobChild(jobs Job, indexList int)  {
	for _, job := range jobs.deeps {
		if jobTable[job.name] == 0 {
			jobList = append(jobList,"")
			index := indexList - 1
			copy(jobList[index+1:],jobList[index:])
			jobList[index] = job.name
			indexList ++
			jobTable[job.name] = 1
		} else {
			for index, str := range jobList {
				if str == job.name {
					if index < indexList {

					} else {
						jobName := job.name
						jobList = append(jobList[:index-1],jobList[index+1:]...)
						jobList = append(jobList,"")
						index := indexList - 1
						copy(jobList[index+1:],jobList[index:])
						jobList[index] = jobName
						indexList ++
					}
				}
			}
		}
	}
	fmt.Println(jobList)
}

//多叉树处理
var jobList2 []string
var jobTable2 = make(map[string]int, 100)	//查表
var jobChange = make(map[string]int, 100)	//去重
type Job2 struct {
	parent []Job2
	name string
	deeps []*Job2
}


func sortJob2 (jobs []Job2) {
	for index, val := range jobs {
		jobTable2[val.name] = index
	}

	for index1, val1 := range jobs {
		for _, val2 := range val1.deeps {
			indexList := jobTable2[val2.name]
			jobs[indexList].parent = append(jobs[indexList].parent,jobs[index1])

		}
	}

	for _,val := range jobs {
		if val.parent == nil {
			fmt.Println(val.name)
			LRD(&val)
		}
	}
	fmt.Println(jobList2)
}

func LRD(root *Job2)  {
	//递归结束条件
	if root == nil {
		return
	}
	//进入下一层
	for _,child := range root.deeps {
		LRD(child)
	}
	//当前层逻辑
	_, err := jobChange[root.name]
	if err != true {
		jobList2 = append(jobList2,root.name)
		jobChange[root.name] = 1
	}

}
