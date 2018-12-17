<template>
  <div style=" width: 100%;height:800px">
    <div class="topButtons">
      <ul>
        <li><a @click="state=0"><i class="icon-bd-o"></i></a></li>
        <li><a @click="state=1"><i class="icon-bd-o"></i></a></li>
        <li><a @click="state=2"><i class="icon-bd-o"></i></a></li>
      </ul>


    </div>
    <!--<chart_line :dataItems="data_line" :face="face" :class="{hideitem:state==1}" ></chart_line>-->
    <!--<kchart :dataItems="data_kLine" :no="$route.params.no"  :class="{hideitem:state==0}"></kchart>-->
    <chart_line :dataItems="data_line" :face="face" v-if="state==0" ></chart_line>
    <kchart :dataItems="data_kLine" :no="$route.params.no"  v-if="state==1"></kchart>
    <chart_volume :dataItems="data_line" :face="face" v-if="state==2" ></chart_volume>
  </div>
</template>

<script>
  import kchart from '@/components/public/charts_k'
  import chart_line from '@/components/public/chart_line'
  import chart_volume from '@/components/public/chart_volume'

  var _this;

  export default {
    components:{kchart,chart_line,chart_volume},
    data:function () {
      return{
        state:1,
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
