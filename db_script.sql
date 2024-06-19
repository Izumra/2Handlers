create table Users(
	Id serial primary key,
	Username varchar(30) not null Unique,
	Password varchar(30)not null
)

insert into Users(Username,Password)
values ('Izumra','Izumra17.')

create table News(
	Id serial primary key,
	Title text not null,
	Content text not null
)

insert into News(Title,Content)
values ('В центре города упал метеорит','В челябинске, в центре города упал метеорит')

create table Categories(
	Id serial primary key,
	Title varchar(50) not null
)

create table NewsCategories(
	NewsId serial references News(Id),
	CategoryId serial references Categories(Id)
)
select * from Users

insert into NewsCategories(NewsId,CategoryId)
values
(1,2),
(1,3),
(1,4),
(1,1),
(1,5)