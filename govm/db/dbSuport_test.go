package db

import(
	"testing"
)

var mysql MysqlSuport=MysqlSuport{}

func TestInit(t *testing.T){
	r:=mysql.Init()
	if !r{
		t.Error("Init Failed")
	}
}

func TestgetIdByNo(t *testing.T){
	mysql.Init()
	id:= mysql.GetIdbyNo(1912261)
	if id!=0 {
		t.Error("GetIdbyNo Failed")
	}
}

func TestGetNoById(t *testing.T){
	mysql.Init()
	no:= mysql.GetNoById(-44334)
	if no!=-1 {
		t.Error("GetNoById Failed")
	}

	no= mysql.GetNoById(1912261)

	if no!=-0 {
		t.Error("GetNoById Failed")
	}
}

func TestAddNewNo(t *testing.T){
	mysql.Init()
	r:= mysql.AddNewNo(1600002)
	if r<0 {
		t.Error("AddNewNo Failed")
	}
}