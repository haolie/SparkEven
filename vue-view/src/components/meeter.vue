<template>
  <div>
    <div class="col-4">
      <div class="panel">
        <div class="panel-heading border-warning">
          <div class="panel-title">采集队列</div>
        </div>
        <div class="panel-body">
          <table class="table table-striped">
            <thead>
            <tr>
              <td>日期</td>
              <td>数量</td>
              <td>状态</td>
            </tr>
            </thead>
            <tbody>
            <tr v-for="item in dateItems" v-bind="{active:item.date==curDateItem.date}" @click="setCurDateItem(item)">
              <td>{{item.date}}</td>
              <td>{{item.count}}</td>
              <td>{{item.state}}</td>
            </tr>
            </tbody>
          </table>
          <el-pagination
            background
            layout="prev, pager, next"
            :pager-count="5"
            @current-change="freashDateItems"
            :current-page="datePageObj.index"
            :page-size="datePageObj.pageSize"
            :total="datePageObj.total">
          </el-pagination>
        </div>
      </div>
    </div>
    <div class="col-8">
      <div class="panel">
        <div class="panel-heading border-warning">
          <div class="panel-title">{{curDateItem?(curDateItem.date+'('+curDateItem.count+')'):""}}</div>
          <el-button type="warning" icon="el-icon-star-off" circle style="float: right; margin-top: -30px;" size="min" @click="codeViewState=!codeViewState"></el-button>
        </div>
        <div class="panel-body" v-if="!codeViewState">
          <table class="table table-striped">
            <thead>
            <tr>
              <td>编号</td>
              <td>状态</td>
              <td>下载</td>
              <td>保存</td>
            </tr>
            </thead>
            <tbody>
            <tr v-for="item in codeRows">
              <td>{{item.no}}</td>
              <td>{{item.state}}</td>
              <td>{{item.downstate}}</td>
              <td>{{item.savestate}}</td>
            </tr>
            </tbody>
          </table>
          <el-pagination
            background
            layout="prev, pager, next"
            :pager-count="5"
            @current-change="showCodeTable"
            :current-page="codePageObj.index"
            :page-size="codePageObj.pageSize"
            :total="codePageObj.total">
          </el-pagination>
        </div>
        <div class="panel-body" v-if="codeViewState">
          <div class="chart" ref="echart_codes">

          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script>
  var _this,interval=30000,pagesize=10;

  export default {
    data:function () {
      return{
        dateItems:[],
        curDateItem:{date:"",count:0,state:0},
        codeRows:[],
        dateViewState:0,
        codeViewState:0,
        datePageObj:{
          total:0,
          index:0,
          pageSize:pagesize
        },
        codePageObj:{
          total:0,
          index:0,
          pageSize:pagesize
        }
      }
    },
    methods:{
      dateState:function (state) {
        return ""
      },
      setCurDateItem:function (item) {
        item.items=[];
        _this.curDateItem=item;

        _this.freashCodeItems()
      },
      showCodeTable:function (index) {
        index:index?index:1;
        var obj= _this.$lib.pageList(_this.curDateItem.items,pagesize,index);
        _this.codeRows=  obj.content
        _this.codePageObj.index=obj.pageIndex;
        _this. freashCodeChart()
      },
      freashCodeChart:function () {
        if(!_this.curDateItem||!_this.curDateItem.items||!_this.curDateItem.items.length) return
        var echart_codes= _this.$echarts.init(_this.$refs.echart_codes);
        var names=["未处理"," 已下载","正在处理","已完成","错误"],datas=[];
        names.forEach(function (name,i) {
          datas.push({name:names[i],value:0})
        })
        for(var i=0;i<_this.curDateItem.items.length;i++){
          if(_this.curDateItem.items[i].savestate<0) datas[0].value+=1;
          else if(_this.curDateItem.items[i].savestate==0)datas[1].value+=1;
          else if(_this.curDateItem.items[i].savestate==1)datas[2].value+=1;
          else if(_this.curDateItem.items[i].savestate==2)datas[3].value+=1;
          else if(_this.curDateItem.items[i].savestate>2)datas[4].value+=1;
        }


        echart_codes.setOption({
          legend: {
            show:false,
            orient: 'vertical',
            left: 'left',
            data: names
          },
          series : [
            {
              name: '完成情况',
              type: 'pie',
              radius : '55%',
              center: ['50%', '60%'],
              data:datas,
              label:{
                normal:{
                  show:true, formatter:'{b}:{d}%({c})'
                }
              },
              itemStyle: {
                emphasis: {
                  shadowBlur: 10,
                  shadowOffsetX: 0,
                  shadowColor: 'rgba(0, 0, 0, 0.5)'
                }
              }
            }
          ]
        })


      },
      freashCodeItems:function () {
        if(!_this.curDateItem||!_this.curDateItem) return;
        _this.$http.get({api:"G_DATECODEITEM",params:{date:_this.curDateItem.date }}).then(data=>{
          _this.curDateItem.items=data;
          _this.codePageObj.total=data.length;
          _this.showCodeTable(1);
        },(res)=>{  _this.$store.commit('showAlertv', {text: res}) });
      },
      freashDateItems:function (index) {
        index=index?index:1;
        _this.$http.get({api:"G_DATEITEMLIST",params:{pageSize:pagesize,pageIndex:index}}).then(data=>{

          _this.dateItems=data.content;
          _this.datePageObj.total=data.totalRecord;
          _this.datePageObj.index=Number(data.pageIndex) ;
        },(res)=>{  _this.$store.commit('showAlertv', {text: res}) });
      }
    },
    created:function () {
      _this=this;
    },
    mounted:function () {
      var fun=function () {
        _this.freashDateItems( _this.datePageObj.index)
        setTimeout(fun,interval)
      }

      fun();
    }
  }
</script>
