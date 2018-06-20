<template>
  <div class="alertLoading" id="loading" v-show="loadDom" style="display: block;">
    <div class="layerbg"></div>
    <div class="layerbox">
      <img src="static/assets/images/loading.gif">
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    loadDom () {
      return this.$store.state.loadDom
    }
  }
}
</script>
