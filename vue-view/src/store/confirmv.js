export default {
  state: {
    confirmv: false,
    confirmv_text: '',
    confirmv_text_length: 0,
    confirmv_callback: function () {}
  },
  mutations: {
    showConfirmv (state, o) {
      state.confirmv_text = o.text
      state.confirmv_text_length = o.text.length
      state.confirmv_callback = o.callback || function () {}
      state.confirmv = true
    },
    hideConfirmv (state) {
      state.confirmv = false
    }
  },
  actions: {
  }
}
