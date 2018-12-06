import DEV_ENV from '@/http/DEV_ENV'
// 'T_S_USER':DEV_ENV ?'sys/userMessage':  'user/menu/getUserInfo',
let appApi = (c) => {
  let a = {
    //采集

    'G_DATEITEMLIST':'meeter/dates',
    'G_DATECODEITEM':'meeter/dateItem',

    "DV_CODEPAGE":"dataview/codes_page",

    "M_UPLOADFILE":"micinfo/fileUpload",
    "M_ITEMSLIST_PAGE":"micinfo/itemList_page",
    "M_ITEM":"micinfo/item",

    'DV_KLINE':'dataview/facesInDate',
    'DV_DATELINE':"dataview/lineData"
  }
  return '/express/' + a[c]
}

export default appApi
