<template>
  <div class="modal fade" style="display:block;" v-show="alertv">
    <div class="modal-dialog animated tada" id="alertvBox" style="position:absolute;" v-bind:style="{ width: width + 'px', left: left + 'px', top: '20%' }">
      <div class="modal-content">
        <div class="modal-header" id="alertvBar" style="cursor: move;">
          <h4 class="modal-title">系统提示</h4>
        </div>
        <div class="modal-body">{{ alertvText }}</div>
        <div class="modal-footer">
          <button type="button" class="btn btn-primary" @click="btnOk">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    alertv () {
      return this.$store.state.alertv
    },
    alertvText () {
      return this.$store.state.alertv_text
    },
    width () {
      let w = this.$store.state.alertv_text_length * 20
      w = w < 150 ? 150 : w
      return w
    },
    alertvCallback () {
      return this.$store.state.alertv_callback
    },
    left () {
      return document.body.clientWidth / 2 - this.width / 2
    }
  },
  methods: {
    btnOk () {
      this.$store.commit('hideAlertv')
      if (this.alertvCallback) {
        this.alertvCallback()
      }
    }
  },
  updated: function () {
    this.$drag.startDrag(document.getElementById('alertvBar'), document.getElementById('alertvBox'))
  }
}
</script>
