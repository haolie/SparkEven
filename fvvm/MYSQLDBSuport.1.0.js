var mysql = require('mysql');
var http = require('http');
var async = require("async");
var url = require("url");
var tool = require('./tools.js');
var dateutil = require('date-utils');
var co=require("co");
var codeface = "codeface";
var time_price = "time_price";
var process = require('process');
var util = require('util');
global.shcode = 912261;

var suporter = function () {

}

suporter.prototype.connctions = [null, null, null,null];
suporter.prototype.cindex = 0;
suporter.prototype.getConnction = function (callback) {
    var curindex = module.exports.cindex;
    module.exports.cindex += 1;
    if (module.exports.cindex >= module.exports.connctions.length)module.exports.cindex = 0;

    if (module.exports.connctions[curindex] == null) {
        module.exports.connctions[curindex] = mysql.createConnection({
            host: 'localhost',
            user: 'mysql',
            password: '123456',
            database: 'finance',
            multipleStatements:true,
            useConnectionPooling: true
        });
        module.exports.connctions[curindex].connect(function (a, b) {
            callback(null, module.exports.connctions[curindex]);
        });
    }
    else
        callback(null, module.exports.connctions[curindex]);
}

suporter.prototype.CodesObj=null;

suporter.prototype.initCodesObj=function (callback) {
     if(current.CodesObj){
         if(callback) callback()
         return;
     }
    current.getConnction(function (err,conn) {
        conn.query('select * from tbl_codes',function (err,items) {
            current.CodesObj={};
            items.each(function (item) {
                current.CodesObj[item._no]=item.id;
            })
        })
    })
}

suporter.prototype.getNoById=function (id) {
    for(var i in current.CodesObj){
        if(current.CodesObj[i]==id)return current.CodesObj[i];
    }

    return "";
}

suporter.prototype.getInsertStr = function (item) {
    if (Number(item.no) < 1000000)item.no = Number(item.no) + 1000000;
    return '(' +
        Number(item.no) + ',' +
        '"' + tool.convertToTIMESTAMP(item.time) + '"' + ',' +
        (item.price * 100) + ',' +
        item.trade_type + ',' +
        (item.turnover_inc / 100) + ',' +
        (item.volume / 100) +
        ')';
}

suporter.prototype.deletebyno = function (no, callback) {
    module.exports.getConnction(function (err, conn) {
        var str = 'delete from time_price where no=' + no;
        conn.query(str, function (err, result) {
            callback(err, result);
        })
    })
}

suporter.prototype.saveTimePrice = function (timevalues, allcallback) {
    if (timevalues == null || timevalues.length == 0) {

        if (allcallback)allcallback(1, "");
        return;
    }

    var lists = tool.getSpiedList(timevalues, 600);

    module.exports.getConnction(function (err, conn) {
        if (err == null) {

            async.mapLimit(lists, 1, function (items, callback) {
                var insert = 'replace INTO time_price(id,_no,state) VALUES';
                for (var i in items) {
                    insert += module.exports.getInsertStr(items[i]);
                    if (i == items.length - 1) insert += ';';
                    else insert += ',';
                }
                conn.query(insert, function (err, result) {
                    callback(err, insert);
                })
            }, function (err, result) {
                if (allcallback)allcallback(err, "");
            })
        }
        else {
            if (allcallback)allcallback(null, "");
            console.log("数据库连接失败");
        }
    });
};


suporter.prototype.getValueByDayNo = function (item, callback) {
    var date = item.date;
    if (typeof (date) != 'string')
        date = date.toLocaleDateString();
    date = new Date(date);
    module.exports.getConnction(function (err, conn) {
        var no = Number(item.no);
        if (no < 1000000)no += 1000000;
        var start = date.toLocaleDateString();
        date.add('d', 1)
        var end = date.toLocaleDateString();
        var str = 'select * from time_price where no=' + no + ' and time>"' + start + '" and time<"' + end + '" order by time;';
        conn.query(str, function (err, result) {
            var items = [];
            if (result && result.length > 0) {
                for (var i in result) {
                    items.push({
                        no: result[i].no.toString().substring(1),
                        time: result[i].time,
                        price: result[i].price / 100,
                        trade_type: result[i].trade_type,
                        turnover_inc: result[i].turnover_inc * 100,
                        volume: result[i].volume * 100
                    })
                }
            }

            callback(err, items);
        });
    });
}

