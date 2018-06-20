
var Tools=function () {
    
}

Tools.prototype.createBaseRespon=function (data,code,msg) {
    var temp={
        "status": {
            "code": code?code:200,
            "message": msg?msg:"系统请求成功！"
        },
        "data": data?data:{}
    }

    return temp;
}

Tools.prototype.pageList=function (list,size,index,filterFun) {
    if(!list||!size) return  {content:[],totalRecord:0,pageIndex:1};
    var count=Number(list)
    count=count?count:list.length;
    if(!count)return {content:[],totalRecord:0,pageIndex:1};;

   if(!index||index<0)index=1;
    var i=(index-1)*size;
    if(i>=count) return module.exports.pageList(list,size,index-1,filterFun);
    var temp=[];
    while (i<count&&temp.length<size){
        if(filterFun){
            temp.push(filterFun(list.length?list[i]:i))
        }
        else temp.push(list[i])
        i++;
    }

    return  {content:temp,totalRecord:count,pageIndex:index};
}

module.exports=new Tools();
