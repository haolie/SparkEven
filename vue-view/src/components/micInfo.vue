<template>
  <div>
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="信息管理" name="first">

      </el-tab-pane>
      <el-tab-pane label="数据分析" name="second">

      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
  var _this;

  export default {
    components:{},
    data:function () {
      return{

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
