export default {
  state: {
    loadDom: false,
    loadingNum: 0
  },
  mutations: {
    showLoading (state) {
      state.loadDom = true
    },
    hideLoading (state) {
      state.loadDom = false
    }
  },
  actions: {
  }
}
