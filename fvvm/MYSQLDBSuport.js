var mysql = require('mysql');
var http = require('http');
var async = require("async");
var url = require("url");
var tool = require('./tools.js');
var dateutil = require('date-utils');
//var co=require("co");
var codeface = "codeface";
var time_price = "time_price";
var process = require('process');
var util = require('util');
global.shcode = 912261;
var maxCodeId=-1;
var maxFaceId=-1;

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
            database: 'test',
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
            items.forEach(function (item) {
                maxCodeId=Math.max(maxCodeId,item.id);
                current.CodesObj[item._no]=item.id;
            })

            conn.query("select max(id) as id from codeface",function (err,items) {
                maxFaceId=items[0].id;if(!maxFaceId)maxFaceId=0;
                if(callback) callback(err,items);
            })
        })
    })
}

suporter.prototype.getNoById=function (id) {
    for(var i in current.CodesObj){
        if(current.CodesObj[i]==id)return i;
    }

    return "";
}

suporter.prototype.getIdByNo=function (no) {
    if(no==global.shcode){
        console.log(global.shcode);
    }
    no=Number(no);
    if(no<1000000)no+=1000000;
    if(current.CodesObj[no]==undefined)
        return -1;
    return current.CodesObj[no];
}




suporter.prototype.getInsertStr = function (item,id) {
    var seconds=new Date(item.time*1000);
    seconds=((seconds.getHours()-9)*60+seconds.getMinutes())*60+seconds.getSeconds();

  //  console.log(item.no)
    return '(' +
        id + ',' +
        seconds+ ',' +
        (Number(item.price.toString().replace(".",""))) + ',' +
        item.trade_type + ',' +
        (item.volume) +
        ')';
}

// suporter.prototype.deletebyno = function (no, callback) {
//     module.exports.getConnction(function (err, conn) {
//         var str = 'delete from time_price where no=' + no;
//         conn.query(str, function (err, result) {
//             callback(err, result);
//         })
//     })
// }

suporter.prototype.saveTimePrice = function (timevalues, allcallback) {
    if (timevalues == null || timevalues.length == 0) {

        if (allcallback)allcallback(1, "");
        return;
    }

    var start=timevalues[0].price;
    timevalues.forEach(function (item) {
        item.price=(item.price-start).toFixed(2);
    })

    var lists = tool.getSpiedList(timevalues, 200);

    module.exports.getConnction(function (err, conn) {
        if (err == null) {

            var date=new Date(timevalues[0].time*1000).toFormat("YYYY-MM-DD");
            current.getFaceId(new Date(timevalues[0].time*1000),timevalues[0].no,function (err,id) {
                async.mapLimit(lists, 1, function (items, callback) {
                    var insert = 'replace INTO time_price(face_id,time,price,trade_type,volume) VALUES';
                    for (var i in items) {
                        insert += module.exports.getInsertStr(items[i],id);
                        if (i == items.length - 1) insert += ';';
                        else insert += ',';
                    }
                    conn.query(insert, function (err, result) {
                        callback(err, insert);
                    })
                }, function (err, result) {
                    if (allcallback)allcallback(err, "");
                })
            })

        }
        else {
            if (allcallback)allcallback(null, "");
            console.log("数据库连接失败");
        }
    });
};

 suporter.prototype.getFaceId=function(date,no,callback){
    if(!tool.isDate(date))date=new Date(date);
    date= date.toFormat("YYYY-MM-DD");

    current.getConnction(function (err,conn) {
        var sql="select b.id from tbl_codes as a right join codeface as b on a.id=b.no_id where a._no=" +current.convertNo(no) +" and b._date='" +date + "'";
        conn.query(sql,function (err,ids) {
            if(!ids.length||ids[0]==undefined){

                console.log(no+"--"+date);
                console.log(no+"--"+date);
                console.log(no+"--"+date);
                console.log(no+"--"+date);
                console.log(no+"--"+date);
                console.log(no+"--"+date);
                console.log(no+"--"+date);
                console.log(no+"--"+date);
                console.log(no+"--"+date);
                console.log(no+"--"+date);
                console.log(sql);
            }
            callback(err,ids[0].id);
        })
    })

}



suporter.prototype.getValueByDayNo = function (item, callback) {
    var date = item.date;
    if(!tool.isDate(date))
        date=new Date(date);

    current.getFaceId(date,item.no,function (err,id) {
        current.getConnction(function (err,conn) {
            conn.query('select * from time_price where face_id="' +id +'" order by time',function (err,result) {
                var items = [];
                var no=item.no>1000000?item.no-1000000:item.no;
                if (result && result.length > 0) {
                    for (var i in result) {
                        items.push({
                            no: no,
                            time: date.addSeconds(result[i].time),
                            price: result[i].price / 100,
                            trade_type: result[i].trade_type,
                            volume: result[i].volume * 100
                        })
                    }
                }

                callback(err, items);
            })
        })
    })

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
                no: current.getNoById(results[i].no_id).toString().substring(1),
                date:  new Date(results[i]._date).toFormat("YYYY-MM-DD"),
                min: results[i]._min / 100,
                max: results[i]._max / 100,
                ud: results[i]._change / 100,
                lastprice: results[i].lastprice / 100,
                face: results[i].face,
                dde_b: results[i].dde_b,
                dde_s: results[i].dde_s,
                state: results[i].state,
                per: results[i].per / 10000,
                startprice: results[i].startprice / 100
            });
        }

    return items;
}

