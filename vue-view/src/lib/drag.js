var paramsDrag = {
  top: 0,
  left: 0,
  currentX: 0,
  currentY: 0,
  flag: false
}

/**
obj.currentStyle[attr]
getComputedStyle(obj,false)[attr] 获取DOM非行间样式
**/
var getCss = function (o, key) {
  return o.currentStyle ? o.currentStyle[key] : document.defaultView.getComputedStyle(o, false)[key]
}

function startDrag_ (bar, target, callback) {
  if (getCss(target, 'left') !== 'auto') {
    paramsDrag.left = getCss(target, 'left')
  }
  if (getCss(target, 'top') !== 'auto') {
    paramsDrag.top = getCss(target, 'top')
  }
  bar.onmousedown = function (event) {
    paramsDrag.flag = true
    if (!event) {
      event = window.even
      bar.onselectstart = function () { // IE和Chrome适用，防止内容被选中默认是true
        return false
      }
    }
    var e = event
    paramsDrag.currentX = e.clientX
    paramsDrag.currentY = e.clientY
  }
  document.onmouseup = function () {
    paramsDrag.flag = false
    if (getCss(target, 'left') !== 'auto') {
      paramsDrag.left = getCss(target, 'left')
    }
    if (getCss(target, 'top') !== 'auto') {
      paramsDrag.top = getCss(target, 'top')
    }
  }
  document.onmousemove = function (event) {
    // var e = event ? event : window.event
    var e = {}
    if (event) {
      e = event
    } else {
      e = window.event
    }
    if (paramsDrag.flag) {
      var nowX = e.clientX
      var nowY = e.clientY
      var disX = nowX - paramsDrag.currentX
      var disY = nowY - paramsDrag.currentY
      if (nowY < 20) {
        return
      }
      target.style.left = parseInt(paramsDrag.left) + disX + 'px'
      target.style.top = parseInt(paramsDrag.top) + disY + 'px'
    }
    if (callback === 'function') {
      let left = parseInt(paramsDrag.left) + disX
      let top = parseInt(paramsDrag.top) + disY
      callback(left, top)
    }
  }
}

let Drag = () => {}

Drag.prototype = {
  startDrag (bar, target, callback) {
    startDrag_(bar, target, callback)
  }
}

export default {
  install (Vue, name = '$drag') {
    Object.defineProperty(Vue.prototype, name, { value: new Drag() })
  }
}
