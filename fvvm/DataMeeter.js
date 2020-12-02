/**
 * Created by LYH on 2016/9/21.
 */

var http = require('http');
var async=require("async");
var hashmap=require("hashmap");
var nohelper = require('./nohelper.js');
var dbsuport = require('./MYSQLDBSuport.js');
var  process = require('process');
var url=require("url");
var zlib = require('zlib');
var path = require('path');
var fs= require('fs');
var fork = require('child_process').fork;
var request=require('request');
var childProgresscount=1;
var tools= require('./tools.js');
var downfile=1,savedb=1;

var stateObj={};

//savestate //-1：未处理  0：已下载 1：正在处理 2：已完成  4：处理错误

//
var tempCodes=['300207','600313','600510'];

var DataMeeter=function(){};

DataMeeter.prototype.checkValueDate=function(callback){
  var url="http://hq.sinajs.cn/list=sh000001";
  http.get(url,function(res){
    var length=0;
    var chunks=[];
    res.on('data', function (chunk) {
      length+=chunk.length;
      chunks.push(chunk);
    });
    res.on('end', function (str) {

      var str=chunks.toString();
      str=str.substring(str.length-25);
      str=str.substring(0,10);
      global.datestr=str;
      var curdate=new Date();
      var curstr=curdate.getFullYear()+"-"+(curdate.getMonth()+1)+"-"+curdate.getDate();
      //if(curstr==global.datestr&&curdate.getHours()<15){
      //    callback(null,true);
      //    return;
      //}

      //var datastr=date.getFullYear()+"-"+(date.getMonth()+1)+"-"+date.getDate();
      dbsuport.getcodeface(global.shcode,str,function(err,result){
        callback(null,result!=null&&result.state);
      })
    });
  });




}

var dC=function () {
  this.index=0,this.dates={},this.items=[],this.finished=false,this.saveDate='';
}

dC.prototype.addDate=function (date) {
  this.dates[date]={
    items:[],
    count:0,
    state:0,//0:为开始下载 1:下载中 2:已下载完成
  }
}

dC.prototype.addItem=function (item) {
  this.items.push(item);
}
dC.prototype.getItem=function(){
  if(this.finished||!this.items.length ) return null;
  var item=this.items.shift();
  this.saveDate=item.date;
  return item;
}


DataMeeter.prototype.dataContext=new dC();

DataMeeter.prototype.downDateFiles=function(date,callback){
  module.exports.console("start file down:"+ date);
  nohelper.getallno(date,function(err,codes){
    module.exports.dataContext.dates[date].count=codes.length;
    module.exports.dataContext.dates[date].state=0;


    if(err){
      callback(err,null);
      return;
    }

    var temps=[];

    codes.forEach(function(c,i){
      if(c.no==global.shcode) {
        module.exports.dataContext.dates[date].count-1;
        return;
      }
      c.file=0;
      c.i=i;
      c.downstate=c.savestate=-1;
      c.index=i;
      module.exports.dataContext.dates[date].items.push(c);
      if(c.state)
        c.downstate=c.savestate=2;
      else
        temps.push(c);
    });

    //待检测代码
    if(!downfile){
      module.exports.console("codes saved "+date);
      callback(0,1);
      return;
    }

    module.exports.console(date +": count="+ codes.length+"; needsave="+temps.length );
    if(temps.length==0){
      module.exports.dataContext.dates[date].state=2;
      dbsuport.updatacodeface({
        no:global.shcode,
        state:1,
        date:date
      },function(err,r){
        callback(null,true);
        module.exports.console(date+ " has save completed");
      });

      return;
    }

    var dateitem={date:date,items:temps};
    module.exports.dataContext.dates[date].state=1;
    async.mapLimit(dateitem.items,2,function(item,mapcb){
      var file= path.join(__dirname, "datafiles/"+date+"_"+item.no +".xls");
      fs.exists(file,function(exist){
        if(exist){

          module.exports.console("exist："+item.no + "  "+ item.i+"/"+dateitem.items.length);
          item.savestate=0;
          item.file=file;
          item.trytimes=0;
          module.exports.dataContext.addItem(item);
          mapcb(null,1);
        }
        else {
          var datetime=new Date(date);
          //  http://stock.gtimg.cn/data/index.php?appn=detail&action=download&c=sz000819&d=20170522
          //url="http://quotes.money.163.com/cjmx/" +
          //    datetime.getFullYear() + "/" +
          //    datetime.toLocaleDateString().replace(/-/g,'') +"/";

          var url="http://stock.gtimg.cn/data/index.php?appn=detail&action=download";
          var codestr="sz"+ item.no;
          if(Number(item.no)>=600000)codestr="sh"+ item.no;
          var datestr=date.replace('-','');
          datestr=datestr.replace('-','');
          url=url+"&c="+codestr+"&d="+datestr;
          var stream = fs.createWriteStream(file);
          request.get(url,{timeout:15000},function(err){
            if(err){
              module.exports.console("下载成功失败失败");
              module.exports.console("下载成功失败失败");
              module.exports.console("下载成功失败失败");
              module.exports.console("下载成功失败失败");
              module.exports.console("下载成功失败失败");
              module.exports.console("下载成功失败失败");
              module.exports.console("下载成功失败失败");
              module.exports.console("下载成功失败失败");
              module.exports.console("下载成功失败失败");
              module.exports.console("下载成功失败失败");
              item.savestate=-1;
              mapcb(null,0);

              fs.exists(file,function(exist){

                try{
                  if(exist) fs.unlink(file,function () {

                  });
                }catch(a) {}
              });
            }

          }).on("error",function(){
            module.exports.console("下载成功失败失败");
            module.exports.console("下载成功失败失败");
            module.exports.console("下载成功失败失败");
            module.exports.console("下载成功失败失败");
            module.exports.console("下载成功失败失败");
            module.exports.console("下载成功失败失败");
          }).pipe(stream).on('close', function(err,result){
            //request(url).pipe(stream).on('close', function(err,result){
            stream.end();
            if(err){
              item.savestate=-1;
              mapcb(null,0);


              fs.exists(file,function(exist){
                try{
                  if(exist) fs.unlink(file,function () {

                  });
                }catch(a) {}

              });
              return;
            }
            // module.exports.console("下载成功："+file+"  "+ item.i+"/"+dateitem.items.length);
            item.savestate=0;
            item.file=file;
            item.trytimes=0;

            module.exports.dataContext.addItem(item);
            mapcb(null,1);
          }).on("error",function (){

          });
        }
      })

    },function(err,result){
      module.exports.dataContext.dates[date].state=1;
      module.exports.console(dateitem.date+"  finid");
      var r=1;
      result.forEach(function (item) {
        r&=item;
      })
      callback(null,result);
    })
  });
}

