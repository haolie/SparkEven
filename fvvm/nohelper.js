/**
 * Created by LYH on 2016/10/10.
 */
var http = require('http');
var async=require("async");
var percount=1000;
var token="";
var fs= require('fs');
var path = require('path');
var dbsuport = require('./MYSQLDBSuport.js');
var tools = require('./tools.js');
var request=require('request');
var querylimit=1;
var repl = require('repl');
require('date-utils');

var columnObj=[
    {name:'涨跌',index:-1,id:'ud'},
    {name:'收盘价:不复权',index:-1,id:'lastprice'},
    {name:'开盘价:不复权',index:-1,id:'startprice'},
    {name:'最高价:不复权',index:-1,id:'max'},
    {name:'最低价:不复权',index:-1,id:'min'},
  {name:'换手率',index:-1,id:'turnoverRate'},
  {name:'成交额',index:-1,id:'turnover'},
  {name:'成交量',index:-1,id:'volume'}
    ]

var nohelper=function Nohelper(){}

nohelper.prototype.currentDate=null;
nohelper.prototype.allItems=null;

nohelper.prototype.getallnofromweb=function(date,callback){
    var index=1;
    var pageCount=0;
    var allcodes=[];
    this.currentDate=date;
    this.allItems={};
    this.getToken(date,function(err,tokenObj){
        tokenObj=tokenObj.data.result;
        token=tokenObj.token;
        pageCount=Math.floor(tokenObj.code_count/percount);
        if(Math.floor(tokenObj.code_count%percount)>0)
          pageCount+=1;

        var indexarray=[];
        for(var i=1;i<=pageCount;i++){
            indexarray.push(i);
        }

        columnObj.forEach(function (co) {
            co.index=-1
        })

        tokenObj.columnsIndexID.forEach(function (col,i) {
            if(col.timestamp==date.replace(/-/g,"")){
               columnObj.forEach(function (co) {
                   if(co.name==col.index_name)co.index=i
               })
            }
        })

        async.mapLimit(indexarray,querylimit,module.exports.getno,function(err,results){
            for(var i=0;i<results.length;i++){
                allcodes= allcodes.concat(results[i]);
            }


            callback(null,allcodes);
        })



    });

};



nohelper.prototype.getno=function(index,callback){

    var fun=function (n) {
        var url="/stockpick/cache?token=" +token
            +"&p=" +index +
            "&perpage=" +
            percount +
            "&showType=[%22%22,%22%22,%22onTable%22,%22onTable%22,%22onTable%22,%22onTable%22]";

        var  options = {
            hostname:'www.iwencai.com',
            port: 80,
            method: 'get',
            path: url,
            KeepAlive: true,
            headers: {
                'Connection':"keep-alive",
                'Host':"www.iwencai.com",

                "User-Agent":'Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36',
                'Cookie':'PHPSESSID=8kj3mqbbsij14id6f3403u03f1; v=AsnAGHLn926_m4sa8g40WcE62P4nFr1IJwrh3Gs_RbDvsudos2bNGLda8ar7; cid=nc6hgumel0h6oscbpn8fhibtf31497502968; ComputerID=nc6hgumel0h6oscbpn8fhibtf31497502968'
            }
        };

        http.get(options,function(resp){
            var length=0;
            var chunks=[];
            resp.on("data",function(chunk){
                length+=chunk.length;
                chunks.push(chunk);
            })

            resp.on("end",function(){
                var buf=Buffer.concat(chunks,length);
                //console.log(buf.toString())
                var temp="";
                try {
                    temp=JSON.parse(buf.toString().toLowerCase());
                }catch(e) {
                    setTimeout(function(){
                        if(n>0){
                            fun(n-1);
                        }
                        else throw e;
                    } ,200)
                    return;
                }
                var result=[];

                var indexFun=function (name) {
                    for(var i=0;i< temp.columnsindexid.length;i++){
                        if(temp.columnsindexid[i].index_name==name)return i;
                    }

                }

                var index_start=indexFun("开盘价:前复权"),
                    index_end=indexFun("收盘价:前复权"),
                    index_updown=indexFun("涨跌"),
                    index_codes=0,
                    index_max=indexFun("最高价:前复权"),
                    index_min=indexFun("最低价:前复权"),
                    index_vlo=indexFun("成交量");

                for(var i=0;i<temp.result.length;i++){
                    if(temp.result[i][columnObj[0].index]=="--") continue;
                    var codestr=temp.result[i][0].toString();
                    codestr=codestr.replace(".sz","");
                    codestr=codestr.replace(".sh","");
                    codeObj={no:codestr,
                        date:module.exports.currentDate,
                        state:0,
                        index:result.length}
                    columnObj.forEach(function (col) {
                        codeObj[col.id]= temp.result[i][col.index]
                    })
                    result.push(codeObj)
                }
                if(callback)
                    setTimeout(function(){callback(null,result)} ,500)

            })
        })


    }

    fun(5);
}