suporter.prototype.getfaces = function (item, callback) {

    module.exports.getConnction(function (err, conn) {
        var str = "SELECT * FROM codeface";
        var filter = "",pagestr;

        var tempno;

        if (item.no) {
            var tempno = Number(item.no);
            if (filter.length > 0)filter += " and"
            filter += " no_id=" + current.getIdByNo(tempno);
        }

        if (item.date) {
            if (filter.length > 0)filter += " and"
            filter += " _date='" + item.date + "'";
        }

        if (item.start) {
            if (filter.length > 0)filter += " and"
            filter += " _date>'" + item.start + "'";
        }

        if (item.end) {
            if (filter.length > 0)filter += " and"
            filter += " _date<'" + item.end + "'";
        }

        if(filter.length)
            filter=" where"+filter


        str += filter+";";
        conn.query(str, function (err, results) {
            var items = module.exports.convertface(results);
            callback(err, items);
        })
    });
}


suporter.prototype.transQeurysql = function (sql, callback) {
    module.exports.getConnction(function (err, conn) {
        conn.query(sql, callback)
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

    current.checkNoSave(items,function(){

        if (!(items instanceof Array)) items = [items];
        var lists = tool.getSpiedList(items, 300);
        module.exports.getConnction(function (err, conn) {
            conn.query('delete from codeface where _date="'+items[0].date +'";',function (err,r) {
                async.mapLimit(lists, 1, function (list, mapcallback) {
                    var str = "INSERT INTO codeface(id,_date,no_id,_min,_max,_change,lastprice,startprice,face,dde,dde_b,dde_s,state)" +
                        " VALUES";
                    for (var i in list) {
                        var item = list[i];
                        str += '(' +
                            (++maxFaceId) +',' +
                            '"' + item.date + '"' + ',' +
                            current.getIdByNo(item.no)+ ',' +
                            (module.exports.getValueSql(item, 'min') * 100) + ',' +
                            (module.exports.getValueSql(item, 'max') * 100) + ',' +
                            (module.exports.getValueSql(item, 'ud') * 100) + ',' +
                            (module.exports.getValueSql(item, 'lastprice') * 100) + ',' +
                            (module.exports.getValueSql(item, 'startprice') * 100) + ',' +
                            module.exports.getValueSql(item, 'face') + ',' +
                            module.exports.getValueSql(item, 'dde') + ',' +
                            module.exports.getValueSql(item, 'dde_b') + ',' +
                            module.exports.getValueSql(item, 'dde_s') + ',' +
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
            })
        });
    });
}

suporter.prototype.convertNo=function (no) {
    no=Number(no);
    if(no<1000000)no+=1000000
    return no;
}

suporter.prototype.checkNoSave=function (items,callback) {
    var codes=[];
    items.forEach(function (item) {
        codes.push(current.convertNo(item.no));
    })
    current.saveNos(codes,callback)
}

suporter.prototype.saveNos=function (codes,callback) {
    var temps=[];
    codes.forEach(function (c) {
       if(current.getIdByNo(c)>=0)return;
        maxCodeId++;
        temps.push({id:maxCodeId,no:c});
    })

    if(temps.length==0){
        callback(0,0)
        return;
    }

    module.exports.getConnction(function (err,conn) {
        temps=tool.getSpiedList(temps,300);
        async.mapLimit(temps,1,function (item,ab) {
          var str="INSERT INTO tbl_codes(id,_no)" +
              " VALUES";

          item.forEach(function (n,i) {
              str+='(' +
                  n.id+","+
                  n.no+
                  ')';
              if(i==item.length-1) str+";";
              else str+=",";
          })

            conn.query(str, function (err, result) {
                item.forEach(function (t) {
                    current.CodesObj[t.no]=t.id;
                })
                ab(err, str);
            })

        },function (err,r) {
            callback(0,0);
        })
    })
}

suporter.prototype.updatacodeface = function (item, callback) {
    //console.log(Date.now() )
    module.exports.getConnction(function (err, conn) {
        var fun = function (option, cb) {
           // console.log(Date.now() )
            var items=option;
            if(option.items)items=option.items;
            if(!util.isArray(items)) items=[items];

            var faids = ["state", "_min", "_max", "face", "dde", "dde_s", "dde_b", "ud|_change","lastprice", "per", "startprice"];
            var lists = tool.getSpiedList(items, 200);
            async.mapLimit(lists,1,function(list,aycb){
                var str="";
                for(var i in list){
                    var item=list[i];
                    var temp="";
                    for(var j in faids){
                        var sf=faids[j];fb=sf;
                        if(sf.indexOf("|")>0){
                            sf=sf.split('|');
                            fb=sf[1];
                            sf=sf[0];
                        }

                        if (item[sf] == undefined) continue;
                        if (temp.length > 0) temp += ",";
                        var value = item[sf];
                        if (sf == "lastprice" || sf == "ud" ||sf== "_min" || sf== "_max" || sf == "startprice") value =Math.floor( value * 100);
                        if (sf == "per") value = value * 10000;
                        temp +="c."+ fb+ "=" + value;
                    }
                    var tmno = Number(item.no);
                    if (tmno < 1000000)tmno += 1000000;
                    temp = "update codeface c join tbl_codes t on c.no_id=t.id set " + temp + " where t._no=" + tmno + " and c._date='" + item.date + "'; ";
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
        fun(item, callback)
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