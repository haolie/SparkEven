
var hideFun=function(){
    var sider= document.querySelector(".fl.sidebar")
    if(sider){
        sider.style.display="none"
    }else setTimeout(hideFun,1000)  
}

hideFun()