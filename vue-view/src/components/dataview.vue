<template>
  <div class="col-12">
    <div class="panel">
      <div class="panel-heading border-warning">
        <div class="panel-title" style="width: 150px;float: left;font-size: 30px">数据列表</div>
        <el-date-picker
          v-model="queryParams.date"
          type="date"
          value-format="yyyy-MM-dd"
          style="float: left;margin-left: 20px"
          placeholder="选择日期">
        </el-date-picker>
        <el-input placeholder="请输入编号" v-model="queryParams.no" style="float: left;margin-left: 20px;width: 220px">
          <template slot="prepend">编号</template>
        </el-input>
        <el-button type="primary" icon="el-icon-search" @click="freashItems(1)">搜索</el-button>
      </div>
      <div class="panel-body">
        <table class="table table-striped">
          <thead>
          <tr>
            <td>编号</td>
            <td>日期</td>
            <td>最小值</td>
            <td>最大值</td>
            <td>开盘</td>
            <td>收盘</td>
            <td>涨跌</td>
            <td>换手率</td>
            <td>成交量</td>
            <td>成交额</td>
          </tr>
          </thead>
          <tbody>
          <tr v-for="item in codeitems"  @click="onrowclick(item)">
            <td>{{item.no}}</td>
            <td>{{item.date}}</td>
            <td>{{item.min}}</td>
            <td>{{item.max}}</td>
            <td>{{item.startprice}}</td>
            <td>{{item.lastprice}}</td>
            <td>{{item.ud}}</td>
            <td>{{item.turnoverRate}}</td>
            <td>{{item.volume}}</td>
            <td>{{item.turnover}}</td>
          </tr>
          </tbody>
        </table>
        <el-pagination
          background
          layout="prev, pager, next"
          :pager-count="11"
          @current-change="freashItems"
          :current-page="queryParams.pageIndex"
          :page-size="queryParams.pageSize"
          :total="total">
        </el-pagination>
      </div>
    </div>
  </div>
</template>

<script>
  var _this;

  export default {
    components:{},
    data:function () {
      return{

        queryParams:{
          pageIndex:1,
          pageSize:10,
          date:"",
          no:""
        },
        codeitems:[],

        total:0
      }
    },
    methods:{
      freashItems:function (index) {
        _this.queryParams.pageIndex=index;
        _this.$http.get({api:"DV_CODEPAGE",params:_this.queryParams}).then(data=>{
          _this.codeitems=data.content;
          _this.total=data.totalRecord;
          _this.queryParams.pageIndex=Number(data.pageIndex) ;
        },(res)=>{  _this.$store.commit('showAlertv', {text: res}) });
      },
      onrowclick(row){
        var path='#/chartPanel/'+row.no+'?date='+row.date
        window.location.href=path
      }
    },
    watch:{
      "queryParams.date":function (d) {
        if(d||d=="null")_this.queryParams.date=""
      }
    },

    created:function () {
      _this=this;
    },
    mounted:function () {
      _this.freashItems(1);
    }
  }
</script>
