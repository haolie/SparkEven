<template>
  <div style=" width: 100%;height:800px">
    <div class="topButtons">
      <ul>
        <li><a @click="state=0"><i class="icon-bd-o"></i></a></li>
        <li><a @click="state=1"><i class="icon-bd-o"></i></a></li>
      </ul>


    </div>
    <chart_line :dataItems="data_line" :face="face" :class="{hideitem:state==1}" ></chart_line>
    <kchart :dataItems="data_kLine" :no="$route.params.no"  :class="{hideitem:state==0}"></kchart>
  </div>
</template>

<script>
  import kchart from '@/components/public/charts_k'
  import chart_line from '@/components/public/chart_line'

  var _this;

  export default {
    components:{kchart,chart_line},
    data:function () {
      return{
        state:0,
        face:null,
        data_kLine:null,
        data_line:null,

      }
    },
    computed:{

    },
    methods:{
       getKlineData(){
         this.$http.request('DV_KLINE',{no:this.$route.params.no,date:this.$route.query.date,count:this.$route.query.count?this.$route.query.count:90},function (data) {
           _this.data_kLine=data.content
         })
       },
      getLineDate(){
        this.$http.request('DV_DATELINE',{no:this.face.no,date:this.face.date},function (data) {
          _this.data_line=data
        })
      }
    },
    watch:{},

    created:function () {
      _this=this;
    },
    mounted:function () {

      var params={
        pageIndex:1,
        pageSize:10,
        date:this.$route.query.date,
        no:this.$route.params.no
      }
      this.$http.request('DV_CODEPAGE',params,function (data) {
        if(data.content.length){
          _this.face=data.content[0]
          _this.getKlineData()
          _this.getLineDate()

        }
      })


    }
  }
</script>
