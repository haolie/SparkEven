
var hideFun=function(){
    var sider= document.querySelector(".fl.sidebar")
    if(sider){
        sider.style.display="none"
    }else setTimeout(hideFun,1000)  
}

hideFun()


var ConsoleObj=function(){

}

ConsoleObj.prototype.downlink=function(){
 var linklist= document.querySelectorAll(".list a")
 var o={}
 var list=[]
 for(var i=0;i<linklist.length;i++){
   var temp= linklist[i].getAttribute("href")
   if (temp.indexOf("#")){
     temp=temp.split("#")[0]
   }

   if(o[temp]) continue
   o[temp]=1
   list.push({
    savepath: "E:/web",
    path: "https://openlayers.org"+temp.substring(1),
    filename: temp.substring(1),
   })
 }
 console.log(list)

 var xmlhttp=new XMLHttpRequest();
 xmlhttp.open("POST","http://localhost:9520/down/downfiles",true);
xmlhttp.setRequestHeader("Content-type","application/json");
xmlhttp.send(list)
}

ConsoleObj.prototype.test=function(){
    console.log("aready")
}

var CObj=new ConsoleObj()



console.log(3333222)