suporter.prototype.checkCount = function (item, callback) {
    module.exports.getConnction(function (err, conn) {

    });

    MongoClient.connect(dburl, function (err, db) {
        var filter = {
            "no": item.no,
            "time": {$gte: Number(item.start), $lte: Number(item.end)},

            //"_id":item.no+"_"+item.start,1477359297
            //"time":{$gte:1477357237},
            //"time":{$gte:1477359297},
            //"price":4.5
        };

        if (err) {
            callback(err, 0);
        }
        else {
            db.collection("time_price").find(filter).count(callback)
        }
        db.close();
    });
}


suporter.prototype.transData = function (code, callback) {

}

suporter.prototype.convertface = function (results) {
    var items = [];
    if (results && results.length)
        for (var i in results) {
            items.push({
                no: results[i]._no.toString().substring(1),
                date:  new Date(results[i]._date).toFormat("YYYY-MM-DD"),
                min: results[i]._min / 100,
                max: results[i]._max / 100,
                ud: results[i].ud / 100,
                lastprice: results[i].lastprice / 100,
                face: results[i].face,
                dde_b: results[i].dde_b,
                dde_s: results[i].dde_s,
                mainforce: results[i].mainforce,
                state: results[i].state,
                per: results[i].per / 10000,
                startprice: results[i].startprice / 100
            });
        }

    return items;
}

suporter.prototype.getfaces = function (item, callback) {

    module.exports.getConnction(function (err, conn) {
        var str = "SELECT * FROM codeface where ";
        var filter = "";

        var tempno;

        if (item.no) {
            var tempno = Number(item.no);
            if (tempno < 1000000)tempno += 1000000;
            if (filter.length > 0)filter += " and "
            filter += "_no=" + tempno;
        }

        if (item.date) {
            if (filter.length > 0)filter += " and "
            filter += "_date='" + item.date + "';";
        }

        if (item.start) {
            if (filter.length > 0)filter += "and "
            filter += "_date>'" + item.start + "';";
        }

        if (item.end) {
            if (filter.length > 0)filter += " and "
            filter += "_date<'" + item.end + "';";
        }

        str += filter;
        conn.query(str, function (err, results) {
            var items = module.exports.convertface(results);
            callback(err, items);
        })
    });
}

suporter.prototype.getfacesbysql = function (sql, callback) {
    module.exports.getConnction(function (err, conn) {
        conn.query(sql, function (err, results) {
            var items = module.exports.convertface(results);
            callback(err, items);
        })
    });
}


suporter.prototype.getcodeface = function (code, date, callback) {
    module.exports.getfaces({no: code, date: date}, function (err, items) {
        if (items && items.length > 0)
            callback(0, items[0]);
        else
            callback(0, null);
    })
}

suporter.prototype.getValueSql = function (o, str, _default) {
    if (_default == undefined)_default = 0;
    var temp = o[str];
    if (temp == undefined)
        return _default;
    return temp;
}

suporter.prototype.savecodefaces = function (items, callback) {
    if (!(items instanceof Array)) items = [items];
    var lists = tool.getSpiedList(items, 300);

    module.exports.getConnction(function (err, conn) {
        async.mapLimit(lists, 1, function (list, mapcallback) {
            var str = "INSERT INTO codeface(_no,_date,_min,_max,ud,lastprice,startprice,face,dde,dde_b,dde_s,mainforce,state)" +
                " VALUES";
            for (var i in list) {
                item = list[i];
                var no = Number(item.no);
                if (no < 1000000)no += 1000000;
                str += '(' +
                    no + ',' +
                    '"' + item.date + '"' + ',' +
                    (module.exports.getValueSql(item, 'min') * 100) + ',' +
                    (module.exports.getValueSql(item, 'max') * 100) + ',' +
                    (module.exports.getValueSql(item, 'ud') * 100) + ',' +
                    (module.exports.getValueSql(item, 'lastprice') * 100) + ',' +
                    (module.exports.getValueSql(item, 'startprice') * 100) + ',' +
                    module.exports.getValueSql(item, 'face') + ',' +
                    module.exports.getValueSql(item, 'dde') + ',' +
                    module.exports.getValueSql(item, 'dde_b') + ',' +
                    module.exports.getValueSql(item, 'dde_s') + ',' +
                    module.exports.getValueSql(item, 'mainforce') + ',' +
                    module.exports.getValueSql(item, 'state') + ')';
                if (i == list.length - 1) str += ';';
                else str += ',';
            }
            conn.query(str, function (err, result) {
                mapcallback(err, str);
            })
        }, function (err, result) {
            if (callback)callback(err, result);
        });
    });

}

