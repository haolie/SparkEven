<template>
  <div>

    <div id="dateCharts_1" style="width: 100%;height: 520px;">

    </div>

  </div>
</template>

<script>
  var _this;

  export default {
    components:{},
    data:function () {
      return{
        date:'',
        codeitems:[],
        total:0,
        listMap:[
          {
            min:-12,
            max:-8.09,
            list:[]
          },
          {
            max:-8.09,
            min:-6.18,
            list:[]
          },
          {
            min:-6.18,
            max:-3.82,
            list:[]
          },
          {
            min:-3.82,
            max:-1.91,
            list:[]
          },
          {
            min:-1.91,
            max:1.91,
            list:[]
          },
          {
            min:1.91,
            max:3.82,
            list:[]
          },
          {
            min:3.82,
            max:6.18,
            list:[]
          },
          {
            min:6.18,
            max:8.09,
            list:[]
          },
          {
            min:8.09,
            max:12,
            list:[]
          }]  //0: >=-10 && <-8.09  1:>=-8.09&& < -6.18  2:>=-6.18 && <-3.82  3:>=-3.82 && < 3.82  4:>=3.82 && <6.18  5:>=6.18 && <8.09  6:>=8.09:
      }
    },
    methods:{
      getData(){
        if(!this.date)return
        _this.$http.get({api:"DV_CODEPAGE",params:{date:this.date}}).then(data=>{
          _this.codeitems=data.content;
          _this.total=data.totalRecord;

          _this.listMap.forEach(function (list) {
            list.list=[]
          })
          _this.codeitems.forEach(function (item) {
            item.per=Math.floor((item.ud/(item.lastprice-item.ud))*10000)/100
            for(var i=0;i<_this.listMap.length;i++){
              if(i==_this.listMap.length-1){
                _this.listMap[i].list.push(item)
                break
              }
              else if(item.per<_this.listMap[i].max) {
                _this.listMap[i].list.push(item)
                break
              }
            }
          })

          _this.setChart()
        },(res)=>{  _this.$store.commit('showAlertv', {text: res}) });

      },
      setChart(){
        if( !this.codeitems.length)return
       var chart=this.$echarts.init(document.getElementById('dateCharts_1'));
        var names=[],values=[]
        for(var i=0;i<_this.listMap.length;i++){
          var o={
            name:_this.listMap[i].min+'-'+_this.listMap[i].max,
            value:_this.listMap[i].list.length
          }
          if(i==0){
            o.name='<'+_this.listMap[i].max
          }
          if(i==_this.listMap.length-1){
            o.name='>='+_this.listMap[i].min
          }
          names.push(o.name)
          values.push(o)
        }

        chart.setOption( {
          tooltip : {
            trigger: 'item',
            formatter: "{a} <br/>{b} : {c} ({d}%)"
          },
          legend: {
            type: 'scroll',
            orient: 'vertical',
            right: 10,
            top: 20,
            bottom: 20,
            data: names
          },
          series : [
            {
              type: 'pie',
              radius : '55%',
              center: ['50%', '50%'],
              data: values,
              itemStyle: {
                emphasis: {
                  shadowBlur: 10,
                  shadowOffsetX: 0,
                  shadowColor: 'rgba(0, 0, 0, 0.5)'
                }
              }
            }
          ]
        }
      )

      },
      getDate(){
        this.date= this.$route.params.date
      }
    },
    watch:{
      '$route'(){
        this.getData()
      },
      date(){
        this.getData()
      }
    },

    created:function () {
      _this=this;
    },
    mounted:function () {
      this.getDate()
    }
  }
</script>
