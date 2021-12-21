drop table if exists user;
create table user(
	id  int primary key auto_increment,
	name varchar(255) not null,
	password varchar(255) not null,
	sex varchar(255) not null,	
	phone varchar(11) not null unique,
	avatar blob not null,
	role int(1) not null
) ;

drop table if exists list;
create table list(
	id  int primary key auto_increment,
	name varchar(255) not null,
	data varchar(255) not null,
	place varchar(255) not null,	
	info varchar(255) not null,
	appointnum not null,
	teamA varchar(255) not null,
	teamB varchar(255) not null 
) ; 

drop table if exists user_gamelist;
create table  user_gamelist (
	id int primary key auto_increment,
	user_phone varchar(11) not null,
	list_id int not null,
	foreign key (user_phone) references user(phone),  
	foreign key (list_id) references list(id),
	unique(user_phone,list_id)     
);

drop table if exists team;
create table team (
	id int primary key auto_increment,
	name varchar(255) not null unique,
	logo varchar(255) not null,
	info text not null
);

drop table if exists player;
create table player (
	id int primary key auto_increment,
	name varchar(255) not null unique,
	avatar varchar(255) not null,
	team varchar(255) ,
	num varchar (3) ,
	position varchar(255) not null,
	age varchar(2) not null,
	foreign key (team) references team(name)
);

drop table if exists player_team;
create table player_team (
	id int primary key auto_increment,
	player_id int not null,
	team_id int not null,
	foreign key (player_id) references player(id),
	foreign key (team_id) references team(id),
	unique(player_id,team_id)
); 







