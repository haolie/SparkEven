/**
 * Created by LYH on 2017/2/6.
 */
var async = require("async");
var nohelper = require('./nohelper.js');
var dbsuport = require('./MYSQLDBSuport.js');
var tool = require('./tools.js');
var lcsv = require('./localCSV.js');
//var xls_tool = require('./xlstool.js');
//var xls_tool = require('xls-to-json');
var process = require('process');
var fs = require('fs');
var path = require('path');
var util = require('util');

var worker = function () {
}
worker.prototype.isworking = false;
worker.prototype.console = function (item) {
  if (process.send)
    module.exports.sendMsg(item, "console");
  else
    console.log(item);
}

worker.prototype.sendMsg = function (msg, type) {
  process.send(JSON.stringify({
    id: process.pid,
    type: type,
    msg: msg
  }));
}

worker.prototype.saveToDb = function (item, allcallback) {
  module.exports.getValuesFromfile(item, function (err, items) {

    if (err) {
      module.exports.console(item.no + ": 获取文件数据失败");
      allcallback(1, 0);
      return;
    }


    dbsuport.saveTimePrice(items, function (err, result) {
        module.exports.console( 'saveTimePrice  '+ (new Date().valueOf()-module.exports.time))
      if (err) {
        module.exports.console( item.no+" 数据保存失败保存失败");
        allcallback(1, 1);
      }
      else {
        //  module.exports.console( item.no+" 保存成功");
        if (items.length) {
          var face = {
            _id: item.no + "_" + item.date,
            no: item.no,
            date: item.date,
            lastprice: item.lastprice,
            startprice: item.startprice,
            _min:item.min,
            _max: item.max,
            state: 1
          };

          //module.exports.console( face.startprice+"pricepricepricepricepricepricepriceprice");

            fs.open(path.join(__dirname, 'saveCatch/'+item.faceId), "a", function (err) {
                module.exports.console( 'fileTime  '+ (new Date().valueOf()-module.exports.time))
                allcallback(0, true);
            });

          // dbsuport.updatacodeface(face, function (err, s) {
			//     module.exports.console( 'updatacodeface  '+ (new Date().valueOf()-module.exports.time))
          //   allcallback(0, true);
          // });
        }
        else {
          allcallback(0, true);
        }
      }


    });
  })
}

worker.prototype.getValuesFromfile = function (item, allcallback) {

  if (!(item.no && item.date)) {
    allcallback(null, null);
    return;
  }
  var file = item.file;
  fs.exists(file, function (exist) {
    if (exist) {
      lcsv.FileToJson(file, function (err, result) {
        if(err){
            module.exports.console("出错出错出错出错" );
          try{
            fs.unlink(file,function () {
                module.exports.console("删除文件：" + file);
            });
          }catch (ex){
              module.exports.console("删除文件出错" + file);
          }

          allcallback(2, null)
          return;
        }
        result.shift();
        var items = [];
        result.forEach(function (row, index) {
          var time = new Date(item.date + " " + row[0]).getTime() / 1000;
          var t_type = 0;
          if (row[5] == "买盘") t_type = 1;
          if (row[5] == "卖盘") t_type = -1;
          items.push({
            _id: item.no + "_" + time,
            no: item.no,
            time: time,
            price: row[1],
            trade_type: t_type,
            turnover_inc: row[4],
            volume: row[3]
          })
        })

        try{
          fs.unlink(file,function () {});
        }catch (ex){}

        allcallback(0, items)
      })
    }
    else {
        module.exports.console("文件不存在："+file);
      allcallback(1, null);
    }
  })
}

worker.prototype.time=0

worker.prototype.start = function () {

  //2018-01-02_601636
  //  dbsuport.initCodesObj(function () {
  //      var temp=JSON.parse('{"faceId":1292400,"no":"688009","date":"2019-07-22","min":7.69,"max":15.21,"ud":6.42,"lastprice":12.27,"face":0,"dde_b":0,"dde_s":0,"volume":923624409,"turnoverRate":77.98,"turnover":9762393840,"state":0,"per":0,"startprice":11.7,"file":"E:/work/github/SparkEven/fvvm/datafiles/2019-07-22_688009.xls","i":3622,"savestate":1,"downstate":-1,"index":3622,"trytimes":0}');
  //      module.exports.saveToDb(temp, function (err, result) {
  //          module.exports.sendMsg({
  //              index: temp.index,
  //              result: err
  //          }, "result");
  //      });
  //
  //  })
  //
  //  return;

  process.on("message", function (msg) {
    //module.exports.console(msg);
    msg = JSON.parse(msg);
    var state = "free";
    if (msg.type == "work") {
      var item = msg.item;
      if (item != null) {
        state = "working";
          module.exports.time=new Date().valueOf()
        module.exports.saveToDb(item, function (err, result) {
          //保存成功 module.exports.console(" saveToDb saveToDb saveToDb saveToDb");
            module.exports.console( 'totaltime  '+ (new Date().valueOf()-module.exports.time))
          module.exports.sendMsg({
            index: item.index,
            result: err
          }, "result");
        });
      }
      module.exports.sendMsg(state, "state");
    }

  })

  module.exports.sendMsg("free", "state");
}

module.exports = new worker();
module.exports.start();
