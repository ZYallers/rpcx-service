#!/bin/bash

jsonFile="$(pwd)/src/service.json"
serviceName=""
serviceAddr=""
logDir=""

echoFun(){
    str=$1
    color=$2
    case ${color} in
        ok)
            echo -e "\033[32m $str \033[0m"
        ;;
        err)
            echo -e "\033[31m $str \033[0m"
        ;;
        tip)
            echo -e "\033[33m $str \033[0m"
        ;;
        title)
            echo -e "\033[42;34m $str \033[0m"
        ;;
        *)
            echo "$str"
        ;;
    esac
}

getJsonValue(){
    echo `cat ${jsonFile} | grep \"$1\"\: | sed -n '1p' | awk -F '": ' '{print $2}' | sed 's/,//g' | sed 's/"//g' | sed 's/ //g'`
}

helpFun(){
    echoFun "操作:" title
    echoFun "    status                                  查看服务状态" tip
    echoFun "    sync                                    同步服务vendor资源" tip
    echoFun "    build                                   编译生成服务程序" tip
    echoFun "    reload                                  平滑重启服务" tip
    echoFun "    quit                                    停止服务" tip
    echoFun "    help                                    查看命令的帮助信息" tip
    echoFun "有关某个操作的详细信息，请使用 help 命令查看" tip
}