DataMeeter.prototype.startFiledown=function(callback){
  module.exports.isdowning=true;
  module.exports.dataItems=[];
  module.exports.getQueryDates(function(err,dates){
    // dates=["2017-04-17","2017-04-18","2017-04-19","2017-04-20","2017-04-21"];
    if(dates==null||dates.length==0){
      callback(null,1)
      return;
    }

    dates.forEach(function(d,index){
      module.exports.dataContext.addDate(d);
      module.exports.console(d);
    })
    async.mapLimit(dates,1,module.exports.downDateFiles,function(err,result){
      callback(err,result);
      module.exports.isdowning=false;
    });
  });
}

DataMeeter.prototype.getQueryDates=function(callback){
  var tempdate=new Date(global.datestr);
  tempdate.add('d',-30);
  tempdate=tempdate.toLocaleDateString();
  tempdate="2019-12-20";
  nohelper.getwebDates(tempdate,function(err,dates){
    var date=[];
     // dates=dates.concat(['2019-01-02','2019-01-03','2019-01-04'])
    if(dates==null&&dates.length==0){
      date.push(global.datestr);
      callback(null,date)
      return;
    }

    var allDates={};
    dbsuport.getfaces({no:global.shcode},function(err,items){
      items.forEach(function (item) {
        allDates[item.date]=1
        if(item.state>0)allDates[item.date]=2;
        if(item.state<=0){
          date.push(item.date);
        }
      })

      dates.forEach(function (d) {
        if(!allDates[d]){
          date.push(d);
        }
      })
      callback(null,date);
    });
  });
}

DataMeeter.prototype.updataFaces=function(callback){
 fs.readdir(path.join(__dirname, 'saveCatch'),function (err,files) {
   if(!files.length){
       callback()
       return
   }
    var arraylist= tools.getSpiedList(files,100)
    async.mapLimit(arraylist,1,function (list,cb) {
        dbsuport.updateFaceState(list,function (err) {
         cb()
        })
    },function () {
        async.mapLimit(files,5,function (file,dcb) {
            try{
                fs.unlink(path.join(__dirname, 'saveCatch/'+file),function () {
                    dcb()
                })
            }catch (e) {
                dcb()
            }

        },function () {
            callback()
        })
    })
 })
}


DataMeeter.prototype.consoleTimes=function(str){
  var ms=new Date().valueOf()-starttime.valueOf();
  var m=0;
  var s=0;
  s=Math.floor(ms/1000);
  ms=Math.floor(ms%1000);
  m=Math.floor(s/60);
  s=Math.floor(s%60);

  module.exports.console(str+ " "+m+":"+s+";"+ms);
}

