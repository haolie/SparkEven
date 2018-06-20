import axios from 'axios'
import store from '@/store'
import DEV_ENV from '@/http/DEV_ENV'
import appApi from '@/appApi'

let Http = () => {}

Http.prototype = {
  get (o) {
    o.type = 'get'
    return this.xhr(o)
  },
  post (o) {
    o.type = DEV_ENV ? 'get' : 'post'
    return this.xhr(o)
  },
  put (o) {
    o.type = DEV_ENV ? 'get' : 'put'
    return this.xhr(o)
  },
  delete (o) {
    o.type = DEV_ENV ? 'get' : 'delete'
    return this.xhr(o)
  },
  getUrl (o) {
    let uri = appApi(o.api)
    if (uri === '') {
      store.commit('showAlertv', {text: '请填写api'})
    }
    if (o.type !== 'post' && o.params.id && !DEV_ENV) {
      uri = uri + '/' + o.params.id
    }
    if (o.type === 'get' && Object.keys(o.params).length > 0) {
      uri = uri + '?' + this.joinP(o.params)
    }
    return uri
  },
  joinP (o) {
    let x = []
    for (let i in o) {
      if (i !== 'id') {
        x.push(`${i}=${o[i]}`)
      }
    }
    return encodeURI(x.join('&'))
  },
  xhr (o) {
    return new Promise((resolve, reject) => {
      let config = {
        headers: o.headers ? o.headers : {
          'Content-type': 'application/json'
        },
        timeout: o.timeout ? o.timeout : 30000
      }
      if (!o.loading) {
        store.state.loadingNum ++
        store.commit('showLoading')
      }
      let instance = axios.create(config)
      let p = o.type === 'get' ? '' : o.params
      instance[o.type]( o.url?o.url:this.getUrl(o), p).then((res) => {
        if (!o.loading) store.state.loadingNum --
        if (store.state.loadingNum === 0) {
          store.commit('hideLoading')
        }
        if(o.url){   //http 直接请求
          resolve(res.data)
        }
        else
        if (Number(res.data.status.code) === 200) {
          resolve(res.data.data)
        } else {
          let msg = res.data.status.message ? res.data.status.message : '对不起，服务器接口出错！请联系技术人员！'
          reject(msg)
        }
      }).catch((err) => {
        store.commit('hideLoading')
        if (err.message.indexOf('timeout') > -1) {
          reject(new Error('请求超时,请重新请求！'))
        } else {
          reject(err.message)
        }
      })
    })
  }
}

export default {
  install (Vue, name = '$http') {
    Object.defineProperty(Vue.prototype, name, {value: new Http()})
  }
}
