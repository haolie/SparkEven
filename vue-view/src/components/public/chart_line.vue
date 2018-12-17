<template>
  <div ref="linechart" class="chart">

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
        this.chart=this.$echarts.init(this.$refs.linechart);
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
        this.dataItems.forEach(function (d) {
          values.push([getTime(d.time),d.price])
        })

        var mid=this.face.lastprice-this.face.ud
        var h=Math.max(Math.abs(this.face.max-mid),Math.abs(this.face.min-mid))


        var option = {
          title: {
            text: '动态数据 + 时间坐标轴'
          },
          grid:{
            width:'100%',
            height:'100%',
          },
          tooltip: {
            trigger: 'axis',
            formatter: function (params) {
              var temp=_this.dataItems[params[0].dataIndex]
              return '时间：'+params[0].axisValue+'<br/>'+'时间：'+temp.time+'<br/>'+'value：'+temp.price+'<br/>'+'volume：'+temp.volume+'<br/>'+'type：'+temp.trade_type;
            },
            axisPointer: {
              animation: false
            }
          },
          xAxis: {
            show:false,
            type: 'category',
            splitLine: {
              show: false
            }
          },
          yAxis: {

            type: 'value',
            min:mid-h,
            max:mid+h,
            boundaryGap: [0, '100%'],
            splitLine: {
              show: false
            }
          },
          series: [{
            name: '模拟数据',
            type: 'line',
            showSymbol: true,
            hoverAnimation: false,
            markLine: {
              data: [
                { yAxis: mid, name: '平均值'}
              ]
            },
            data: values
          }]
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