DataMeeter.prototype.console=function(item){
  if(process.send)
    process.send(item);
  else
    console.log(item);
}

DataMeeter.prototype.startwork=function(){
  module.exports.isWorking=true;
  module.exports.checkValueDate(function(err,result){

    module.exports.console("checkValueDate:"+result);
    //if(result) {
    //    module.exports.isWorking=false;
    //    return;
    //}

    //module.exports.dataContext.finished=false;
    //module.exports.dataContext.items=[];
    //module.exports.dataContext.index=0;


    module.exports.console("start data meet");
    var downCall=function(err,dates){
      // if(err){
      //     module.exports.startFiledown(downCall)
      // }
      // else {
      //     module.exports.console("downFinished");
      //    // module.exports.commitDateItems();
      //     module.exports.isWorking=false;
      // }
    }

    module.exports.startFiledown(downCall)
  });
}

DataMeeter.prototype.sendWorker=function(p){
  var work=module.exports.dataContext.getItem();
  if(work)work.savestate=1;
  else work=null;


    if(work&&work.date=='2019-12-13'){
      debugger
    }
  p.item=work;
  p.worker.send(JSON.stringify({type:"work",item:work}));

  return true;
  // tool.console(JSON.stringify({type:"work",item:work}));
};


DataMeeter.prototype.createChild=function(p){
  if(p==undefined||p==null){
    p={};
    module.exports.progress.push(p) ;//启动子进程
  }
  p.state="working";
  p.worker= fork( path.join(__dirname, "DataMeeter_worker.js"))
  p.worker.on("message",function(msg,b){
    //console.log(msg);
      try{
          msg=JSON.parse(msg);
          if(msg.type=="state"){
              module.exports.progress.forEach(function(p,index){
                  if(p.worker.pid==msg.id){
                      p.state=msg.msg;
                  }
              })
          }
          else  if(msg.type=="result"){
              //module.exports.console(JSON.stringify(msg));
              module.exports.progress.forEach(function(p,index){
                  if(p.worker.pid==msg.id){
                      var tm= p.item;
                      if(!tm)return;
                      tm.savestate=2;//-1：未处理  0：已下载 1：正在处理 2：已完成
                      var ts=" 保存成功";
                      if(msg.msg.result){
                          tm.trytimes+=1;
                          if(tm.trytimes>=3)msg.savestate=4;
                          else {
                              msg.savestate=0;
                              module.exports.dataContext.items.push(tm)
                          }
                          ts=" 保存失败 times "+tm.trytimes
                      }

                      module.exports.console(tm.date+" "+tm.no+ ts +" "+tm.index);
                      p.item=null;
                      module.exports.sendWorker(p);
                  }
              })

          }
          else  if(msg.type=="console"){
              module.exports.console(msg.msg);
          }
      }catch (e) {
          module.exports.console(msg);
      }


  });
  p.worker.on("create",function(){
    p.state="free";
  })
  p.worker.on("exit",function(){
    if(p.item){
      try {
        fs.unlink(p.item.file,function () {

        });
      }
      catch (e){}
      p.item.savestate=0;
    }
    p.worker.kill();
    p.worker=null;
    module.exports.createChild(p);
  });

}

DataMeeter.prototype.startChildWorker=function(){
  for (var i = 0; i < childProgresscount; i++) {
    module.exports.createChild();
  }

  var fun=function () {
      setTimeout(function () {
          module.exports.updataFaces(function () {
              fun()
          })
      },5000)
  }

  //fun()
}

DataMeeter.prototype.isdowning=false;
DataMeeter.prototype.isWorking=null;
DataMeeter.prototype.progress=[];
DataMeeter.prototype.start=function() {
  module.exports.isWorking = true;

  //setInterval(function () {
  //    if (module.exports.isWorking) return;
  //    module.exports.startwork();
  //}, 180000);
  module.exports.startwork();
  module.exports.startChildWorker();
  //
  //

  var nulltimes=0;
  var fun=function () {
    if(module.exports.dataContext.finished) return;

    var state=0;
    module.exports.progress.forEach(function(p,index){
      if(p.state=="free"){
        state|= module.exports.sendWorker(p);
      }
      else state=1;
    })

    if(!state)nulltimes+=1;
    else nulltimes=0;
    if(nulltimes>8)   module.exports.dataContext.finished=!module.exports.isWorking;

    setTimeout(fun,1000);
  }
  setTimeout(fun,10000);
}

// var fileset=function(){
//   var path='';
//
//   fs.
//
//
// }

module.exports=new DataMeeter();
nohelper.startCodeIdListen()
dbsuport.initCodesObj(function () {
    module.exports.updataFaces(function () {
        module.exports.start();

    })
})


