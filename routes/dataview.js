var express = require('express');
var router = express.Router();
var dbsurport= require('../fvvm/MYSQLDBSuport');
var tools= require('../common/tools');

dbsurport.initCodesObj()

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.send('respond with a resource');
});

/* 查询日期状态 */
router.get('/codes_page', function(req, res, next) {
    var filter=[];
    if(req.query.date&&req.query.date.length)filter.push("_date='"+req.query.date+"'");
    if(req.query.no&&req.query.no.length)filter.push(" no_id=" + dbsurport.getIdByNo(req.query.no));
    if(filter.length)
    filter=" where "+filter.join(" and ")
    var start=(req.query.pageIndex-1)*req.query.pageSize;


    var mysql ="select *,(select count(id) from codeface" +
        filter +
        ") as total from codeface"
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
            result.data.content=dbsurport.convertface(items);
            result.data.totalRecord=items.length?items[0].total:0;
        }
        res.send(result);
    })

});


module.exports = router;
