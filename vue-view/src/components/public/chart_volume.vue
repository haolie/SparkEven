<template>
  <div ref="piechart" class="chart">

  </div>
</template>

<script>
  var _this;

  export default {
    components:{},
    props:['dataItems','face'],
    data:function () {
      return{
        chart:null
      }
    },
    methods:{
      setChart(){
        if(!this.dataItems||!this.face)return

        if(this.chart){
          this.chart.dispose()
        }
        this.chart=this.$echarts.init(this.$refs.piechart);
        var values=[]
        var getTime=function (time) {
          var s=time%60
          var m=((time-s)/60)%60
          var h=Math.ceil( (time-m*60-s)/(60*60))+9
          s=s>=10?s:'0'+s
          m=m>=10?m:'0'+m
          h=h>=10?h:'0'+h
          //return _this.face.date+' '+  h+':'+m+':'+s
          return  h+':'+m+':'+s
        }

        var datas=[{value:0,name:'m'},{value:0,name:'a'}]
        this.dataItems.forEach(function (d) {
          values.push([getTime(d.time),d.price])
          if(d.time>10000)datas[1].value+=d.volume
          else datas[0].value+=d.volume
        })

        var mid=this.face.lastprice-this.face.ud
        var h=Math.max(Math.abs(this.face.max-mid),Math.abs(this.face.min-mid))



        var option =  {
          title : {
            text: '某站点用户访问来源',
            subtext: '纯属虚构',
            x:'center'
          },
          tooltip : {
            trigger: 'item',
            formatter: "{a} <br/>{b} : {c} ({d}%)"
          },
          legend: {
            orient: 'vertical',
            left: 'left',
            data: ['m','a']
          },
          series : [
            {
              name: '访问来源',
              type: 'pie',
              radius : '55%',
              center: ['50%', '60%'],
              data:datas,
              itemStyle: {
                emphasis: {
                  shadowBlur: 10,
                  shadowOffsetX: 0,
                  shadowColor: 'rgba(0, 0, 0, 0.5)'
                }
              }
            }
          ]
        };

        this.chart.setOption(option)
      }
    },
    watch:{
      dataItems(){
        this.setChart()
      }
    },
    created:function () {
      _this=this;
    },
    mounted:function () {
      this.setChart()
    }
  }
</script>
