var express = require('express');
var router = express.Router();
var dbsurport= require('../fvvm/MYSQLDBSuport');
var tools= require('../common/tools');
var lcsv = require('../fvvm/localCSV.js');
var fs= require('fs');

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.send('respond with a resource');
});

/* 批量上传 */
router.get('/fileUpload', function(req, res, next) {

    if(!req.query.file||!req.query.file.length) {
        res.send(tools.createBaseRespon([],801,"参数出错"))
    }
    else {
        fs.exists(req.query.file,function(exist){
            if(!exist){
                res.send(tools.createBaseRespon([],801,"文件未找到"))
                return
            }
            lcsv.FileToJson(req.query.file, function (err, result) {
                console.log(result)
                console.log(result)
                console.log(result)
                console.log(req.query.file)
                console.log(req.query.file)

                if(err){
                    res.send(tools.createBaseRespon([],801,"文件解析出错"))
                    return
                }
                var items=[]
                result.forEach(function (item,i) {
                    items.push({no:item[0],date:item[1]})
                })

                addItems(items,function (err,result) {
                    if(err){
                        res.send(tools.createBaseRespon([],801,"数据保存失败"))
                        return
                    }
                    else
                        res.send(tools.createBaseRespon(items))
                })

            })

        })


    }


});

/* 分页列表获取. */
router.get('/itemList_page', function(req, res, next) {
    var filter=[];
    if(req.query.date&&req.query.date.length)filter.push("_date='"+req.query.date+"'");
    if(req.query.no&&req.query.no.length)filter.push(" no_id=" + dbsurport.getIdByNo(req.query.no));
    if(filter.length)
        filter=" where "+filter.join(" and ")
    var start=(req.query.pageIndex-1)*req.query.pageSize;


    var mysql ="select tbl_codes._no,micinfo.infoDate,(select count(id) from micinfo inner join tbl_codes on micinfo.no_id=tbl_codes.id" +
        filter +
        ") as total from micinfo inner join tbl_codes on micinfo.no_id=tbl_codes.id"
        +filter+
        " order by id limit " +
        start +
        "," +
        req.query.pageSize +
        ";";

    console.log(mysql)

    dbsurport.transQeurysql(mysql,function (err,items) {
        var result= tools.createBaseRespon();
        if(err){
            result.code=701;
            result.msg="数据获取失败";
            result.data.content=0;
            result.data.totalRecord=0;
        }
        else {
            var list=[];
            for(var i=0;i<items.length;i++){
                list.push({no:items[i]._no,date:items[i].infoDate})
            }
            result.data.content=list
            result.data.totalRecord=items.length?items[0].total:0;
        }
        res.send(result);
    })

});

/* 删除. */
router.delete('/item', function(req, res, next) {

    if(!req.query.date||!req.query.no){
        res.send(tools.createBaseRespon([],801,"参数丢失"));
        return
    }

    var sql="delete a.* from micinfo a left join tbl_codes b on a.no_id=b.id where b.no=" +
        req.query.no +
        " and a.infoDate=" +
        req.query.date +
        ";"

    console.log(sql)
    dbsurport.transQeurysql(sql,function (err,result) {
        if(err){
           res.send( tools.createBaseRespon(0,200,"数据库删除失败"))
        }
        else  res.send( tools.createBaseRespon(1));
    })

});

/* 添加. */
router.post('/item', function(req, res, next) {

    if(!req.body.no||!req.body.date){
        res.send(tools.createBaseRespon(0,801,"参数丢失"));
        return;
    }

    addItems([res.body],function () {
        if(err){
            res.send( tools.createBaseRespon(0,200,"数据库添加失败"))
        }
        else  res.send( tools.createBaseRespon(1));
    })


});

function addItems(items,callback) {
    var sql="replace INTO micinfo(no_id,infoDate) VALUES"
    items.forEach(function (item,i) {
        sql+='(' +
            dbsurport.getIdByNo( item.no) +
            ',"' +
            item.date +
            '")'
        if(i==items.length-1)
            sql+=";"
        else sql+=","
    })

    console.log(sql)
    dbsurport.transQeurysql(sql,callback)
}




module.exports = router;
