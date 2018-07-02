<template>
  <div>
    <el-tabs >
      <el-tab-pane label="信息管理" name="first">
        <div class="col-6">
          <div class="panel">
            <div class="panel-heading border-warning">
              <div class="panel-title">列表</div>
              <input type="file" style="display: none" @change="onFilePathChanged" id="upload_input"/>
            </div>
            <div class="panel-body">
              <div class="demo-input-suffix">
               <label style="float: left"></label>
                <el-input style="width: 160px;float: left"
                  placeholder="编号"
                  v-model="addItem.no">
                </el-input>
                <el-date-picker style="width: 160px;float: left"
                  v-model="addItem.date"
                  align="right"
                  type="date"
                  placeholder="选择日期"
                  >
                </el-date-picker>
                <el-button type="primary" style="width: 75px;float: left">添加</el-button>
                <el-button type="primary" style="width: 75px;float: left" @click="uploadFile()">上传</el-button>
              </div>
              <table class="table table-striped">
                <thead>
                <tr>
                  <td>编号</td>
                  <td>日期</td>
                </tr>
                </thead>
                <tbody>
                <tr v-for="item in itemRows">
                  <td>{{item.no}}</td>
                  <td>{{item.date}}</td>
                </tr>
                </tbody>
              </table>
              <el-pagination
                background
                layout="prev, pager, next"
                :pager-count="5"
                @current-change="freashItems"
                :current-page="queryParams.pageIndex"
                :page-size="queryParams.pageSize"
                :total="total">
              </el-pagination>
            </div>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="数据分析" name="second">

      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
  var _this;

  export default {
    components:{},
    data:function () {
      return{
        curItem:{},
        uploadFilePath:"",
        addItem:{
          no:"",
          date:""
        },
        itemRows:[],
        total:0,
        queryParams:{
          pageIndex:1,
          pageSize:10
        }
      }
    },
    methods:{
      freashItems:function (index) {
        _this.queryParams.pageIndex=index;
        _this.$http.get({api:"M_ITEMSLIST_PAGE",params:_this.queryParams}).then(data=>{
          _this.itemRows=data.content;
          _this.total=data.totalRecord;
          _this.queryParams.pageIndex=Number(data.pageIndex) ;
        },(res)=>{  _this.$store.commit('showAlertv', {text: res}) });
      },
      uploadFile:function (file) {

        document.getElementById("upload_input").click();
      },
      addItemfun:function (item) {
        if(_this.addItem.no&&_this.addItem.date){
          _this.$http.post({api:"M_ITEM",params:_this.addItem}).then(data=>{
            _this.addItem.no="";
            _this.addItem.date="";
          },(res)=>{  _this.$store.commit('showAlertv', {text: res}) });
        }
      },
      deleteItem:function () {

      },
      analysisItem:function (item) {

      },
      onFilePathChanged:function (e) {
        document.getElementById('upload_input').select()
        console.log(document.selection.createRange().text)
        if(e.target.value.length){
          _this.$http.get({api:"M_UPLOADFILE",params:{file:e.target.value}}).then(data=>{

          },(res)=>{  _this.$store.commit('showAlertv', {text: res}) });
        }
      },
      setCurDateItem:function (item) {

      }
    },
    watch:{},

    created:function () {
      _this=this;
    },
    mounted:function () {
      _this.freashItems(1);
    }
  }

  function getObjectURL(file) {
    var url = null;
    if (window.createObjcectURL != undefined) {
      url = window.createOjcectURL(file);
    } else if (window.URL != undefined) {
      url = window.URL.createObjectURL(file);
    } else if (window.webkitURL != undefined) {
      url = window.webkitURL.createObjectURL(file);
    }
    return url;
  }
</script>
