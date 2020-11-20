var http = require('http');
var async = require("async");
var nohelper = require('./nohelper.js');
var url = require("url");
var path = require('path');
var fs = require('fs');
var tool = require('./tools.js');


var grouper = function () {

}

var group = function () {
  var date = "2018-11-20", path = "D:\\d\\", index = 0;
  fs.readFile(path + "old.txt", function (err, bytesRead) {

   // var oldcodes = bytesRead.toString().split("\r");
    var oldObj = {},oldcodes={}
    // oldcodes.forEach(function (c) {
    //   c=c.replace(/\n/,"");
    //   c=c.replace(/\n/,"");
    //   c=c.replace(/ /,"");
    //   oldcodes[c] = 1;
    // })

    nohelper.getallnofromweb("date", function (err, codes) {
      var list = [];
      codes.forEach(function (c) {

        if (!oldcodes[c.no]) list.push(c)
      })

      list = tool.getSpiedList(list, 500);


      async.mapLimit(list, 1, function (group, callback) {
        var str = "";
        group.forEach(function (g) {
          str += g.no + "\r\n "
        })
        index += 1;
        var tmp = path + date + "_" + index + ".txt";
        fs.writeFile(tmp, str, function (err) {
          console.log(tmp);
          callback(err, 1)
        });

      }, function (err, result) {
        console.log("finished")
      })
    })

  });
}

group();
