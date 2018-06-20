var express = require('express');
var router = express.Router();
var meeter= require('../fvvm/DataMeeter.js');
var tools= require('../common/tools');

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.send('respond with a resource');
});

/* 查询日期状态 */
router.get('/dates', function(req, res, next) {
    var list=[];
    for(var d in  meeter.dataContext.dates){
        list.push({date:d,count:meeter.dataContext.dates[d].count,state:meeter.dataContext.dates[d].state});
    }
    list=tools.pageList(list,req.query.pageSize,req.query.pageIndex);

    res.send(tools.createBaseRespon( list) );
});

/* 查询某日的条目信息 */
router.get('/dateItem', function(req, res, next) {
    var date=req.query.date;
    var list=[];
    if(meeter.dataContext.dates[date]){
        for(var i=0;i<meeter.dataContext.dates[date].items.length;i++){
          list.push(meeter.dataContext.dates[date].items[i]);
        }
    }

    res.send(tools.createBaseRespon( list) );
});



module.exports = router;
