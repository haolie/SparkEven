CREATE TABLE time_price(
no MEDIUMINT NOT NULL,
time TIMESTAMP NOT NULL,
price MEDIUMINT,
trade_type TINYINT,
turnover_inc MEDIUMINT UNSIGNED,
volume MEDIUMINT UNSIGNED,
primary key (no,time)
);
CREATE TABLE codeface(
_no MEDIUMINT NOT NULL,
_date DATE NOT NULL,
_min MEDIUMINT,
_max MEDIUMINT,
ud MEDIUMINT,
lastprice MEDIUMINT,
face TINYINT,
dde INT,
dde_b INT UNSIGNED,
dde_s INT UNSIGNED,
mainforce INT UNSIGNED,
_state tinyint,
per smallint
primary key (_no,_date)
);

2.0

GRANT all privileges ON *.* TO 'mysql'@'%' IDENTIFIED BY '123456' WITH GRANT OPTION;

mysqld --initialize --user=mysql --console

CREATE TABLE tbl_codes(
id SMALLINT UNSIGNED NOT NULL,
_no MEDIUMINT UNSIGNED NOT NULL,
state tinyint DEFAULT 0,
primary key (id)
);

insert into tbl_codes(id,_no) values(0,"1912261");

CREATE TABLE time_price(
face_id MEDIUMINT UNSIGNED NOT NULL,
time SMALLINT UNSIGNED NOT NULL,
price SMALLINT,
trade_type TINYINT,
volume MEDIUMINT UNSIGNED,
primary key (face_id,time)
);

CREATE TABLE codeface(
id MEDIUMINT UNSIGNED NOT NULL,
_date DATE NOT NULL,
 no_id SMALLINT UNSIGNED NOT NULL,
_min MEDIUMINT UNSIGNED,
_max MEDIUMINT UNSIGNED,
_change MEDIUMINT,
lastprice MEDIUMINT UNSIGNED,
startprice MEDIUMINT UNSIGNED,
dde BIGINT,
dde_b BIGINT UNSIGNED,
dde_s BIGINT UNSIGNED,
face TINYINT UNSIGNED,
volume BIGINT UNSIGNED,
turnoverRate SMALLINT UNSIGNED,
turnover BIGINT UNSIGNED,
state TINYINT,
per smallint,
primary key (id)
);

alter table codeface add volume BIGINT UNSIGNED;
alter table codeface add turnoverRate SMALLINT UNSIGNED;
alter table codeface add turnover BIGINT UNSIGNED;
alter table codeface add INDEX no_id(no_id);

CREATE TABLE micInfo(
  no_id SMALLINT UNSIGNED NOT NULL,
  infoDate DATE NOT NULL,
  primary key (infoDate,no_id)
  );
  
 select count(*),_date from codeface group by _date;

--mysqld  --initialize
--mysqld --install
--net start MySQL
--net stop MySQL
--sc delete MySQL
--mysqld -remove

select b._change,b.lastprice,c.* from tbl_codes as a inner join codeface as b on a.id=b.no_id inner join time_price as c on b.id=c.face_id where a._no=1600960 and b._date='2018-11-09';

select b._change,b.lastprice,c.* from tbl_codes as a inner join codeface as b on a.id=b.no_id inner join time_price as c on b.id=c.face_id where a._no=1600960 and b._date='2018-11-09';

select a._no,b.id,b._date from tbl_codes as a inner join codeface as b on a.id=b.no_id where b._date='2020-01-03' and a._no=1912261;

select count(*) from time_price where face_id in(select c.face_id as face_id from tbl_codes as a inner join codeface as b on a.id=b.no_id inner join time_price as c on b.id=c.face_id where a._no=1600960 and b._date='2018-11-09');
mysqldump -u mysql -p test timeprice2018_1 > d:/db/timeprice2018_1.sql
delete from time_price where face_id=1704852;
update codeface set state=0 where id=1704852;
select id,state from codeface where id=1704852;  
select id,state from codeface where id=912261;  

//912261   1704858