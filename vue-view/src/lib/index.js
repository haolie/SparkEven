let Libs = () => {}

Libs.prototype = {
  /**
   * 获取url中的param
   * name：key
   * return：value
   */
  getUrlParam (name) {
    let url = window.location
    let reg = new RegExp('(^|&)' + name + '=([^&]*)(&|$)')
    let r = url.search.substr(1).match(reg)
    return (r === null || r.length < 2) ? null : r[2]
  },

  pageList (list,size,index,filterFun) {
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
}

export default {
  install (Vue, name = '$lib') {
    Object.defineProperty(Vue.prototype, name, { value: new Libs() })
  }
}