suporter.prototype.updatacodeface = function (item, callback) {
    console.log(Date.now() )
    module.exports.getConnction(function (err, conn) {
        var fun = function (option, cb) {
            console.log(Date.now() )
            var items=option;
            if(option.items)items=option.items;
            if(!util.isArray(items)) items=[items];

            var faids = ["state", "_min", "_max", "face", "dde", "dde_s", "dde_b", "ud", "mainforce", "lastprice", "per", "startprice"];
            var lists = tool.getSpiedList(items, 200);
            async.mapLimit(lists,1,function(list,aycb){
                var str="";
                for(var i in list){
                    var item=list[i];
                    var temp="";
                    for(var j in faids){
                        if (item[faids[j]] == undefined) continue;
                        if (temp.length > 0) temp += ",";
                        var value = item[faids[j]];
                        if (faids[j] == "lastprice" || faids[j] == "ud" || faids[j] == "_min" || faids[j] == "_max" || faids[j] == "startprice") value =Math.floor( value * 100);
                        if (faids[j] == "per") value = value * 10000;
                        temp += faids[j] + "=" + value;
                    }
                    var tmno = Number(item.no);
                    if (tmno < 1000000)tmno += 1000000;
                    temp = "update codeface set " + temp + " where _no=" + tmno + " and _date='" + item.date + "'; ";
                    str+=temp;
                }
                conn.query(str, function (err, r) {
                    console.log(Date.now() )
                    if (aycb)aycb(err, r)
                })
            },function(err,result){
                callback(err,result);
            })
        }

        if (item.uncheck) return fun(item, callback)
        module.exports.getcodeface(item.no, item.date, function (err, face) {
            if (face) {
                fun(item, callback)
            }
            else {
                module.exports.savecodefaces(item, callback);
            }
        })
    });

}

var current = module.exports = new suporter();

//var upDateItem = function (item, callback) {
//
//    if (item.no == global.shcode) return callback(0, 0);
//
//
//    current.getValueByDayNo(item, function (err, valueitems) {
//
//        // console.log(item)
//        if (valueitems.length == 0) {
//            console.log("数据丢失:" + item.date + " " + item.no)
//        }
//        //  return callback(0,0);
//
//        item.startprice = valueitems[0].price;
//        //current.updatacodeface(item, function (err, result) {
//        //    var log = "更新成功";
//        //    if (err) log = "更新失败";
//        //    console.log(log + ":" + item.date + " " + item.no)
//        //    callback(err, result);
//        //
//        //})
//
//        callback(err, item);
//    })
//}

//current.getfaces({no: global.shcode}, function (err, dates) {
//    var time= Date.now();
//    async.mapLimit(dates, 2, function (d, dc) {
//
//       // if (new Date(d.date) < new Date("2017-1-1 0:00:00")) return dc(0, 1);
//        if( new Date (d.date)<new Date("2017-04-20 22:00:00") ) return dc(0,1);
//        module.exports.getConnction(function (err, conn) {
//            conn.query('CALL codeface_update("' + d.date+'");', function(err, rows) {
//                var t=Date.now();
//                console.log((t-time)/1000+"  更新完成" + ":" + d.date)
//                time=t;
//                dc(err,rows);
//            });
//        });
//
//
//        return;
//
//
//        current.getfaces({date: d.date}, function (err, codeitem) {
//            async.mapLimit(codeitem, 4, upDateItem, function (err, result) {
//        //        console.log(Date.now() + ":" + d.date)
//
//
//                current.updatacodeface({items:result,uncheck:true},function(err,r){
//                     t=Date.now();
//                    console.log("更新完成" + ":" + d.date)
//                    console.log("耗时" + ":" + (t-time)/1000);
//                    time=t;
//                    dc(err,  result);
//                })
//            })
//        })
//    }, function (err, result) {
//        console.log("全部更新完成")
//        console.log("全部更新完成")
//        console.log("全部更新完成")
//        console.log("全部更新完成")
//        console.log("全部更新完成")
//        console.log("全部更新完成")
//    })
//})