nohelper.prototype.getToken=function(date,callback){
    //var url="http://www.iwencai.com/stockpick/search?typed=1&preParams=&ts=1&f=1&qs=result_rewrite&selfsectsn=&querytype=&searchfilter=&tid=stockpick&w=%E5%87%80%E9%87%8F";
    // var url='/stockpick/search?typed=1&preParams=&ts=1&f=3&qs=pc_%7Esoniu%7Estock%7Estock%7Ehistory%7Equery&selfsectsn=&querytype=&searchfilter=&tid=stockpick&w=dde' +
    //     date +
    //     '+%E6%B6%A8%E8%B7%8C' +
    //     date;


    //2018年7月6日
   // var url= 'http://www.iwencai.com/stockpick/load-data?typed=1&preParams=&ts=1&f=1&qs=result_rewrite&selfsectsn=&querytype=&searchfilter=&tid=stockpick&w=dde' +
   //     date +
   //     '+%E6%B6%A8%E8%B7%8C' +
   //     date +
   //     '&queryarea='

    // var url= 'http://www.iwencai.com/stockpick/load-data?typed=1&preParams=&ts=1&f=1&qs=result_rewrite&selfsectsn=&querytype=&searchfilter=&tid=stockpick&w=%E6%B6%A8%E8%B7%8C%20' +
    //     date +
    //     '%20%E6%B6%A8%E5%B9%85' +
    //     date +
    //     '&queryarea='
   // var url= 'http://www.iwencai.com/stockpick/load-data?typed=1&preParams=&ts=1&f=1&qs=result_rewrite&selfsectsn=&querytype=&searchfilter=&tid=stockpick&w=dde' +
    //    //     date +
    //    //     '+%E6%B6%A8%E8%B7%8C' +
    //    //     date +
    //    //     '&queryarea='
    var url= 'http://www.iwencai.com/stockpick/load-data?typed=0&preParams=&ts=1&f=1&qs=result_original&selfsectsn=&querytype=stock&searchfilter=&tid=stockpick&w=%E6%B6%A8%E8%B7%8C' +
      date +
      '+%E4%BB%B7%E6%A0%BC' +
      date +
      '+%E6%88%90%E4%BA%A4%E9%87%8F' +
      date +
      '+%E6%8D%A2%E6%89%8B' +
      date +
      '+&queryarea='

   // url='/stockpick/search?typed=1&preParams=&ts=1&f=3&qs=pc_%7Esoniu%7Estock%7Estock%7Ehistory%7Equery&selfsectsn=&querytype=&searchfilter=&tid=stockpick&w=dde2017-07-25+%E6%B6%A8%E8%B7%8C2017-07-25';
    var length=0;
    var chunks=[];
    var  options = {
        hostname:'www.iwencai.com',
        port: 80,
        method: 'get',
        path: url,
        KeepAlive: true,
        ondata:function (chunk) {
            length+=chunk.length;
            chunks.push(chunk);
        },
        headers: {
            'Connection':"keep-alive",
            'Host':"www.iwencai.com",

            "User-Agent":'Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36',
            //'Cookie':'PHPSESSID=8kj3mqbbsij14id6f3403u03f1; v=AsnAGHLn926_m4sa8g40WcE62P4nFr1IJwrh3Gs_RbDvsudos2bNGLda8ar7; cid=nc6hgumel0h6oscbpn8fhibtf31497502968; ComputerID=nc6hgumel0h6oscbpn8fhibtf31497502968'
            'Cookie':'cid=2eb414ff675d8125147fc689c2569d7a1538712169; ComputerID=2eb414ff675d8125147fc689c2569d7a1538712169; PHPSESSID=4a956d32cb5bff85cbf6de4cd97716c3; v=AiKotn5ejcUgMZEZsAsIJiJsc6OHcyQW2H8a4my6ThVAP8gVVAN2nagHatM_'
        }
    };

    var requsetCallBack=function (err,r) {
        var buf=Buffer.concat(chunks, length);
        var str= buf.toString();
        //str=str.substring(str.indexOf("var allResult ="));
        //str=str.substring(str.indexOf("{"),str.indexOf(";\n"));

        try {
            var temp=JSON.parse(str);
            if(callback)callback(null,temp);
        }
        catch (e) {
            var r = repl.start({ prompt: 'token  failed> ', eval: function myEval(cmd, context, filename, callback) {
                    callback(null, cmd);
                }, writer: function myWriter(output) {
                    options.headers.Cookie=output
                    r.close()
                    tools.HttpRequest(options,requsetCallBack)
                    return
                } });
        }

    }

    tools.HttpRequest(options,requsetCallBack)

    return;

   // url+="涨跌";
   http.get(options,function(resp){
       var length=0;
       var chunks=[];
       resp.on("data",function(chunk){
           length+=chunk.length;
           chunks.push(chunk);
       });
       resp.on("end",function(){
           var buf=Buffer.concat(chunks, length);
           var str= buf.toString();
           str=str.substring(str.indexOf("var allResult ="));
           str=str.substring(str.indexOf("{"),str.indexOf(";\n"));
           var temp=JSON.parse(str);

           if(callback)callback(null,temp);
       })
   })
}