initFun(){
    echoFun "service config:" title
    if [[ ! -f "$jsonFile" ]];then
        echoFun "file [$jsonFile] is not exist" err
        exit 1
    fi

    serviceName=$(getJsonValue "name")
    if [[ "$serviceName" == "" ]];then
        echoFun "serviceName is empty" err
        exit 1
    fi
    echoFun "ServiceName: $serviceName" tip

    sip=$(ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:"|head -1)
    if [[ "$sip" == "" ]];then
        echoFun "SystemIP is empty" err
        exit 1
    fi

    serviceAddr=$(getJsonValue "addr")
    if [[ "$serviceAddr" == "" ]];then
        echoFun "ServiceAddr is empty" err
        exit 1
    fi

    port=`echo ${serviceAddr}|awk -F ':' '{print $2}'`
    serviceAddr="$sip:$port"
    echoFun "ServiceAddr: $serviceAddr" tip

    logDir=$(getJsonValue "logDir")
    if [[ "$logDir" == "" ]];then
        echoFun "LogDir is empty" err
        exit 1
    fi

    echoFun "LogDir: $logDir" tip
}

statusFun(){
    echoFun "ps process:" title
    if [[ `pgrep ${serviceName}|wc -l` -gt 0 ]];then
        ps -p $(pgrep ${serviceName}|sed ':t;N;s/\n/,/;b t'|sed -n '1h;1!H;${g;s/\n/,/g;p;}') -o user,pid,ppid,%cpu,%mem,vsz,rss,tty,stat,start,time,command
    fi

    echoFun "lsof process:" title
    port=`echoFun ${serviceAddr}|awk -F ':' '{print $2}'`
    lsof -i:${port}
}

syncFun(){
    cd ./src
    echoFun "go mod vendor:" title
    if [[ ! -f "./go.mod" ]];then
        go mod init src
    fi
    go mod tidy
    rm -rf ./vendor
    go mod vendor
    echoFun "go mod vendor finished" ok
}

buildFun(){
    echoFun "git pull:" title
    branch=$1
    env=$2

    if [[ "$branch" == "" ]];then
        echoFun "branch of the build is empty" err
        exit 1
    fi
    if [[ "$branch" == "local" ]];then
        echoFun "ignore git pull, direct build by local" tip
    else
        git remote update origin --prune # 更新远程分支列表
        git checkout ${branch} # 切换分支
        git pull # 拉取最新版本
        echoFun "git pull [$branch] finish" ok
    fi

    echoFun "build runner:" title
    cd ./src

    tmpName="${serviceName}_tmp_$(date +'%Y-%m-%d-%H-%M-%S')"

    if [[ "$env" == "dev" ]];then
        echoFun 'build in develop environment' tip
        CGO_ENABLED=0 go build -v -installsuffix cgo -ldflags '-w' -i -o ../bin/${tmpName} -tags=jsoniter ./main.go
    else
        ### build编译参数参考资料：
        # 无依赖编译：https://blog.csdn.net/weixin_42506905/article/details/93135684
        # build参数详解：https://blog.csdn.net/zl1zl2zl3/article/details/83374131
        # ldflags参数：https://blog.csdn.net/javaxflinux/article/details/89177863
        CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -i -o ../bin/${tmpName} -tags=jsoniter ./main.go
    fi

    cd ../
    if [[ ! -f "./bin/${tmpName}" ]];then
        echoFun "build tmp runner ($(pwd)/bin/${tmpName}) failed" err
        exit 1
    fi

    mv -f ./bin/${tmpName} ./bin/${serviceName}
    if [[ ! -f "./bin/${serviceName}" ]];then
        echoFun "mv tmp runner failed" err
        exit 1
    fi

    echoFun "build runner ($(pwd)/bin/${serviceName}) finished" ok
}

sendMsg(){
    envVar="$(echo $hxsenv)"
    if [[ "$envVar" == "" ]];then
        envVar="development"
    fi
    env="Env: $envVar"
    name="Name: $serviceName"
    addr="Addr: $serviceAddr"
    hostnameVar=$(hostname)
    if [[ "$hostnameVar" == "" ]];then
        envVar="unknown"
    fi
    hostname="HostName: $hostnameVar"
    sipVar=$(ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:"|head -1)
    if [[ "$sipVar" == "" ]];then
        sipVar="unknown"
    fi
    sip="SystemIP: $sipVar"
    time="Time: $(date "+%Y/%m/%d %H:%M:%S")"
    token=$(getJsonValue "gracefulRobotToken")
    url="https://oapi.dingtalk.com/robot/send?access_token=$token"
    content="$1\n---------------------------\n$env\n$name\n$addr\n$hostname\n$time\n$sip"
    cnt=$(echo ${content//\"/\\\"})
    header="Content-Type: application/json"
    curl -o /dev/null -m 3 -s "$url" -H "$header" -d "{\"msgtype\":\"text\",\"text\":{\"content\":\"$cnt\"}}"
}

reloadFun(){
    echoFun "runner reloading:" title

    if [[ ! -f "./bin/$serviceName" ]];then
        echoFun "runner [`pwd`/bin/$serviceName] is not exist" err
        exit 1
    fi

    # 日志目录
    if [[ ! -d "$logDir" ]];then
        mkdir -p ${logDir}
    fi
    if [[ ! -d "$logDir" ]];then
        echoFun "logDir [$logDir] is not exist" err
        exit 1
    fi

    # 日志文件
    logfile=${logDir}/${serviceName}.log
    if [[ ! -f "$logfile" ]];then
        touch ${logfile}
    fi
    echoFun "logfile: $logfile" tip

    # 执行权限
    if [[ ! -x "./bin/$serviceName" ]];then
        chmod u+x ./bin/${serviceName}
    fi

    # 终止程序
    quitFun

    # 防止Jenkins默认会在Build结束后Kill掉所有的衍生进程
    export BUILD_ID=dontKillMe

    nohup ./bin/${serviceName} >> ${logfile} 2>&1 &
    echoFun "service $serviceName($serviceAddr) is reloaded, pid: `echo $!`" ok

    # 检查健康接口是否访问正常
    sleep 3s
    id=$(date +%Y%m%d%H%M%S)
    post="{\"jsonrpc\":\"2.0\",\"method\":\"$serviceName.health\",\"params\":{},\"id\":$id}"
    resp=`curl -m 3 -s "http://$serviceAddr" --header "X-JSONRPC-2.0:true" -d "$post"`
    echoFun "curl \"http://$serviceAddr\" health: $resp" tip
    sendMsg "curl \"http://$serviceAddr\" health: $resp"
}

quitFun(){
    port=`echo ${serviceAddr}|awk -F ':' '{print $2}'`
    counter=0
    while true;
    do
        pid=`lsof -i tcp:${port}|grep LISTEN|awk '{print $2}'`
        if [[ ${pid} -gt 0 ]];then
            if [[ ${counter} -ge 30 ]];then
                kill -9 ${pid}
                echoFun "service $serviceName has been killed for 30s and is ready to be forcibly killed" tip
                sendMsg "service $serviceName has been killed for 30s and is ready to be forcibly killed"
                break
            else
                kill ${pid}
                counter=$(($counter+1))
                echoFun "killing service $serviceName($port), pid($pid), $counter tried" tip
                sleep 1s
            fi
        else
            echoFun "service $serviceName($port) service is stopped" ok
            break
        fi
    done
}

initFun
case $1 in
    status)
        statusFun
    ;;
    sync)
        syncFun
    ;;
    build)
        buildFun $2 $3
    ;;
    quit)
        quitFun
    ;;
    reload)
        reloadFun
    ;;
    *)
        helpFun
    ;;
esac