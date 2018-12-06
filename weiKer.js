var fork = require('child_process').fork;
var tool= require('./fvvm/tools.js');
var path = require('path');

var meeter= null

var statrMeeter=function(){
    tool.console("statrMeeter")
    meeter=  fork( path.join(__dirname, "fvvm/DataMeeter_worker.js"))
    meeter.on('console',function (msg) {
        tool.console(msg)
    })
    meeter.on("message",function(msg,b){
        msg=JSON.parse(msg);
        if(!msg.type){
            tool.console(msg)
        }
        else if(msg.type=="state"){
        }
        else  if(msg.type=="result"){
        }
        else  if(msg.type=="console"){
            tool.console(msg.msg);
        }
    });
    meeter.on("create",function(){
        tool.console("create")
    })
    meeter.on("exit",function(a,b){
        tool.console("exit")
        setTimeout(statrMeeter,30000)
    });
}

statrMeeter()
var t=new Date().valueOf()
console.log(t)

