import DEV_ENV from '@/http/DEV_ENV'
// 'T_S_USER':DEV_ENV ?'sys/userMessage':  'user/menu/getUserInfo',
let appApi = (c) => {
  let a = {
    //采集

    'G_DATEITEMLIST':'meeter/dates',
    'G_DATECODEITEM':'meeter/dateItem',
  }
  return '/express/' + a[c]
}

export default appApi
