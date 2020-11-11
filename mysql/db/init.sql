use prod_db;

create table labels(
    id int primary key not null,
    label_text varchar(64) not null
);

create table tasks(
    id int primary key not null,
    name varchar(128) not null,
    description varchar(256) not null,
    submitTime datetime not null,
    label int not null,
    foreign key (label) references labels (id) on delete cascade
);

insert into
    labels(id,label_text)
values
    (0,"なし");

insert into
    labels(id,label_text)
values
    (1,"緊急");

insert into
    labels(id,label_text)
values
    (2,"期限なし");

insert into
    tasks(id,name,description,submitTime,label)
values
    (1,"睡眠","よく寝る","2020-08-12 21:00:00",1);

insert into
    tasks(id,name,description,submitTime,label)
values
    (2,"ウホウホミッドナイト","ウホウホ！！！？！？？！？！？！？！？！？ウホウホウホウホウホウホウホ！！！！ウホホホホホホホホホｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗウッキーーーーーーーーー！！！！！！！！！！！！！！！！！ｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗｗ","2020-08-31 00:00:00",2);