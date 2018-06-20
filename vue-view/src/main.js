// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.

import Vue from 'vue'
import App from './App'
import store from '@/store'
import router from './router'
import http from '@/http'
import lib from '@/lib'
import drag from '@/lib/drag'
import dateUtils from 'date-utils';
import alertv from '@/components/public/alertv'
import confirmv from '@/components/public/confirmv'
import loading from '@/components/public/loading'
import ElementUI from 'element-ui'
import echarts from 'echarts'


Vue.config.productionTip = false

Vue.prototype.$echarts = echarts
Vue.use(http)
Vue.use(lib)
Vue.use(drag)
Vue.use(ElementUI);

new Vue({
  el: '#app',
  store,
  router,
  components: { App ,alertv,confirmv,loading},
  template: '<App/>'
})
