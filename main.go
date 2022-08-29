package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"go-watchdog/config"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func execcomm (name  string, arg1,arg2 string) *exec.Cmd {
	//log.Println("exec commond:  ",name,arg1,arg2)
	return exec.Command(name,arg1,arg2)
}
func svcsoperate (name,operate string) *exec.Cmd {
	return execcomm("/usr/sbin/service",name,operate)
	//cmdrun(cmd)
}
func cmdrun (cmd *exec.Cmd) {
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}
}
func cmdrunresult(cmd *exec.Cmd) int {
	stdout, _ :=cmd.StdoutPipe()
	cmdrun(cmd)
	result, _ := ioutil.ReadAll(stdout) // 读取输出结果
	printrestolog(cmd,result)
    defer stdout.Close()
	return checkrunstatus(cmd)
}
func checkrunstatus(cmd *exec.Cmd) int {
	var res int
	if err := cmd.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			log.Println("cmd exit status, Failed cmd Args: ",cmd.Args)
			res = ex.Sys().(syscall.WaitStatus).ExitStatus() //获取命令执行返回状态，相当于shell: echo $?
		}
	}
    return res
}
func printrestolog (cmd *exec.Cmd , result []byte) {
	resdata := string(result)
	log.Println("cmd run Arg: ",cmd.Args,"  result:" ,resdata)
}


func execcomms (name1,arg1,name2,arg2,name3,arg3,arg3_1,name4,arg4 string) []*exec.Cmd {
	var 	commands 	[]*exec.Cmd
	commands = append(commands, exec.Command(name1, arg1))
	commands = append(commands, exec.Command(name2, arg2))
	commands = append(commands, exec.Command(name3, arg3,arg3_1))
	commands = append(commands, exec.Command(name4, arg4))
	return commands
}
func psstatus (name string) []*exec.Cmd {
	if (name == "restart"){
		fmt.Println("not restart")
		name="status"
	}

	log.Println("/usr/bin/ps","-ef","/usr/bin/grep",name,"/usr/bin/grep" ,"-v","grep","/usr/bin/wc","-l")
	return execcomms(
		"/usr/bin/ps","-ef",
		"/usr/bin/grep",name,
		"/usr/bin/grep" ,"-v","grep",
		"/usr/bin/wc","-l",
	)
}
func cmdsrunresult(commands []*exec.Cmd)  int {
	for i:= 1;i < len(commands);i++{
		commands[i].Stdin, _ = commands[i-1].StdoutPipe()
	}
	commands[len(commands)-1].Stdout=os.Stdout
	res :=cmdsrun(commands)
	return res
}
func cmdsrun (commands []*exec.Cmd) int {
	for i:=1;i<len(commands);i++{
		err := commands[i].Start()
		//log.Println(commands[i].Process.Pid)
		if err != nil {
			return 1
		}
	}
	commands[0].Run()

	var res int
	for i:=1;i<len(commands);i++{
		err := commands[i].Wait()
		if err != nil {
			if ex, ok := err.(*exec.ExitError); ok {
				res = ex.Sys().(syscall.WaitStatus).ExitStatus() //获取命令执行返回状态，相当于shell: echo $?
			}
			log.Println(err)
			return 1
		}
	}
	return res
}
func caseconfigwatchdog(arg,operate string, conf *config.Config) int {
	switch {
	case strings.EqualFold(arg,"sshd")	:
		log.Println("will be to check :'",arg,"' ",operate)
		return watchdogcheck(arg,conf.Sshdsvc,conf.Sshdpsst,operate)
	case strings.EqualFold(arg,"mysqld")	:
		log.Println("will be to check :'",arg,"' ",operate)
		return watchdogcheck(arg,conf.Mysqldsvc,conf.Mysqldpsst,operate)
	case strings.EqualFold(arg,"grafana")	:
		log.Println("will be to check :'",arg,"' ",operate)
		return watchdogcheck(arg,conf.Grafanadsvc,conf.Grafanapsst,operate)
	case strings.EqualFold(arg,"nodemanager")	:
		log.Println("will be to check :'",arg,"' ",operate)
		return watchdogcheck(arg,conf.Nodemanagersvc,conf.Nodemanagerpsst,operate)
	case strings.EqualFold(arg,"resourcemanager")	:
		log.Println("will be to check :'",arg,"' ",operate)
		return watchdogcheck(arg,conf.Resourcemmanagersvc,conf.Resourcemmanagerpsst,operate)
	case strings.EqualFold(arg,"namenode")	:
		log.Println("will be to check :'",arg,"' ",operate)
		return watchdogcheck(arg,conf.Namenodesvc,conf.Nanenodepsst,operate)
	case strings.EqualFold(arg,"datanode")	:
		log.Println("will be to check :'",arg,"' ",operate)
		return watchdogcheck(arg,conf.Datanodesvc,conf.Datanodepsst,operate)
	case strings.EqualFold(arg,"sentry")	:
		log.Println("will be to check :'",arg,"' ",operate)
		return watchdogcheck(arg,conf.Sentrysvc,conf.Sentrypsst,operate)
	default	:
		log.Fatalln("not a avaliable value:  ",arg,"  cehck work will skip")
		return -1
	}
}

func watchdogcheck(program,arg1,arg2,operate string) int{
	statussultsum :=0
	switch  {
	case operate=="status" :
		if !(len(arg1)==0){
			//log.Println("program '",program,"' service status check")
			cmd := svcsoperate(arg1,operate)
			//log.Println("后面还要开发的逻辑判断符",cmdrunresult(cmd))
			statussultsum +=cmdrunresult(cmd)
		}else {
			log.Println("program '",program,"' service status check argument is wrond check skip")
		}
		//log.Println("查看一次执行结果",statussultsum)
		if !(len(arg2)==0){
			//log.Println("program '",program,"' ps ef status check")
			commands := psstatus(arg2)
			//log.Println("后面还要开发的逻辑判断符",cmdsrunresult(commands))
			statussultsum +=cmdsrunresult(commands)
		}else {
			log.Println("program '",program,"' ps ef  status check argument is wrond check skip")
		}
		//log.Println("查看两次执行结果",statussultsum)
		return statussultsum
	case operate=="restart" :
		if !(len(arg1)==0){
			log.Println("program '",program,"' service will be restart")
			cmd := svcsoperate(arg1,operate)
			statussultsum +=cmdrunresult(cmd)
		}else {
			log.Println("program '",program,"' service restart argument is wrond restart skip")
		}
		return statussultsum

	default :
		return -999
	}
}

var(
	Cfg 		*ini.File

)

func main() {
	var err error

	for (true) {

	Cfg,err=ini.Load("conf.ini")
	if err != nil{
		log.Fatal("Fail to Load ‘conf/app.ini’:",err)
	}
	config :=new(config.Config)
	err=Cfg.MapTo(config)
	//log.Printf("%T",config)
	//log.Println(config)




	for i:=0;i<len(config.App);i++{
		check :=caseconfigwatchdog(config.App[i],"status",config)
		switch  {
		case  check > 0 :
			log.Println(config.App[i]," program check is Error")
			log.Println("service restart ",config.App[i])
			caseconfigwatchdog(config.App[i],"restart",config)
		case check == 0 :
			log.Println(config.App[i]," program check is succeeded")
		case check < 0 :
			log.Println(config.App[i]," program check is WARN check program all wround not in watchdog list")
		}

	}


	time.Sleep(config.Polling* time.Minute)
	}
}