nohelper.prototype.savenofile=function(items,callback){
    fs.writeFile(path.join(__dirname, 'history/'+global.datestr), JSON.stringify(items), function (err) {
        if(callback){
            callback(err);
        }
    });

}

nohelper.prototype.getwebDates=function(start,callback){
   // start="2017-01-25";
    var uri="http://quotes.money.163.com/service/chddata.html?code=0000001&start=" +
        new Date(start).toFormat("YYYY-MM-DD").replace(/-/g,"")+
        "&end=" +
        new Date(global.datestr).toFormat("YYYY-MM-DD").replace(/-/g,"")+
        "&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;VOTURNOVER;VATURNOVER";

    var file=path.join(__dirname, 'datafiles/sh000001.xls')
    console.log(file)
    var stream = fs.createWriteStream(file);
    request(uri).pipe(stream).on('close', function(err,result){
        fs.readFile(file, function (err,bytesRead) {
            var dates=[];
            var strs= bytesRead.toString("utf8").split("\r\n");
            for (var i=strs.length-1;i>=1;i--){
                var temp=strs[i].split(',');
                if(temp.length>2) dates.push(temp[0]);
            }
            callback(null,dates)
        });
    });

}

nohelper.prototype.getallnofromlocal=function(datestr, callback){
    fs.exists(path.join(__dirname, 'history/'+datestr),function(exist){
        if(exist){
            fs.readFile(path.join(__dirname, 'history/'+datestr), function (err,bytesRead) {
                callback(err,JSON.parse( bytesRead.toString()))
            });
        }
        else {
            callback(1,null);
        }
    })
}

nohelper.prototype.getallno=function(date,getallnocallback){
    // module.exports.getallnofromweb(date,function(err,items){
    //     if(err)getallnocallback("获取出错",null);
    //     else{
    //         items.push({no:global.shcode,state:0 ,date:items[0].date})
    //         dbsuport.savecodefaces(items,function(err,result){
    //             getallnocallback(null,items);
    //         })
    //     }
    // })
    //
    // return;

    dbsuport.getfaces({date:date},function(err,items){
        var temps=[];
        if(err==null&&items&&items.length>0){

            //var temps=[];
            //var obj={};
            //items.forEach(function(d,i){
            //    obj[d.no]=d;
            //})
            //
            //module.exports.getallnofromweb(date,function(err,webitems){
            //
            //    webitems.forEach(function(d,i){
            //        if(obj[d.no])
            //            temps.push(obj[d.no]);
            //        else
            //            temps.push(d)
            //    })
            //    getallnocallback(err,temps);
            //});

            getallnocallback(err,items);

        }
        else
            // setTimeout(function (args) {
            //     getallnocallback("获取出错",null);
            // },500)
        module.exports.getallnofromweb(date,function(err,items){
            if(err)getallnocallback("获取出错",null);
            else{
                items.push({no:global.shcode,state:0 ,date:items[0].date})
                dbsuport.savecodefaces(items,function(err,result){
                    getallnocallback(null,items);
                })
            }
        })
    })
}




module.exports=new nohelper();

// module.exports.getToken("2017-01-23",function (a,ba) {
//
// })

//global.datestr="2017-01-26"
//module.exports.getwebDates("2017-01-23")
