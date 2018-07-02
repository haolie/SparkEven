<template>
  <div>
   <ul class="menu-contaier">
     <li class="menu-item" v-for="item in menu"><a v-bind:href="'#/'+item.key" >{{item.name}}</a></li>
   </ul>
  </div>
</template>

<script>
  var _this,interval=5000;

  export default {
    components:{},
    data:function () {
      return{
        menu:[
          {name:"概览",key:"overView",children:{}},
          {name:"采集",key:"gather",children:{}},
          {name:"数据查询",key:"dataview",children:{}},
          {name:"微讯",key:"micInfo",children:{}},
        ]
      }
    },
    methods:{

    },
    wathc:{},

    created:function () {
      _this=this;
    },
    mounted:function () {

    }
  }
</script>
