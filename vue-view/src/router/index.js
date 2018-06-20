import Vue from 'vue'
import Router from 'vue-router'
import meeter from '@/components/meeter'
import micInfo from '@/components/micInfo'
import overView from '@/components/overView'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      component: overView
    },
    {
      path: '/overView',
      name: 'overView',
      component: overView
    },
    {
      path: '/gather',
      name: 'gather',
      component: meeter
    },
    {
      path: '/micInfo',
      name: 'micInfo',
      component: micInfo
    }
  ]
})
