<style>
  @import "../static/assets/css/elm.css";
  @import "../static/assets/css/style.css";
  @import "../static/assets/css/analysis.css";
</style>

<template>


  <div id="app">
    <el-container >
      <el-aside width="200px"><tilemenu></tilemenu></el-aside>
      <el-container>
        <!--<el-header>Header</el-header>-->
        <el-main><router-view/></el-main>
        <!--<el-footer>Footer</el-footer>-->
      </el-container>
    </el-container>

  </div>
</template>

<script>
  import tilemenu from '@/components/menu'
  import alertv from '@/components/public/alertv'
  import confirmv from '@/components/public/confirmv'
  import loading from '@/components/public/loading'

export default {
  name: 'App',
  components: {
    tilemenu,alertv,confirmv,loading
  }
}
</script>

