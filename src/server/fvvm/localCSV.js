/**
 * Created by youhao on 2017/5/30.
 */

var async=require("async");
var tool = require('./tools.js');
var fs= require('fs');
var readline = require('readline');
var iconv = require('iconv-lite');


CSV=function(){

}

CSV.prototype.FileToJson=function(file,callback,code){
    if(!code)code='gbk';
    try {
        fs.readFile(file, function (err,bytesRead) {
            var string = iconv.decode(bytesRead, code);
            var strs= string.split('\n')
            var items=[];
            for (var i=0;i<strs.length;i++){
                var temp=strs[i].split('\t');
                items.push(temp);
            }
            if(callback)callback(err,items);
        });
    }
    catch(ex){
        callback(1,null);
    }
}

CSV.prototype.saveToCSV=function(objs,property,savepath,callback){
    var rows=[];
    var str=property.join("\t");
    rows.push(str);

    objs.forEach(function(item,index){
        var temp=[];
        for (var i=0;i<property.length;i++)
        {
            if(item[property[i]]==undefined) temp.push("   ")
            else temp.push(item[property[i]]);
        }
        rows.push(temp.join("\t"));
    })

    fs.writeFile(savepath, rows.join('\n'), function (err) {
        callback(err,savepath);
    });

}

module.exports=new CSV();