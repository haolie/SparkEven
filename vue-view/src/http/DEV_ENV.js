let DEV_ENV = 0
let host = window.location.hostname
let port = window.location.port
if (host.indexOf('localhost') > -1 || port === '8106') {
  DEV_ENV = 1
}

export default DEV_ENV
