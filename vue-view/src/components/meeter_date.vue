<template>
  <div>
    <el-row>
      <el-col :span="12"><div class="grid-content bg-purple">
        1</div></el-col>
      <el-col :span="12"><div class="grid-content bg-purple-light">2</div></el-col>
    </el-row>
  </div>
</template>

<script>
  var _this;

  export default {
    data:function () {
      return{
        dateItems:[]
      }
    },
    created:function () {
      _this=this;
    },
    mounted:function () {

    }
  }
</script>
