export default {
  state: {
    alertv: false,
    alertv_text: '',
    alertv_text_length: 0,
    alertv_callback: function () {}
  },
  mutations: {
    showAlertv (state, o) {
      state.alertv_text = o.text
      state.alertv_text_length = o.text.length
      state.alertv_callback = o.callback || function () {}
      state.alertv = true
    },
    hideAlertv (state) {
      state.alertv = false
    }
  },
  actions: {
  }
}
