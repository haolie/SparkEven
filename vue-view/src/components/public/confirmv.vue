<template>
  <div class="modal fade" style="display:block;" v-show="confirmv">
    <div class="modal-dialog animated flash" id="confirmvBox" style="position:absolute;" v-bind:style="{ width: width + 'px', left: left + 'px', top: '20%' }">
      <div class="modal-content">
        <div class="modal-header" id="confirmvBar" style="cursor: move;">
          <h4 class="modal-title">系统提示</h4>
        </div>
        <div class="modal-body">{{ confirmvText }}</div>
        <div class="modal-footer">
          <button type="button" class="btn btn-primary" @click="btnOk">确定</button>
          <button type="button" class="btn btn-minor" @click="btnNO">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    confirmv () {
      return this.$store.state.confirmv
    },
    confirmvText () {
      return this.$store.state.confirmv_text
    },
    width () {
      let w = this.$store.state.confirmv_text_length * 20
      w = w < 150 ? 150 : w
      return w
    },
    confirmvCallback () {
      return this.$store.state.confirmv_callback
    },
    left () {
      return document.body.clientWidth / 2 - this.width / 2
    }
  },
  methods: {
    btnOk () {
      this.$store.commit('hideConfirmv')
      if (this.confirmvCallback) {
        this.confirmvCallback()
      }
    },
    btnNO () {
      this.$store.commit('hideConfirmv')
    }
  },
  updated: function () {
    this.$drag.startDrag(document.getElementById('confirmvBar'), document.getElementById('confirmvBox'))
  }
}
</script>
